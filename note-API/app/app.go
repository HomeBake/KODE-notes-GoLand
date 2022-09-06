package app

import (
	"database/sql"
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
	utils.InfoLog.Printf("User trying connect to %s", r.URL)
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
		utils.InfoLog.Printf("User trying connect to nonexistent URL %s", r.URL)
		return
	}
}

func (a *App) InitializeBd(db string) {
	utils.InfoLog.Printf("Database initialize on %s mode", db)
	if db == "postgres" {
		err := os.Setenv("bdType", "postgres")
		err = os.Setenv("dialect", "postgres")
		err = os.Setenv("host", "localhost")
		err = os.Setenv("port", "5432")
		err = os.Setenv("user", "postgres")
		err = os.Setenv("password", "root")
		err = os.Setenv("dbname", "postgres")
		if err != nil {
			utils.ErrorLog.Printf("APP ERROR: %s", err)
		}
	} else {
		err := os.Setenv("bdType", "dummy")
		if err != nil {
			utils.ErrorLog.Printf("APP ERROR: %s", err)
		}
		database.FillDB()
	}
}

func (a *App) Run(addr string) {
	mux := http.ServeMux{}
	mux.HandleFunc("/api/", ServeHTTP)
	utils.InfoLog.Printf("Server listen %s address", addr)
	err := http.ListenAndServe(addr, &mux)
	if err != nil {
		utils.ErrorLog.Printf("Cannot listen address %s", addr)
		os.Exit(1)
	}
}
