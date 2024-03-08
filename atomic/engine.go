package atomic

import (
	"github.com/Metadiv-Atomic-Engine/jwt"
	"github.com/gin-gonic/gin"
)

/*
Engine is the core of the atomic package.
There is only one instance of the Engine object in the entire application.
*/
var Engine = newEngine()

func newEngine() *engine {

	priv, pub, err := jwt.CreateRSAKeyPair()
	if err != nil {
		panic(err)
	}

	e := &engine{
		Gin: gin.Default(),
		PEM: &pem{
			Private: priv,
			Public:  pub,
		},
		Handlers: &handlers{
			Apis: make(map[string]*ApiHandler),
			Wss:  make(map[string]*WsHandler),
			Jobs: make(map[string]*JobHandler),
			Mids: make(map[string]*MiddlewareHandler),
		},
		Endpoints: &endpoints{
			Apis:  make(map[string]*ApiEndpoint),
			Wss:   make(map[string]*WsEndpoint),
			Crons: make(map[string]*CronEndpoint),
			Inits: make(map[string]*InitEndpoint),
			Mids:  make(map[string]*MiddlewareEndpoint),
		},
		Errors:      make(map[string]*Error),
		Environment: make(map[string]*Environment),
		Modules:     make(map[string]*Module),
		DBMigrates:  make([]any, 0),
		MEMMigrates: make([]any, 0),
	}

	e.Logger = &engineLogger{
		Engine: e,
	}

	e.initEnvironments()
	e.initErrors()
	e.initDatabase()
	return e
}
