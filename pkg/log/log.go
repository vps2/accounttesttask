//Простая обёртка для стандартного логера, добавляющая в вывод теги [INFO], [WARNING],
//[ERROR], [FATAL] в зависимости от наименования метода выводаю
// Сделано для упрощения поиска по файлу.
package log
import (
	"io"
	"log"
	"strings"
)

const (
	spaceTag   = " "
	infoTag    = "[INFO]"
	warningTag = "[WARNING]"
	errorTag   = "[ERROR]"
	fatalTag   = "[FATAL]"
)

func Info(msg string) {
	log.Printf("%s %s%s", infoTag, msg, LineBreak)
}

func Warning(msg string) {
	log.Printf("%s %s%s", warningTag, msg, LineBreak)
}

func Error(msg string) {
	log.Printf("%s %s%s", errorTag, msg, LineBreak)
}

func Fatal(msg string) {
	log.Fatalf("%s, %s%s", fatalTag, msg, LineBreak)
}

func Infof(format string, v ...interface{}) {
	log.Printf(concat(infoTag, spaceTag, format), v...)
}

func Warningf(format string, v ...interface{}) {
	log.Printf(concat(warningTag, spaceTag, format), v...)
}

func Errorf(format string, v ...interface{}) {
	log.Printf(concat(errorTag, spaceTag, format), v...)
}

func Fatalf(format string, v ...interface{}) {
	log.Fatalf(concat(fatalTag, spaceTag, format), v...)
}

func SetOutput(w io.Writer) {
	log.SetOutput(w)
}

func concat(v ...string) string {
	builder := strings.Builder{}
	for _, val := range v {
		builder.WriteString(val)
	}

	return builder.String()
}
