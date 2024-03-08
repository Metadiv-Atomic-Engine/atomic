package atomic

import "github.com/gin-gonic/gin"

type ApiHandler struct {
	UUID       string          `json:"uuid"`
	Module     *Module         `json:"module"`
	Handler    gin.HandlerFunc `json:"-"`
	Note       string          `json:"note"`
	Typescript *TypescriptOpt  `json:"typescript"`
}

type TypescriptOpt struct {
	Models       []any    `json:"models"`
	FunctionName string   `json:"function_name"`
	Paths        []string `json:"paths"`
	Forms        []string `json:"forms"`
	Body         string   `json:"body"`
	Response     string   `json:"response"`
}

type WsHandler struct {
	UUID    string          `json:"uuid"`
	Module  *Module         `json:"module"`
	Handler gin.HandlerFunc `json:"-"`
	Note    string          `json:"note"`
}

type JobHandler struct {
	UUID    string  `json:"uuid"`
	Module  *Module `json:"module"`
	Handler func()  `json:"-"`
	Note    string  `json:"note"`
}

type MiddlewareHandler struct {
	UUID    string          `json:"uuid"`
	Module  *Module         `json:"module"`
	Handler gin.HandlerFunc `json:"-"`
	Note    string          `json:"note"`
}
