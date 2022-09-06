package main

import (
	"note-API/app"
)

func main() {
	a := app.App{}
	a.InitializeBd("dummy")
	a.Run("localhost:8080")
}
