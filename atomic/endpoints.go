package atomic

import "net/http"

func NewGetApi(
	uuid string,
	handlerUUID string,
	route string,
	cache *CacheOpt,
	rate *RateLimitOpt,
) {
	handler := Engine.Handlers.Apis[handlerUUID]
	if handler == nil {
		panic("ApiHandler not found: " + handlerUUID)
	}
	if _, ok := Engine.Endpoints.Apis[uuid]; ok {
		panic("ApiEndpoint uuid already exists: " + uuid)
	}
	Engine.Endpoints.Apis[uuid] = &ApiEndpoint{
		UUID:    uuid,
		Route:   route,
		Method:  http.MethodGet,
		Cache:   cache,
		Rate:    rate,
		Handler: handler,
	}
}

func NewPostApi(
	uuid string,
	handlerUUID string,
	route string,
	cache *CacheOpt,
	rate *RateLimitOpt,
) {
	handler := Engine.Handlers.Apis[handlerUUID]
	if handler == nil {
		panic("ApiHandler not found: " + handlerUUID)
	}
	if _, ok := Engine.Endpoints.Apis[uuid]; ok {
		panic("ApiEndpoint uuid already exists: " + uuid)
	}
	Engine.Endpoints.Apis[uuid] = &ApiEndpoint{
		UUID:    uuid,
		Route:   route,
		Method:  http.MethodPost,
		Cache:   cache,
		Rate:    rate,
		Handler: handler,
	}
}

func NewPutApi(
	uuid string,
	handlerUUID string,
	route string,
	cache *CacheOpt,
	rate *RateLimitOpt,
) {
	handler := Engine.Handlers.Apis[handlerUUID]
	if handler == nil {
		panic("ApiHandler not found: " + handlerUUID)
	}
	if _, ok := Engine.Endpoints.Apis[uuid]; ok {
		panic("ApiEndpoint uuid already exists: " + uuid)
	}
	Engine.Endpoints.Apis[uuid] = &ApiEndpoint{
		UUID:    uuid,
		Route:   route,
		Method:  http.MethodPut,
		Cache:   cache,
		Rate:    rate,
		Handler: handler,
	}
}

func NewDeleteApi(
	uuid string,
	handlerUUID string,
	route string,
	cache *CacheOpt,
	rate *RateLimitOpt,
) {
	handler := Engine.Handlers.Apis[handlerUUID]
	if handler == nil {
		panic("ApiHandler not found: " + handlerUUID)
	}
	if _, ok := Engine.Endpoints.Apis[uuid]; ok {
		panic("ApiEndpoint uuid already exists: " + uuid)
	}
	Engine.Endpoints.Apis[uuid] = &ApiEndpoint{
		UUID:    uuid,
		Route:   route,
		Method:  http.MethodDelete,
		Cache:   cache,
		Rate:    rate,
		Handler: handler,
	}
}

func NewWsApi(
	uuid string,
	handlerUUID string,
	route string,
) {
	handler := Engine.Handlers.Wss[handlerUUID]
	if handler == nil {
		panic("WsHandler not found: " + handlerUUID)
	}
	if _, ok := Engine.Endpoints.Wss[uuid]; ok {
		panic("WsEndpoint uuid already exists: " + uuid)
	}
	Engine.Endpoints.Wss[uuid] = &WsEndpoint{
		UUID:    uuid,
		Route:   route,
		Handler: handler,
	}
}

func NewCronJob(
	uuid string,
	handlerUUID string,
	spec string,
	initExecute bool,
) {
	handler := Engine.Handlers.Jobs[handlerUUID]
	if handler == nil {
		panic("JobHandler not found: " + handlerUUID)
	}
	if _, ok := Engine.Endpoints.Crons[uuid]; ok {
		panic("CronEndpoint uuid already exists: " + uuid)
	}
	Engine.Endpoints.Crons[uuid] = &CronEndpoint{
		UUID:        uuid,
		Spec:        spec,
		InitExecute: initExecute,
		Handler:     handler,
	}
}

func NewInitJob(
	uuid string,
	handlerUUID string,
	after bool,
) {
	handler := Engine.Handlers.Jobs[handlerUUID]
	if handler == nil {
		panic("JobHandler not found: " + handlerUUID)
	}
	if _, ok := Engine.Endpoints.Inits[uuid]; ok {
		panic("InitEndpoint uuid already exists: " + uuid)
	}
	Engine.Endpoints.Inits[uuid] = &InitEndpoint{
		UUID:    uuid,
		After:   after,
		Handler: handler,
	}
}

func NewMiddleware(
	uuid string,
	handlerUUID string,
	matchRoutes []string,
	skipRoutes []string,
) {
	handler := Engine.Handlers.Mids[handlerUUID]
	if handler == nil {
		panic("MiddlewareHandler not found: " + handlerUUID)
	}
	if _, ok := Engine.Endpoints.Mids[uuid]; ok {
		panic("MiddlewareEndpoint uuid already exists: " + uuid)
	}
	Engine.Endpoints.Mids[uuid] = &MiddlewareEndpoint{
		UUID:        uuid,
		MatchRoutes: matchRoutes,
		SkipRoutes:  skipRoutes,
		Handler:     handler,
	}
}
