package atomic

import "time"

type ApiEndpoint struct {
	UUID    string        `json:"uuid"`
	Route   string        `json:"route"`
	Method  string        `json:"method"`
	Cache   *CacheOpt     `json:"cache"`
	Rate    *RateLimitOpt `json:"rate"`
	Handler *ApiHandler   `json:"handler"`
}

type CacheOpt struct {
	Duration time.Duration `json:"duration"`
}

type RateLimitOpt struct {
	Rate     int64         `json:"rate"`
	Duration time.Duration `json:"duration"`
}

type WsEndpoint struct {
	UUID    string     `json:"uuid"`
	Route   string     `json:"route"`
	Handler *WsHandler `json:"handler"`
}

type CronEndpoint struct {
	UUID        string      `json:"uuid"`
	Spec        string      `json:"spec"`
	InitExecute bool        `json:"init_execute"`
	Handler     *JobHandler `json:"handler"`
}

type InitEndpoint struct {
	UUID    string      `json:"uuid"`
	After   bool        `json:"after"`
	Handler *JobHandler `json:"handler"`
}

type MiddlewareEndpoint struct {
	UUID        string             `json:"uuid"`
	MatchRoutes []string           `json:"match_routes"`
	SkipRoutes  []string           `json:"skip_routes"`
	Handler     *MiddlewareHandler `json:"-"`
}
