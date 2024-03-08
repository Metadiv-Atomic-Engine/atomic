package atomic

import (
	"time"

	"github.com/Metadiv-Atomic-Engine/sql"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Context[T any] struct {
	Engine *engine
	Gin    *gin.Context
	Module *Module

	DB  *gorm.DB
	MEM *gorm.DB

	Page     *sql.Pagination
	Sort     *sql.Sorting
	Locale   string
	Request  *T
	Response *Response

	Logger *contextLogger

	/*
		These fields only be used in this package.
	*/
	startAt time.Time
	hasResp bool
	isFile  bool
}

type Response struct {
	Success    bool            `json:"success"`
	Duration   int64           `json:"duration"`
	Pagination *sql.Pagination `json:"pagination,omitempty"`
	Error      *ErrorResp      `json:"error,omitempty"`
	Data       any             `json:"data,omitempty"`
}

type ErrorResp struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
