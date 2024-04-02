package atomic

import (
	"github.com/Metadiv-Atomic-Engine/sql"
	"github.com/gin-gonic/gin"
)

func (e *engine) initDatabase() {

	if e.EnvString(GIN_MODE) == gin.ReleaseMode || e.EnvString(DB_TYPE) == "mysql" {
		DB, err := sql.MySQL(
			e.EnvString(MYSQL_HOST),
			e.EnvString(MYSQL_PORT),
			e.EnvString(MYSQL_USERNAME),
			e.EnvString(MYSQL_PASSWORD),
			e.EnvString(MYSQL_DATABASE),
			e.EnvBool(DB_SILENT),
		)
		if err != nil {
			panic(err)
		}
		e.DB = DB

		MEM, err := sql.SqliteMem(e.EnvBool(DB_SILENT))
		if err != nil {
			panic(err)
		}
		e.MEM = MEM

		return
	}

	DB, err := sql.Sqlite("test.db", e.EnvBool(DB_SILENT))
	if err != nil {
		panic(err)
	}
	e.DB = DB

	MEM, err := sql.Sqlite("mem.db", e.EnvBool(DB_SILENT))
	if err != nil {
		panic(err)
	}
	e.MEM = MEM
}
