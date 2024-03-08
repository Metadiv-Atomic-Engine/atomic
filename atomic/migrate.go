package atomic

func NewDBMigration(models ...any) {
	Engine.DBMigrates = append(Engine.DBMigrates, models...)
}

func NewMEMMigration(models ...any) {
	Engine.MEMMigrates = append(Engine.MEMMigrates, models...)
}
