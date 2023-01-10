package converter

import (
	"fmt"
	"log"
	"rocket-nano/internal/util/convert"
	"runtime"
	"runtime/debug"
)

func Exit(v interface{}, args ...interface{}) {
	switch v := v.(type) {
	case error:
		err := v
		pc := make([]uintptr, 1)
		runtime.Callers(2, pc)
		funcName := runtime.FuncForPC(pc[0]).Name()
		funcName = convert.LastPart(funcName, "rocket-nano/")

		_, file, line, ok := runtime.Caller(1)
		if !ok {
			log.Fatal(fmt.Errorf("get file & line failed"))
		}
		file = convert.LastPart(file, "rocket-nano/")
		fileLine := fmt.Sprintf("%s:%d", file, line)

		newErr := fmt.Errorf("\n - line: %s\n - func: %s \n - error: %w \n%v",
			fileLine, funcName, err, string(debug.Stack()))

		log.Fatal(newErr)
	case string:
		format := v
		pc := make([]uintptr, 1)
		runtime.Callers(2, pc)
		funcName := runtime.FuncForPC(pc[0]).Name()
		funcName = convert.LastPart(funcName, "rocket-nano/")

		_, file, line, ok := runtime.Caller(1)
		if !ok {
			log.Fatal(fmt.Errorf("get file & line failed"))
		}
		file = convert.LastPart(file, "rocket-nano/")
		fileLine := fmt.Sprintf("%s:%d", file, line)

		err := fmt.Errorf(format, args...)
		newErr := fmt.Errorf("\n - line: %s\n - func: %s \n - error: %w \n%v",
			fileLine, funcName, err, string(debug.Stack()))

		log.Fatal(newErr)
	}
}
