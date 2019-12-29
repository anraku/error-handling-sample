package main

import (
	"errors"
	"fmt"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type m map[string]interface{}

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
	e.Use(ErrorHandler)
	e.GET("/", handler)
	e.Logger.Fatal(e.Start(":1323"))
}

func ErrorHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var appErr *AppErr
		if err := next(c); err != nil {
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
			return appErr.Unwrap()
		}
		return nil
	}
}

type AppErr struct {
	Level   ErrLevel
	Code    int
	Message string
	err     error
}

func (e *AppErr) Error() string {
	return fmt.Sprintf("[%s]: %d %s", e.Level, e.Code, e.Message)
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
