package try

import (
	"fmt"
	"log/slog"
	"reflect"
	"runtime"
)

func Panic(err error) {
	e(err)
}

func Wrap(err error) error {
	we := wrapError{error: err}
	runtime.Callers(2, we.pc[:])
	return we
}

func Wrapf(format string, a ...any) error {
	we := wrapError{error: fmt.Errorf(format, a...)}
	runtime.Callers(2, we.pc[:])
	return we
}

func Catch(fn ...func(err error)) {
	if len(fn) == 0 {
		r(recover(), func(w wrapError) {
			DefaultCatchHandler(w)
		})
		return
	}

	r(recover(), func(w wrapError) {
		fn[0](w)
	})
}

var DefaultCatchHandler = func(err wrapError) {
	slog.Error("try: recovered: "+err.Error(), "stack", getStackTrace(err.pc[:]))
}

const MaxDepth = 32

func getStackTrace(stack []uintptr) []string {
	throwList := make([]string, 0, MaxDepth)
	frames := runtime.CallersFrames(stack)
	for {
		frame, more := frames.Next()
		if !more {
			break
		}

		// TODO: add lib to skip throw
		throwList = append(throwList, fmt.Sprintf("%s:%s:%d", frame.Function, frame.File, frame.Line))
	}
	return throwList
}

type fake struct{}

var (
	goroot      = runtime.GOROOT()
	packageName = reflect.TypeOf(fake{}).PkgPath()
)
