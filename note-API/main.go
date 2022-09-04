package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"time"
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

const (
	dialect  = "postgres"
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "root"
	dbname   = "postgres"
)

type apiHandler struct{}

type note struct {
	ID        int    `json:"id"`
	BODY      string `json:"body"`
	TITLE     string `json:"title"`
	ISPRIVATE bool   `json:"isPrivate"`
	EXPIRE    int    `json:"expire"`
	USERID    int    `json:"userid"`
}

type userDate struct {
	ID       int    `json:"ID"`
	LOGIN    string `json:"login"`
	PASSWORD string `json:"password"`
}

type accessDate struct {
	USERACCESSID int    `json:"userAccessId"`
	MODE         string `json:"mode"`
}

func isItemExists(arrayType interface{}, item interface{}) bool {
	arr := reflect.ValueOf(arrayType)

	if arr.Kind() != reflect.Array {
		panic("Invalid data-type")
	}

	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true
		}
	}

	return false
}

func isAuth(r *http.Request) (bool, userDate) {
	var userDate userDate
	isOK := false
	userDate.LOGIN, userDate.PASSWORD, isOK = r.BasicAuth()
	if isOK != true {
		return false, userDate
	}
	isOK, userDate.ID = isUser(userDate.LOGIN, userDate.PASSWORD)
	return isOK, userDate
}

func isUser(l string, p string) (bool, int) {
	bd := OpenConnection()
	query := "SELECT ID FROM USERS WHERE LOGIN = $1 AND PASSWORD = $2"
	rows := bd.QueryRow(query, l, p)
	id := 0
	err := rows.Scan(&id)
	if err != nil {
		return false, 0
	}
	return true, id
}

func isUserHaveAccess(userID int, noteID int) bool {
	if isUserNote(userID, noteID) {
		return true
	}
	db := OpenConnection()
	query := "SELECT ID FROM ACCESS WHERE USERID = $1 AND NOTEID = $2"
	row := db.QueryRow(query, userID, noteID)
	id := 0
	err := row.Scan(&id)
	if err != nil || id == 0 {
		return false
	}
	return true
}

func isUserNote(userID int, noteID int) bool {
	var note note
	db := OpenConnection()
	query := "SELECT ID, USERID FROM NOTE WHERE ID = $1"
	row := db.QueryRow(query, noteID)
	err := row.Scan(&note.ID, &note.USERID)
	if err != nil || note.ID == 0 || note.USERID != userID {
		return false
	}
	return true
}

func deleteNoteInTime(noteID int, second int, ch chan bool) {
	time.Sleep(time.Duration(second) * time.Second)
	db := OpenConnection()
	query := "DELETE FROM NOTE WHERE ID = $1"
	_, _ = db.Exec(query, noteID)
	defer db.Close()
	ch <- true
}

func OpenConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	bd, err := sql.Open(dialect, psqlInfo)
	if err != nil {
		panic(err)
	}
	err = bd.Ping()
	if err != nil {
		panic(err)
	}

	return bd
}

func (h *apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	switch {
	case r.Method == http.MethodGet && getNotesRe.MatchString(r.URL.Path):
		h.GetNotes(w, r)
		return
	case r.Method == http.MethodGet && getNoteRe.MatchString(r.URL.Path):
		h.GetNote(w, r)
		return
	case r.Method == http.MethodPost && setAccessRe.MatchString(r.URL.Path):
		h.setAccess(w, r)
		return
	case r.Method == http.MethodPost && loginUserRe.MatchString(r.URL.Path):
		h.LoginUser(w, r)
		return
	case r.Method == http.MethodPost && addNoteRe.MatchString(r.URL.Path):
		h.AddNote(w, r)
		return
	case r.Method == http.MethodPut && updateNoteRe.MatchString(r.URL.Path):
		h.UpdateNote(w, r)
		return
	case r.Method == http.MethodDelete && deleteNoteRe.MatchString(r.URL.Path):
		h.DeleteNote(w, r)
		return
	case r.Method == http.MethodPost && registerUserRe.MatchString(r.URL.Path):
		h.RegisterUser(w, r)
		return
	case r.Method == http.MethodPost && loginUserRe.MatchString(r.URL.Path):
		h.LoginUser(w, r)
		return
	default:
		fmt.Fprintf(w, r.URL.Path)
		return
	}
}

func (h *apiHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var requestUserDate userDate
	err := json.NewDecoder(r.Body).Decode(&requestUserDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	bd := OpenConnection()

	query := "SELECT ID FROM USERS WHERE LOGIN = $1"
	var id = 0
	rows := bd.QueryRow(query, requestUserDate.LOGIN)
	err = rows.Scan(&id)
	if err == nil || id != 0 {
		notesBytes, _ := json.MarshalIndent("Такой логин занят", "", "\t")
		w.Write(notesBytes)
		defer bd.Close()
		return
	}
	query = "INSERT INTO USERS (login, password) VALUES ($1, $2)"
	_, err = bd.Exec(query, requestUserDate.LOGIN, requestUserDate.PASSWORD)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)
	defer bd.Close()
	return
}

func (h *apiHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var userDate userDate
	err := json.NewDecoder(r.Body).Decode(&userDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	isOK, userID := isUser(userDate.LOGIN, userDate.PASSWORD)
	userDate.ID = userID
	if isOK {
		notesBytes, _ := json.MarshalIndent(userDate, "", "\t")
		w.Write(notesBytes)
		return
	} else {
		notesBytes, _ := json.MarshalIndent("Неверный логин или пароль", "", "\t")
		w.Write(notesBytes)
		return
	}
}

func (h *apiHandler) GetNotes(w http.ResponseWriter, r *http.Request) {
	bd := OpenConnection()
	isAuth, user := isAuth(r)
	if isAuth != true {
		http.Error(w, "no auth", http.StatusUnauthorized)
		return
	}
	sortField := r.URL.Query().Get("sort")
	var query = ""

	if sortField == "" {
		query = "SELECT * FROM NOTE WHERE USERID = $1"
	} else {
		sortFields := [10]string{
			"-id",
			"id",
			"-title",
			"title",
			"-body",
			"body",
			"-expire",
			"expire",
			"isprivate",
			"-isprivate",
		}
		if isItemExists(sortFields, sortField) {
			query = "SELECT * FROM NOTE WHERE USERID = $1 ORDER BY " + sortField
		} else {
			panic("Неверный тип сортировки: " + sortField)
		}

	}
	rows, err := bd.Query(query, user.ID)
	if err != nil {
		panic(err)
	}
	var notes []note

	for rows.Next() {
		var note note
		err := rows.Scan(&note.ID, &note.BODY, &note.TITLE, &note.EXPIRE, &note.ISPRIVATE, &note.USERID)
		if err != nil {
			log.Fatal(err)
		}
		notes = append(notes, note)
	}

	notesBytes, _ := json.MarshalIndent(notes, "", "\t")
	w.Write(notesBytes)
	defer rows.Close()
	defer bd.Close()
	return
}

func (h *apiHandler) GetNote(w http.ResponseWriter, r *http.Request) {
	isAuth, user := isAuth(r)
	if isAuth != true {
		http.Error(w, "no auth", http.StatusUnauthorized)
		return
	}
	matches := getNoteRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		http.NotFound(w, r)
		return
	}
	var note note
	id := matches[1]
	bd := OpenConnection()
	const query = "SELECT * FROM NOTE WHERE ID = $1"
	row := bd.QueryRow(query, id)
	err := row.Scan(&note.ID, &note.BODY, &note.TITLE, &note.EXPIRE, &note.ISPRIVATE, &note.USERID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var noteBytes []byte
	if note.ISPRIVATE == false {
		noteBytes, _ = json.MarshalIndent(note, "", "\t")
	} else {
		if isUserHaveAccess(user.ID, note.ID) {
			noteBytes, _ = json.MarshalIndent(note, "", "\t")
		} else {
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}
	w.Write(noteBytes)
	defer bd.Close()
	return
}

func (h *apiHandler) setAccess(w http.ResponseWriter, r *http.Request) {
	isAuth, user := isAuth(r)
	if !isAuth {
		http.Error(w, "no auth", http.StatusUnauthorized)
		return
	}
	matches := getNoteRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		http.NotFound(w, r)
		return
	}
	noteID, err := strconv.Atoi(matches[1])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !isUserNote(user.ID, noteID) {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	var accessDate accessDate
	err = json.NewDecoder(r.Body).Decode(&accessDate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	db := OpenConnection()
	query := "SELECT ID FROM ACCESS WHERE USERID = $1 AND NOTEID = $2"
	row := db.QueryRow(query, accessDate.USERACCESSID, noteID)
	id := 0
	_ = row.Scan(&id)
	switch accessDate.MODE {
	case "ADD":
		if id != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		query = "INSERT INTO ACCESS (userid, noteid) VALUES ($1, $2)"
	case "DELETE":
		if id == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		query = "DELETE FROM ACCESS WHERE USERID = $1 AND NOTEID = $2"
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = db.Exec(query, accessDate.USERACCESSID, noteID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	return

}

func (h *apiHandler) deleteAccess(w http.ResponseWriter, r *http.Request) {
	isAuth, user := isAuth(r)
	if !isAuth {
		http.Error(w, "no auth", http.StatusUnauthorized)
		return
	}
	matches := getNoteRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		http.NotFound(w, r)
		return
	}
	noteID, err := strconv.Atoi(matches[1])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !isUserNote(user.ID, noteID) {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	var userAccessID accessDate
	err = json.NewDecoder(r.Body).Decode(&userAccessID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	db := OpenConnection()
	query := "SELECT ID FROM ACCESS WHERE USERID = $1 AND NOTEID = $2"
	row := db.QueryRow(query, userAccessID.USERACCESSID, noteID)
	id := 0
	_ = row.Scan(&id)
	if id != 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	query = "INSERT INTO ACCESS (userid, noteid) VALUES ($1, $2)"
	_, err = db.Exec(query, userAccessID.USERACCESSID, noteID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	return

}

func (h *apiHandler) AddNote(w http.ResponseWriter, r *http.Request) {
	isAuth, user := isAuth(r)
	if isAuth != true {
		http.Error(w, "no auth", http.StatusUnauthorized)
		return
	}
	var note note
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	db := OpenConnection()
	const query = `INSERT INTO NOTE (body, title, expire, isPrivate, userid) VALUES ($1, $2, $3, $4, $5) RETURNING ID`
	raw := db.QueryRow(query, note.BODY, note.TITLE, note.EXPIRE, note.ISPRIVATE, user.ID)
	noteID := 0
	err = raw.Scan(&noteID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if note.EXPIRE > 0 {
		ch := make(chan bool)
		go deleteNoteInTime(noteID, note.EXPIRE, ch)
	}
	w.WriteHeader(http.StatusOK)
	defer db.Close()
	return
}

func (h *apiHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	isAuth, user := isAuth(r)
	if isAuth != true {
		http.Error(w, "no auth", http.StatusUnauthorized)
		return
	}
	db := OpenConnection()
	var note note
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	const query = `UPDATE NOTE SET BODY = $2, TITLE = $3, EXPIRE = $4, ISPRIVATE = $5 WHERE ID = $1 AND USERID = $6 RETURNING ID`
	row := db.QueryRow(query, note.ID, note.BODY, note.TITLE, note.EXPIRE, note.ISPRIVATE, user.ID)
	id := 0
	err = row.Scan(&id)
	if err != nil || id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	defer db.Close()
	return
}

func (h *apiHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	isAuth, user := isAuth(r)
	if isAuth != true {
		http.Error(w, "no auth", http.StatusUnauthorized)
		return
	}
	matches := getNoteRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		http.NotFound(w, r)
		return
	}
	noteID := matches[1]
	bd := OpenConnection()
	const query = "DELETE FROM NOTE WHERE ID = $1 AND USERID = $2 RETURNING ID"
	row := bd.QueryRow(query, noteID, user.ID)
	id := 0
	err := row.Scan(&id)
	if err != nil || id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	noteBytes, _ := json.MarshalIndent(id, "", "\t")
	w.Write(noteBytes)
	defer bd.Close()
	return
}

// Endpoints

//GET API/notes           RESPONSE-body - {notes:
//<editor-fold desc="Description">
//											{
//												id: number,
//												body: string,
//												access: boolean,
//												expireIn: number(UNIX)
//											}[]
//										}
//</editor-fold>
//GET API/notes/{noteId}  RESPONSE-body - {notes:
//<editor-fold desc="Description">
//											{
//												id: number,
//												body: string,
//												access: boolean,
//												expireIn: number(UNIX)
//											}[]
//										}
//</editor-fold>
//POST API/notes          REQUEST-body - {
//<editor-fold desc="Description">
//											body: string,
//											access: boolean,
//											expireIn: number(UNIX)
//										}
//						  RESPONSE-body: {result: boolean}
//</editor-fold>
//PUT API/notes           REQUEST-body - {
//<editor-fold desc="Description">
//											id: number,
//											body: string,
//											access:boolean,
//											expireIn: number(UNIX)
//										 }
//						  RESPONSE-body: {result:boolean}
//</editor-fold>
//DELETE API/notes/{noteId}   RESPONSE-body - {result: boolean}
//POST API/register       REQUEST-body - {
//<editor-fold desc="Description">
//											login: string,
//											password: string
//										  }
//						  RESPONSE-body: {result: boolean}
//</editor-fold>
//POST API/login          RESPONSE-body -

func main() {
	mux := http.ServeMux{}
	mux.Handle("/api/", &apiHandler{})
	http.ListenAndServe("localhost:8080", &mux)
}
