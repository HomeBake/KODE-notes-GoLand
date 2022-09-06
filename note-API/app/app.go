package app

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"note-API/utils"

	"net/http"
	database "note-API/database"
	"note-API/handler"
	"os"
)

type App struct {
	DB *sql.DB
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch {
	case r.Method == http.MethodGet && utils.GetNotesRe.MatchString(r.URL.Path):
		handler.GetNotes(w, r)
		return
	case r.Method == http.MethodGet && utils.GetNoteRe.MatchString(r.URL.Path):
		handler.GetNote(w, r, utils.GetNoteRe)
		return
	case r.Method == http.MethodPost && utils.SetAccessRe.MatchString(r.URL.Path):
		handler.SetAccess(w, r, utils.SetAccessRe)
		return
	case r.Method == http.MethodPost && utils.LoginUserRe.MatchString(r.URL.Path):
		handler.LoginUser(w, r)
		return
	case r.Method == http.MethodPost && utils.AddNoteRe.MatchString(r.URL.Path):
		handler.AddNote(w, r)
		return
	case r.Method == http.MethodPut && utils.UpdateNoteRe.MatchString(r.URL.Path):
		handler.UpdateNote(w, r)
		return
	case r.Method == http.MethodDelete && utils.DeleteNoteRe.MatchString(r.URL.Path):
		handler.DeleteNote(w, r, utils.SetAccessRe)
		return
	case r.Method == http.MethodPost && utils.RegisterUserRe.MatchString(r.URL.Path):
		handler.RegisterUser(w, r)
		return
	case r.Method == http.MethodPost && utils.LoginUserRe.MatchString(r.URL.Path):
		handler.LoginUser(w, r)
		return
	default:
		fmt.Fprintf(w, r.URL.Path)
		return
	}
}

func (a *App) InitializeBd(db string) {
	if db == "postgres" {
		os.Setenv("bdType", "postgres")
		os.Setenv("dialect", "postgres")
		os.Setenv("host", "localhost")
		os.Setenv("port", "5432")
		os.Setenv("user", "postgres")
		os.Setenv("password", "root")
		os.Setenv("dbname", "postgres")
	} else {
		os.Setenv("bdType", "dummy")
		database.FillDB()
	}
}

func (a *App) Run(addr string) {
	mux := http.ServeMux{}
	mux.HandleFunc("/api/", ServeHTTP)
	err := http.ListenAndServe(addr, &mux)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
