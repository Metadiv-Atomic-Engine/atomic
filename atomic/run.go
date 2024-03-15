package atomic

import (
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	limiter "github.com/ulule/limiter/v3"
	_gin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func (e *engine) Run() {
	e.registerCORS()
	e.dependencyChecking()
	e.databaseMigrate()
	e.executeBeforeJobs()
	e.registerApis()
	e.registerCronJobs()
	e.executeAfterJobs()
	e.Gin.Run(e.EnvString(GIN_HOST) + ":" + e.EnvString(GIN_PORT))
}

func (e *engine) registerCORS() {
	e.Gin.Use(cors.New(cors.Config{
		AllowOrigins: strings.Split(e.EnvString(CORS_ALLOWED_ORIGINS), ","),
		AllowMethods: strings.Split(e.EnvString(CORS_ALLOWED_METHODS), ","),
		AllowHeaders: strings.Split(e.EnvString(CORS_ALLOWED_HEADERS), ","),
	}))
}

func (e *engine) databaseMigrate() {
	if len(e.DBMigrates) > 0 {
		if err := e.DB.AutoMigrate(e.DBMigrates...); err != nil {
			panic(err)
		}
	}
	if len(e.MEMMigrates) > 0 {
		if err := e.MEM.AutoMigrate(e.MEMMigrates...); err != nil {
			panic(err)
		}
	}
}

func (e *engine) executeBeforeJobs() {
	for _, job := range e.Endpoints.Inits {
		if !job.After {
			job.Handler.Handler()
		}
	}
}

func (e *engine) executeAfterJobs() {
	for _, job := range e.Endpoints.Inits {
		if job.After {
			job.Handler.Handler()
		}
	}
}

func (e *engine) registerCronJobs() {
	cron := cron.New()
	for i := range e.Endpoints.Crons {
		cron.AddFunc(e.Endpoints.Crons[i].Spec, e.Endpoints.Crons[i].Handler.Handler)
		if e.Endpoints.Crons[i].InitExecute {
			e.Endpoints.Crons[i].Handler.Handler()
		}
	}
	cron.Start()
}

func (e *engine) registerApis() {
	for uuid, api := range e.Endpoints.Apis {

		/*
			Request route
		*/
		route := strings.TrimRight(api.Route, "/")

		handlers := make([]gin.HandlerFunc, 0)

		/*
			Rate limit
		*/
		if api.Rate != nil {
			handlers = append([]gin.HandlerFunc{
				_gin.NewMiddleware(limiter.New(memory.NewStore(), limiter.Rate{Period: api.Rate.Duration, Limit: api.Rate.Rate})),
			}, handlers...)
		}

		/*
			Middlewares
		*/
		for key, mid := range e.Endpoints.Mids {
			var skip = false
			for i := range mid.SkipRoutes {
				if match, _ := regexp.Match(mid.SkipRoutes[i], []byte(route)); match {
					skip = true
					break
				}
			}
			if skip {
				break
			}

			var match bool = false
			for i := range mid.MatchRoutes {
				if match, _ = regexp.Match(mid.MatchRoutes[i], []byte(route)); match {
					break
				}
			}
			if !match {
				continue
			}

			handlers = append([]gin.HandlerFunc{e.Endpoints.Mids[key].Handler.Handler}, handlers...)
		}

		/*
			Cache
		*/
		if api.Cache != nil {
			handlers = append(handlers, cache.CachePage(persistence.NewInMemoryStore(time.Second), api.Cache.Duration, e.Endpoints.Apis[uuid].Handler.Handler))
		} else {
			handlers = append(handlers, e.Endpoints.Apis[uuid].Handler.Handler)
		}

		/*
			Methods
		*/
		switch api.Method {
		case http.MethodGet:
			e.Gin.GET(route, handlers...)
		case http.MethodPost:
			e.Gin.POST(route, handlers...)
		case http.MethodPut:
			e.Gin.PUT(route, handlers...)
		case http.MethodDelete:
			e.Gin.DELETE(route, handlers...)
		}
	}
}

func (e *engine) dependencyChecking() {
	var installedMaps = make(map[string]*Module) // symbol -> module
	for key := range e.Modules {
		installedMaps[e.Modules[key].Symbol] = e.Modules[key]
	}

	for key := range e.Modules {
		for _, dep := range e.Modules[key].Dependencies {
			if _, ok := installedMaps[dep.Symbol]; !ok {
				panic("Module {" + e.Modules[key].Symbol + "} has dependency {" + dep.Symbol + "} but not installed")
			}
		}
	}
}
