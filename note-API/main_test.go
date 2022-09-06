package main

import (
	"note-API/app"
	_ "note-API/database"
	"os"
	"testing"
)

var a app.App

func TestMain(m *testing.M) {
	a.InitializeBd("dummy")
	code := m.Run()
	os.Exit(code)
}
