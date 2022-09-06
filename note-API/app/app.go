package app

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"

	"net/http"
	"note-API/handler"
	utils "note-API/utils"
	"os"
	"regexp"
)

var (
	getNotesRe     = regexp.MustCompile(`/notes/?$`)
	getNoteRe      = regexp.MustCompile(`/notes/(\d+)*$`) // /notes/1
	setAccessRe    = getNoteRe
	addNoteRe      = getNotesRe
	updateNoteRe   = getNotesRe
	deleteNoteRe   = getNoteRe
	registerUserRe = regexp.MustCompile(`/user/register/?$`)
	loginUserRe    = regexp.MustCompile(`/user/login/?$`)
)

type App struct {
	DB *sql.DB
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch {
	case r.Method == http.MethodGet && getNotesRe.MatchString(r.URL.Path):
		handler.GetNotes(w, r)
		return
	case r.Method == http.MethodGet && getNoteRe.MatchString(r.URL.Path):
		handler.GetNote(w, r, getNoteRe)
		return
	case r.Method == http.MethodPost && setAccessRe.MatchString(r.URL.Path):
		handler.SetAccess(w, r, setAccessRe)
		return
	case r.Method == http.MethodPost && loginUserRe.MatchString(r.URL.Path):
		handler.LoginUser(w, r)
		return
	case r.Method == http.MethodPost && addNoteRe.MatchString(r.URL.Path):
		handler.AddNote(w, r)
		return
	case r.Method == http.MethodPut && updateNoteRe.MatchString(r.URL.Path):
		handler.UpdateNote(w, r)
		return
	case r.Method == http.MethodDelete && deleteNoteRe.MatchString(r.URL.Path):
		handler.DeleteNote(w, r, setAccessRe)
		return
	case r.Method == http.MethodPost && registerUserRe.MatchString(r.URL.Path):
		handler.RegisterUser(w, r)
		return
	case r.Method == http.MethodPost && loginUserRe.MatchString(r.URL.Path):
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
		utils.FillDB()
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
