package atomic

import (
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/Metadiv-Atomic-Engine/jwt"
	"github.com/Metadiv-Atomic-Engine/nanoid"
	"github.com/Metadiv-Atomic-Engine/sql"
	"github.com/gin-gonic/gin"
)

func NewContext[T any](ginCtx *gin.Context, module *Module) *Context[T] {
	var page *sql.Pagination
	var sort *sql.Sorting
	if ginCtx.Request.Method == "GET" {
		page = new(sql.Pagination)
		sort = new(sql.Sorting)
		ginCtx.ShouldBindQuery(page)
		ginCtx.ShouldBindQuery(sort)
	}

	locale := ginCtx.GetHeader(HEADER_LOCALE)
	if locale == "" {
		locale = "en"
	}

	return &Context[T]{
		Engine:  Engine,
		Module:  module,
		Gin:     ginCtx,
		DB:      Engine.DB,
		MEM:     Engine.MEM,
		Page:    page,
		Sort:    sort,
		Locale:  locale,
		Request: GinRequest[T](ginCtx),
		Logger:  &contextLogger{Module: module},
		startAt: time.Now(),
	}
}

/*
Get Jwt from the request header - Authorization
*/
func (ctx *Context[T]) Jwt() IJwt {
	token := ctx.Gin.GetHeader("Authorization")
	if token == "" {
		return nil
	}
	token = strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ReplaceAll(
				token, "Bearer ", ""), "bearer ", ""), "BEARER ", "")

	claims := jwt.ParseTokenByPublicPEM(token, ctx.Engine.PEM.Public)
	if claims == nil {
		return nil
	}

	j := new(Jwt)
	j.UserID = uint(claims.Get("user_id").(float64))
	j.WorkspaceID = uint(claims.Get("workspace_id").(float64))
	j.UUID = claims.Get("uuid").(string)
	j.IP = claims.Get("ip").(string)
	j.Agent = claims.Get("agent").(string)
	return j
}

/*
Get the workspace id from the jwt
*/
func (ctx *Context[T]) WorkspaceID() uint {
	j := ctx.Jwt()
	if j == nil {
		return 0
	}
	return j.GetWorkspaceID()
}

/*
Get the user id from the jwt
*/
func (ctx *Context[T]) UserID() uint {
	j := ctx.Jwt()
	if j == nil {
		return 0
	}
	return j.GetUserID()
}

/*
Get the IP from the request header
*/
func (ctx *Context[T]) IP() string {
	return ctx.Gin.ClientIP()
}

/*
Get the agent from the request header
*/
func (ctx *Context[T]) Agent() string {
	return ctx.Gin.Request.UserAgent()
}

/*
OK is a helper function to respond with a 200 status code.
*/
func (ctx *Context[T]) OK(data any, page ...*sql.Pagination) {
	if ctx.hasResp {
		ctx.Logger.ERROR("context has already responded")
		return
	}

	// pagination
	var pageResponse *sql.Pagination
	if len(page) > 0 {
		pageResponse = page[0]
	}
	ctx.Response = &Response{
		Success:    true,
		Duration:   time.Since(ctx.startAt).Milliseconds(),
		Pagination: pageResponse,
		Data:       data,
	}
	ctx.hasResp = true
}

/*
OKFile is a helper function to respond with a 200 status code and file.
*/
func (ctx *Context[T]) OKFile(bytes []byte, filename ...string) {
	if ctx.hasResp {
		ctx.Logger.ERROR("context has already responded")
		return
	}

	var name string
	if len(filename) == 0 || filename[0] == "" {
		name = nanoid.NewSafe()
	} else {
		name = filename[0]
	}

	ctx.Gin.Header("Content-Disposition", "filename="+name)
	ctx.Gin.Data(http.StatusOK, determineFileType(name), bytes)
	ctx.hasResp = true
	ctx.isFile = true
}

/*
OKDownload is a helper function to respond with a 200 status code and file.
*/
func (ctx *Context[T]) OKDownload(bytes []byte, filename ...string) {
	if ctx.hasResp {
		ctx.Logger.ERROR("context has already responded")
		return
	}

	var name string
	if len(filename) == 0 || filename[0] == "" {
		name = nanoid.NewSafe()
	} else {
		name = filename[0]
	}

	ctx.Gin.Header("Content-Disposition", "filename="+name)
	ctx.Gin.Data(http.StatusOK, "application/octet-stream", bytes)
	ctx.hasResp = true
	ctx.isFile = true
}

/*
Err is a helper function to respond with an error status code (400).
*/
func (ctx *Context[T]) Err(code string, additionalMessage ...string) {
	if ctx.hasResp {
		ctx.Logger.ERROR("context has already responded")
		return
	}

	ctx.Response = &Response{
		Success:  false,
		Duration: time.Since(ctx.startAt).Milliseconds(),
		Error: &ErrorResp{
			Code:    code,
			Message: selectErrorMessage(ctx, code, additionalMessage...),
		},
	}
	ctx.hasResp = true
}

/*
InternalServerErr is a helper function to respond with an error status code (500).
*/
func (ctx *Context[T]) InternalServerErr(additionalMessage ...string) {
	if ctx.hasResp {
		ctx.Logger.ERROR("context has already responded")
		return
	}

	ctx.Response = &Response{
		Success:  false,
		Duration: time.Since(ctx.startAt).Milliseconds(),
		Error: &ErrorResp{
			Code:    ERR_INTERNAL_SERVER_ERROR,
			Message: selectErrorMessage(ctx, ERR_INTERNAL_SERVER_ERROR, additionalMessage...),
		},
	}
	ctx.hasResp = true
}

/*
Unauthorized is a helper function to respond with an error status code (401).
*/
func (ctx *Context[T]) Unauthorized(additionalMessage ...string) {
	if ctx.hasResp {
		ctx.Logger.ERROR("context has already responded")
		return
	}

	ctx.Response = &Response{
		Success:  false,
		Duration: time.Since(ctx.startAt).Milliseconds(),
		Error: &ErrorResp{
			Code:    ERR_UNAUTHORIZED,
			Message: selectErrorMessage(ctx, ERR_UNAUTHORIZED, additionalMessage...),
		},
	}
	ctx.hasResp = true
}

/*
Forbidden is a helper function to respond with an error status code (403).
*/
func (ctx *Context[T]) Forbidden(additionalMessage ...string) {
	if ctx.hasResp {
		ctx.Logger.ERROR("context has already responded")
		return
	}

	ctx.Response = &Response{
		Success:  false,
		Duration: time.Since(ctx.startAt).Milliseconds(),
		Error: &ErrorResp{
			Code:    ERR_FORBIDDEN,
			Message: selectErrorMessage(ctx, ERR_FORBIDDEN, additionalMessage...),
		},
	}
	ctx.hasResp = true
}

func selectErrorMessage[T any](ctx *Context[T], code string, additionalMessage ...string) string {
	var message string
	err := ctx.Engine.Errors[code]
	if err == nil {
		ctx.Logger.ERROR("error code not found")
	} else {
		switch ctx.Locale {
		case "en":
			message = err.Eng
		case "zht":
			message = err.Zht
		case "zhs":
			message = err.Zhs
		default:
			message = err.Eng
		}
	}

	var addMessage string
	if len(additionalMessage) > 0 {
		addMessage = " " + strings.Join(additionalMessage, " ")
	}

	return message + addMessage
}

func determineFileType(filename string) string {
	ext := filepath.Ext(filename)
	switch ext {
	case ".aac":
		return "audio/aac"
	case ".abw":
		return "application/x-abiword"
	case ".arc":
		return "application/x-freearc"
	case ".avif":
		return "image/avif"
	case ".avi":
		return "video/x-msvideo"
	case ".azw":
		return "application/vnd.amazon.ebook"
	case ".bin":
		return "application/octet-stream"
	case ".bz":
		return "application/x-bzip"
	case ".bz2":
		return "application/x-bzip2"
	case ".cda":
		return "application/x-cdf"
	case ".csh":
		return "application/x-csh"
	case ".css":
		return "text/css"
	case ".csv":
		return "text/csv"
	case ".doc":
		return "application/msword"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".eot":
		return "application/vnd.ms-fontobject"
	case ".epub":
		return "application/epub+zip"
	case ".gz":
		return "application/gzip"
	case ".gif":
		return "image/gif"
	case ".htm":
		return "text/html"
	case ".html":
		return "text/html"
	case ".ico":
		return "image/vnd.microsoft.icon"
	case ".ics":
		return "text/calendar"
	case ".jar":
		return "application/java-archive"
	case ".jpeg":
		return "image/jpeg"
	case ".jpg":
		return "image/jpeg"
	case ".js":
		return "text/javascript"
	case ".json":
		return "application/json"
	case ".jsonld":
		return "application/ld+json"
	case ".mid":
		return "audio/midi"
	case ".midi":
		return "audio/midi"
	case ".mjs":
		return "text/javascript"
	case ".mp3":
		return "audio/mpeg"
	case ".mp4":
		return "video/mp4"
	case ".mpeg":
		return "video/mpeg"
	case ".mpkg":
		return "application/vnd.apple.installer+xml"
	case ".odp":
		return "application/vnd.oasis.opendocument.presentation"
	case ".ods":
		return "application/vnd.oasis.opendocument.spreadsheet"
	case ".odt":
		return "application/vnd.oasis.opendocument.text"
	case ".oga":
		return "audio/ogg"
	case ".ogv":
		return "video/ogg"
	case ".ogx":
		return "application/ogg"
	case ".opus":
		return "audio/opus"
	case ".otf":
		return "font/otf"
	case ".png":
		return "image/png"
	case ".pdf":
		return "application/pdf"
	case ".php":
		return "application/x-httpd-php"
	case ".ppt":
		return "application/vnd.ms-powerpoint"
	case ".pptx":
		return "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	case ".rar":
		return "application/vnd.rar"
	case ".rtf":
		return "application/rtf"
	case ".sh":
		return "application/x-sh"
	case ".svg":
		return "image/svg+xml"
	case ".tar":
		return "application/x-tar"
	case ".tif":
		return "image/tiff"
	case ".tiff":
		return "image/tiff"
	case ".ts":
		return "video/mp2t"
	case ".ttf":
		return "font/ttf"
	case ".txt":
		return "text/plain"
	case ".vsd":
		return "application/vnd.visio"
	case ".wav":
		return "audio/wav"
	case ".weba":
		return "audio/webm"
	case ".webm":
		return "video/webm"
	case ".webp":
		return "image/webp"
	case ".woff":
		return "font/woff"
	case ".woff2":
		return "font/woff2"
	case ".xhtml":
		return "application/xhtml+xml"
	case ".xls":
		return "application/vnd.ms-excel"
	case ".xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case ".xml":
		return "application/xml"
	case ".xul":
		return "application/vnd.mozilla.xul+xml"
	case ".zip":
		return "application/zip"
	case ".3gp":
		return "video/3gpp"
	case ".3g2":
		return "video/3gpp2"
	case ".7z":
		return "application/x-7z-compressed"
	default:
		return "application/octet-stream"
	}
}
