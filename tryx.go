package try

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"runtime"
)

const MaxDepth = 32

func Panic(v any) {
	var rerr error
	rerr, ok := v.(error)
	if !ok {
		rerr = fmt.Errorf("%v", v)
	}
	e(rerr)
}

func Wrap(err error) error {
	we := TryError{error: err}
	runtime.Callers(2, we.pc[:])
	return we
}

func Wrapf(format string, a ...any) error {
	we := TryError{error: fmt.Errorf(format, a...)}
	runtime.Callers(2, we.pc[:])
	return we
}

func Catch(fn func(err TryError)) {
	if fn == nil {
		fn = DefaultCatchHandler
	}

	r(recover(), func(w TryError) {
		fn(w)
	})
}

var SlogKey = "error"

var DefaultCatchHandler = func(err TryError) {
	slog.Error("try: panic "+err.Error(), SlogKey, err)
}

func (m TryError) MarshalJSON() ([]byte, error) {
	v := struct {
		Error      string   `json:"root"`
		Stacktrace []string `json:"stack"`
	}{
		Error:      m.Error(),
		Stacktrace: getStackTrace(m.pc[:]),
	}
	return json.Marshal(v)
}

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
