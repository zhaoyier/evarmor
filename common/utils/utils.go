package utils

import "runtime"

// GetCallerName get caller function name
func GetCallerName() string {
	pc, _, _, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name()

	// pc := make([]uintptr, 1)
	// runtime.Callers(2, pc)
	// f := runtime.FuncForPC(pc[0])
	// return f.Name()
}
