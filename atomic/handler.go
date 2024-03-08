package atomic

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func NewApiHandler[T any](
	module *Module,
	uuid string,
	note string,
	handler func(ctx *Context[T]),
	typescript *TypescriptOpt,
) (handlerUUID string) {
	if uuid == "" {
		panic("Error uuid is empty")
	}
	if _, ok := Engine.Handlers.Apis[uuid]; ok {
		panic(fmt.Sprintf("ApiEndpoint uuid %s already exists (module: %s)", uuid, module.Symbol))
	}
	Engine.Handlers.Apis[uuid] = &ApiHandler{
		UUID:       uuid,
		Module:     module,
		Handler:    apiToHandler[T](handler, module),
		Note:       note,
		Typescript: typescript,
	}
	return uuid
}

func NewWsHandler[T any](
	module *Module,
	uuid string,
	note string,
	handler func(ctx *Context[T], ws *websocket.Conn),
) (handlerUUID string) {
	if uuid == "" {
		panic("Error uuid is empty")
	}
	if _, ok := Engine.Handlers.Wss[uuid]; ok {
		panic(fmt.Sprintf("WsEndpoint uuid %s already exists (module: %s)", uuid, module.Symbol))
	}
	Engine.Handlers.Wss[uuid] = &WsHandler{
		UUID:    uuid,
		Module:  module,
		Handler: wsToHandler[T](handler, module),
		Note:    note,
	}
	return uuid
}

func NewJobHandler(
	module *Module,
	uuid string,
	note string,
	handler func(),
) (handlerUUID string) {
	if uuid == "" {
		panic("Error uuid is empty")
	}
	if _, ok := Engine.Handlers.Jobs[uuid]; ok {
		panic(fmt.Sprintf("CronJobEndpoint uuid %s already exists (module: %s)", uuid, module.Symbol))
	}
	Engine.Handlers.Jobs[uuid] = &JobHandler{
		UUID:    uuid,
		Module:  module,
		Handler: handler,
		Note:    note,
	}
	return uuid
}

func NewMiddlewareHandler[T any](
	module *Module,
	uuid string,
	note string,
	handler func(ctx *Context[T]),
) (handlerUUID string) {
	if uuid == "" {
		panic("Error uuid is empty")
	}
	if _, ok := Engine.Handlers.Mids[uuid]; ok {
		panic(fmt.Sprintf("MiddlewareEndpoint uuid %s already exists (module: %s)", uuid, module.Symbol))
	}
	Engine.Handlers.Mids[uuid] = &MiddlewareHandler{
		UUID:    uuid,
		Module:  module,
		Handler: middlewareToHandler(handler, module),
		Note:    note,
	}
	return uuid
}

func apiToHandler[T any](f func(ctx *Context[T]), module *Module) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c := NewContext[T](ctx, module)
		f(c)

		// if file is served, no need to respond
		if c.hasResp && c.isFile {
			return
		}

		// unexpected, service did not respond
		if !c.hasResp || c.Response == nil {
			ctx.JSON(500, gin.H{
				"message": "service did not respond",
			})
			return
		}

		// success
		if c.Response.Success {
			ctx.JSON(200, c.Response)
			return
		}

		// error, but no error object
		if c.Response.Error == nil {
			ctx.JSON(500, gin.H{
				"message": "service did not respond with error",
			})
			return
		}

		if c.Response.Error.Code == ERR_INTERNAL_SERVER_ERROR {
			ctx.JSON(500, c.Response)
			return
		}

		if c.Response.Error.Code == ERR_UNAUTHORIZED {
			ctx.JSON(401, c.Response)
			return
		}

		if c.Response.Error.Code == ERR_FORBIDDEN {
			ctx.JSON(403, c.Response)
			return
		}

		// error
		ctx.JSON(400, c.Response)
	}
}

func wsToHandler[T any](f func(ctx *Context[T], ws *websocket.Conn), module *Module) gin.HandlerFunc {
	wsUpGrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	return func(c *gin.Context) {
		ws, err := wsUpGrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer ws.Close()

		ctx := NewContext[T](c, module)
		f(ctx, ws)
	}
}

func middlewareToHandler[T any](f func(ctx *Context[T]), module *Module) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c := NewContext[T](ctx, module)
		f(c)
	}
}
