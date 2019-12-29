package main

import (
	"errors"
	"fmt"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// mapのシンタックスシュガー
type m map[string]interface{}

// handler→usecase→serviceという流れで実装
func handler(c echo.Context) error {
	err := usecase("hoge")
	if err != nil {
		return err
	}
	return c.JSON(200, m{"message": "success"})
}

func usecase(s string) error {
	err := service(s)
	if err != nil {
		return ErrSomething.Wrap(err)
	}
	return nil
}

func service(s string) error {
	return errors.New("occur error")
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.HTTPErrorHandler = ErrorHandler
	e.GET("/", handler)
	e.Logger.Fatal(e.Start(":1323"))
}

// エラーハンドリングの実装
func ErrorHandler(err error, c echo.Context) {
	var appErr *AppErr
	if errors.As(err, &appErr) {
		switch appErr.Level {
		case Fatal:
			fmt.Printf("[%s] %d %+v\n", appErr.Level, appErr.Code, appErr.Unwrap())
		case Error:
			fmt.Printf("[%s] %d %+v\n", appErr.Level, appErr.Code, appErr.Unwrap())
		case Warning:
		}
	} else {
		appErr = ErrUnknown
	}
	c.JSON(appErr.Code, m{"message": appErr.Message})
}

// ここから下はエラー型とエラー変数の定義
type AppErr struct {
	Level   ErrLevel
	Code    int
	Message string
	err     error
}

func (e *AppErr) Error() string {
	return fmt.Sprintf("[%s] %d: %+v", e.Level, e.Code, e.err)
}

type ErrLevel string

const (
	Fatal   ErrLevel = "FATAL"
	Error   ErrLevel = "ERROR"
	Warning ErrLevel = "WARNING"
)

var ErrSomething = &AppErr{
	Level:   Error,
	Code:    500,
	Message: "something error",
}

var ErrUnknown = &AppErr{
	Level:   Fatal,
	Code:    500,
	Message: "unknown error",
}

func (e *AppErr) Wrap(err error) error {
	e.err = err
	return e
}

func (e *AppErr) Unwrap() error {
	return e.err
}
