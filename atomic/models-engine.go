package atomic

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

/*
engine is the core of the atomic package.
We do not want to expose the engine to the outside world.
Instead, we want to expose a Engine object to the outside world.
*/
type engine struct {
	Gin *gin.Engine `json:"-"`
	DB  *gorm.DB    `json:"-"`
	MEM *gorm.DB    `json:"-"`

	PEM         *pem                    `json:"pem"`
	Handlers    *handlers               `json:"handlers"`
	Endpoints   *endpoints              `json:"endpoints"`
	Errors      map[string]*Error       `json:"errors"`      // uuid -> error
	Environment map[string]*Environment `json:"environment"` // key -> environment
	Logger      *engineLogger           `json:"-"`

	DBMigrates  []any `json:"-"`
	MEMMigrates []any `json:"-"`

	Modules map[string]*Module `json:"modules"` // symbol -> module
}

type pem struct {
	Private string `json:"private"`
	Public  string `json:"public"`
}

type handlers struct {
	Apis map[string]*ApiHandler        `json:"apis"` // uuid -> handler
	Wss  map[string]*WsHandler         `json:"wss"`  // uuid -> handler
	Jobs map[string]*JobHandler        `json:"jobs"` // uuid -> handler
	Mids map[string]*MiddlewareHandler `json:"mids"` // uuid -> handler
}

type endpoints struct {
	Apis  map[string]*ApiEndpoint        `json:"apis"`  // uuid -> endpoint
	Wss   map[string]*WsEndpoint         `json:"wss"`   // uuid -> endpoint
	Crons map[string]*CronEndpoint       `json:"crons"` // uuid -> endpoint
	Inits map[string]*InitEndpoint       `json:"inits"` // uuid -> endpoint
	Mids  map[string]*MiddlewareEndpoint `json:"mids"`  // uuid -> endpoint
}
