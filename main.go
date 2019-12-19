package main

import (
	"fmt"

	"golang.org/x/xerrors"

	"github.com/labstack/echo"
)

type m map[string]interface{}

func handler(c echo.Context) error {
	err := usecase("hoge")
	if err != nil {
		fmt.Printf("%+v\n", err)
		return c.JSON(500, m{"message": "App Errpr"})
	}
	return c.JSON(200, m{"message": "success"})
}

func usecase(s string) error {
	err := service(s)
	if err != nil {
		return xerrors.Errorf("[ERROR] usecase error: %w", err)
	}
	return nil
}

func service(s string) error {
	return xerrors.New("occur error")
}

func main() {
	e := echo.New()
	e.GET("/", handler)
	e.Logger.Fatal(e.Start(":1323"))
}

type AppError struct {
	Level   string
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}

var ErrSomething = &AppError{
	Level:   "Fatal",
	Code:    500,
	Message: "something error",
}
