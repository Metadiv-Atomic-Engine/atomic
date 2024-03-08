package atomic

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

const dayFormat = "2006-01-02"

const dayTimeFormat = "2006-01-02 15:04:05"

type engineLogger struct {
	Engine *engine
}

func (l *engineLogger) INFO(msg ...any) {
	if len(msg) > 0 {
		msg = append([]any{fmt.Sprintf("[INFO] %s ", time.Now().Format(dayTimeFormat))}, msg...)
	}
	m := fmt.Sprint(msg...)
	fmt.Println(m)
	writeFile(m)
}

func (l *engineLogger) DEBUG(msg ...any) {
	if len(msg) > 0 {
		msg = append([]any{fmt.Sprintf("[DEBUG] %s ", time.Now().Format(dayTimeFormat))}, msg...)
	}
	if l.Engine.EnvString(GIN_MODE) != gin.ReleaseMode {
		fmt.Println(msg...)
	}
}

func (l *engineLogger) ERROR(msg ...any) {
	if len(msg) > 0 {
		msg = append([]any{fmt.Sprintf("[ERROR] %s ", time.Now().Format(dayTimeFormat))}, msg...)
	}
	m := fmt.Sprint(msg...)
	fmt.Println(m)
	writeFile(m)
}

type contextLogger struct {
	Module *Module
}

func (l *contextLogger) INFO(msg ...any) {
	if len(msg) > 0 {
		msg = append([]any{fmt.Sprintf("[INFO] %s {%s} ", time.Now().Format(dayTimeFormat), l.Module.Symbol)}, msg...)
	}
	m := fmt.Sprint(msg...)
	fmt.Println(m)
	writeFile(m)
}

func (l *contextLogger) DEBUG(msg ...any) {
	if len(msg) > 0 {
		msg = append([]any{fmt.Sprintf("[DEBUG] %s {%s} ", time.Now().Format(dayTimeFormat), l.Module.Symbol)}, msg...)
	}
	m := fmt.Sprint(msg...)
	fmt.Println(m)
	writeFile(m)
}

func (l *contextLogger) ERROR(msg ...any) {
	if len(msg) > 0 {
		msg = append([]any{fmt.Sprintf("[ERROR] %s {%s} ", time.Now().Format(dayTimeFormat), l.Module.Symbol)}, msg...)
	}
	m := fmt.Sprint(msg...)
	fmt.Println(m)
	writeFile(m)
}

func writeFile(msg string) {
	if _, err := os.Stat(getFile()); os.IsNotExist(err) {
		os.WriteFile(getFile(), []byte(""), os.ModePerm)
	}
	os.WriteFile(getFile(), []byte(msg), os.ModeAppend)
}

func getFile() string {
	return "./logs/" + time.Now().Format(dayFormat) + ".log"
}
