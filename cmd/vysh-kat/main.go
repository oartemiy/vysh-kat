package main

import (
	"vysh-kat/internal/app"
)

func main() {
	application := app.NewApp(app.NewContainer())
	application.Run()
}
