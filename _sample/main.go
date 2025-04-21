package main

import (
	"log/slog"
	"os"

	"github.com/go-stdx/try"
)

func main() {
	text()
	json()
}

func text() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
	defer try.Catch(nil)
	PanicTry()
}

func json() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	defer try.Catch(nil)
	PanicTry()
}

func FailfromPanic() {
	panic("Fail from Panic")
}

func FailfromTry() {
	try.E(Fail())
}

func Fail() (err error) {
	if true {
		return try.Wrapf("err from Fail()")
	} else {
		return nil
	}
}

func PanicTry() {
	try.Panic(Fail())
}
