package atomic

import (
	"github.com/Metadiv-Atomic-Engine/env"
	"github.com/gin-gonic/gin"
)

const (
	GIN_MODE             = "GIN_MODE"
	GIN_HOST             = "GIN_HOST"
	GIN_PORT             = "GIN_PORT"
	MYSQL_HOST           = "MYSQL_HOST"
	MYSQL_PORT           = "MYSQL_PORT"
	MYSQL_USERNAME       = "MYSQL_USERNAME"
	MYSQL_PASSWORD       = "MYSQL_PASSWORD"
	MYSQL_DATABASE       = "MYSQL_DATABASE"
	DB_TYPE              = "DB_TYPE" // mysql, sqlite
	DB_ENCRYPT_KEY       = "DB_ENCRYPT_KEY"
	DB_SILENT            = "DB_SILENT"
	CORS_ALLOWED_ORIGINS = "CORS_ALLOWED_ORIGINS"
	CORS_ALLOWED_METHODS = "CORS_ALLOWED_METHODS"
	CORS_ALLOWED_HEADERS = "CORS_ALLOWED_HEADERS"
)

func NewEnvString(module *Module, key string, note string, defaultValue ...string) (envKey string) {
	if _, ok := Engine.Environment[key]; ok {
		panic("Environment key already exists: " + key)
	}
	value := env.String(key, defaultValue...)
	Engine.Environment[key] = &Environment{
		Key:    key,
		Type:   "string",
		Value:  value,
		Note:   note,
		Module: module,
	}
	return key
}

func NewEnvInt(module *Module, key string, note string, defaultValue ...int) (envKey string) {
	if _, ok := Engine.Environment[key]; ok {
		panic("Environment key already exists: " + key)
	}
	value := env.Int(key, defaultValue...)
	Engine.Environment[key] = &Environment{
		Key:    key,
		Type:   "int",
		Value:  value,
		Note:   note,
		Module: module,
	}
	return key
}

func NewBool(module *Module, key string, note string, defaultValue ...bool) (envKey string) {
	if _, ok := Engine.Environment[key]; ok {
		panic("Environment key already exists: " + key)
	}
	value := env.Bool(key, defaultValue...)
	Engine.Environment[key] = &Environment{
		Key:    key,
		Type:   "bool",
		Value:  value,
		Note:   note,
		Module: module,
	}
	return key
}

func NewFloat(module *Module, key string, note string, defaultValue ...float64) (envKey string) {
	if _, ok := Engine.Environment[key]; ok {
		panic("Environment key already exists: " + key)
	}
	value := env.Float64(key, defaultValue...)
	Engine.Environment[key] = &Environment{
		Key:    key,
		Type:   "float",
		Value:  value,
		Note:   note,
		Module: module,
	}
	return key
}

func (e *engine) EnvString(key string) string {
	env := e.Environment[key]
	if env == nil {
		panic("Environment key not found: " + key)
	}
	if env.Type != "string" {
		panic("Environment key is not a string: " + key)
	}
	return env.Value.(string)
}

func (e *engine) EnvInt(key string) int {
	env := e.Environment[key]
	if env == nil {
		panic("Environment key not found: " + key)
	}
	if env.Type != "int" {
		panic("Environment key is not an int: " + key)
	}
	return env.Value.(int)
}

func (e *engine) EnvBool(key string) bool {
	env := e.Environment[key]
	if env == nil {
		panic("Environment key not found: " + key)
	}
	if env.Type != "bool" {
		panic("Environment key is not a bool: " + key)
	}
	return env.Value.(bool)
}

func (e *engine) EnvFloat(key string) float64 {
	env := e.Environment[key]
	if env == nil {
		panic("Environment key not found: " + key)
	}
	if env.Type != "float" {
		panic("Environment key is not a float: " + key)
	}
	return env.Value.(float64)
}

func (e *engine) initEnvironments() {
	/*
		Gin
	*/
	e.Environment[GIN_MODE] = &Environment{
		Key:   GIN_MODE,
		Type:  "string",
		Value: env.String(GIN_MODE, "debug"),
		Note:  "Gin mode",
	}
	e.Environment[GIN_HOST] = &Environment{
		Key:   GIN_HOST,
		Type:  "string",
		Value: env.String(GIN_HOST, "127.0.0.1"),
		Note:  "Gin host",
	}
	e.Environment[GIN_PORT] = &Environment{
		Key:   GIN_PORT,
		Type:  "string",
		Value: env.String(GIN_PORT, "5000"),
		Note:  "Gin port",
	}

	/*
		CORS
	*/
	e.Environment[CORS_ALLOWED_ORIGINS] = &Environment{
		Key:   CORS_ALLOWED_ORIGINS,
		Type:  "string",
		Value: env.String(CORS_ALLOWED_ORIGINS, "*"),
		Note:  "CORS allowed origins",
	}
	e.Environment[CORS_ALLOWED_METHODS] = &Environment{
		Key:   CORS_ALLOWED_METHODS,
		Type:  "string",
		Value: env.String(CORS_ALLOWED_METHODS, "GET,POST,PUT,DELETE,PATCH,OPTIONS"),
		Note:  "CORS allowed methods",
	}
	e.Environment[CORS_ALLOWED_HEADERS] = &Environment{
		Key:   CORS_ALLOWED_HEADERS,
		Type:  "string",
		Value: env.String(CORS_ALLOWED_HEADERS, "Origin,Content-Length,Content-Type,Authorization,X-Locale"),
		Note:  "CORS allowed headers",
	}

	/*
		Database
	*/
	e.Environment[DB_TYPE] = &Environment{
		Key:   DB_TYPE,
		Type:  "string",
		Value: env.String(DB_TYPE, ""),
		Note:  "Database type",
	}
	e.Environment[MYSQL_HOST] = &Environment{
		Key:   MYSQL_HOST,
		Type:  "string",
		Value: env.String(MYSQL_HOST, "127.0.0.1"),
		Note:  "MySQL host",
	}
	e.Environment[MYSQL_PORT] = &Environment{
		Key:   MYSQL_PORT,
		Type:  "string",
		Value: env.String(MYSQL_PORT, "3306"),
		Note:  "MySQL port",
	}
	e.Environment[MYSQL_USERNAME] = &Environment{
		Key:   MYSQL_USERNAME,
		Type:  "string",
		Value: env.String(MYSQL_USERNAME, "root"),
		Note:  "MySQL username",
	}
	e.Environment[MYSQL_PASSWORD] = &Environment{
		Key:   MYSQL_PASSWORD,
		Type:  "string",
		Value: env.String(MYSQL_PASSWORD, "root"),
		Note:  "MySQL password",
	}
	e.Environment[MYSQL_DATABASE] = &Environment{
		Key:   MYSQL_DATABASE,
		Type:  "string",
		Value: env.String(MYSQL_DATABASE, "atomic"),
		Note:  "MySQL database",
	}

	var defaultSilent bool = e.EnvString(GIN_MODE) == gin.ReleaseMode
	e.Environment[DB_SILENT] = &Environment{
		Key:   DB_SILENT,
		Type:  "bool",
		Value: env.Bool(DB_SILENT, defaultSilent),
		Note:  "Database silent mode",
	}

	var encryptKey string
	if e.EnvString(GIN_MODE) == gin.ReleaseMode {
		encryptKey = env.String(DB_ENCRYPT_KEY, "")
		if encryptKey == "" {
			panic("DB_ENCRYPT_KEY is required in release mode")
		}
	} else {
		encryptKey = env.String(DB_ENCRYPT_KEY, "not_secret")
	}

	e.Environment[DB_ENCRYPT_KEY] = &Environment{
		Key:   DB_ENCRYPT_KEY,
		Type:  "string",
		Value: encryptKey,
		Note:  "Database encrypt key",
	}
}
