package handler

import (
	"encoding/json"
	"net/http"
	db "note-API/database"
	"note-API/models"
	"note-API/utils"
	"regexp"
	"strconv"
)

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

func isAuth(r *http.Request) (bool, models.UserData) {
	userDate := models.UserData{
		ID:       1,
		LOGIN:    "1",
		PASSWORD: "1",
	}
	isOK := false
	userDate.LOGIN, userDate.PASSWORD, isOK = r.BasicAuth()
	if isOK != true {
		utils.InfoLog.Print("UNAUTHORIZED: user try get access without login")
		return false, userDate
	}
	isOK, userDate.ID = db.IsUser(userDate.LOGIN, userDate.PASSWORD)
	return isOK, userDate
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var requestUserDate models.UserData
	err := json.NewDecoder(r.Body).Decode(&requestUserDate)
	if err != nil {
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.BadRequestMessage())
		return
	}
	ok := false
	ok, err = db.AddUser(requestUserDate.LOGIN, requestUserDate.PASSWORD)
	if err != nil || ok == false {
		utils.InfoLog.Printf("NEW USER: try create user but already exist %s", requestUserDate.LOGIN)
		notesBytes, _ := json.MarshalIndent("Такой логин занят", "", "\t")
		utils.ReturnJsonResponse(w, http.StatusOK, notesBytes)
		return
	}
	utils.InfoLog.Printf("NEW USER: register user with login: %s", requestUserDate.LOGIN)
	utils.ReturnJsonResponse(w, http.StatusOK, utils.SuccessMessage())
	return
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var userDate models.UserData
	err := json.NewDecoder(r.Body).Decode(&userDate)
	if err != nil {
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.BadRequestMessage())
		return
	}
	isOK, userID := db.IsUser(userDate.LOGIN, userDate.PASSWORD)
	userDate.ID = userID
	if isOK {
		notesBytes, _ := json.MarshalIndent(userDate, "", "\t")
		utils.ReturnJsonResponse(w, http.StatusOK, notesBytes)
		utils.InfoLog.Printf("USER LOGIN: login: %s", userDate.LOGIN)
		return
	} else {
		utils.InfoLog.Printf("USER LOGIN: wrong password for login:  %s", userDate.LOGIN)
		notesBytes, _ := json.MarshalIndent("Неверный логин или пароль", "", "\t")
		utils.ReturnJsonResponse(w, http.StatusOK, notesBytes)
		return
	}
}

func GetNotes(w http.ResponseWriter, r *http.Request) {
	isAuth, user := isAuth(r)
	if isAuth != true {
		utils.ReturnJsonResponse(w, http.StatusUnauthorized, utils.UnauthorizedMessage())
		return
	}
	sortField := r.URL.Query().Get("sort")
	notes, err := db.GetNotes(sortField, user.ID)
	if err != nil {
		utils.InfoLog.Printf("GET NOTES: error: %s ", err)
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.ErrorMessage(err))
		return
	}
	notesBytes, _ := json.MarshalIndent(notes, "", "\t")
	utils.ReturnJsonResponse(w, http.StatusOK, notesBytes)
	return
}

func GetNote(w http.ResponseWriter, r *http.Request, getNoteRe *regexp.Regexp) {
	isAuth, user := isAuth(r)
	if isAuth != true {
		utils.ReturnJsonResponse(w, http.StatusUnauthorized, utils.UnauthorizedMessage())
		return
	}
	matches := getNoteRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		utils.InfoLog.Printf("GET NOTE: Bad query: %s", r.URL.Path)
		utils.ReturnJsonResponse(w, http.StatusNotFound, utils.NotFoundMessage())
		return
	}
	id := matches[1]
	if !db.IsNote(id) {
		utils.InfoLog.Printf("GET NOTE: Try GET nonexistent note: %s", id)
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.BadRequestMessage())
		return
	}
	note, _ := db.GetNote(id)
	var noteBytes []byte
	if note.ISPRIVATE == false {
		noteBytes, _ = json.MarshalIndent(note, "", "\t")
		utils.ReturnJsonResponse(w, http.StatusOK, noteBytes)
		return
	} else {
		if db.IsUserHaveAccess(user.ID, note.ID) {
			noteBytes, _ = json.MarshalIndent(note, "", "\t")
			utils.ReturnJsonResponse(w, http.StatusOK, noteBytes)
			return
		} else {
			utils.InfoLog.Printf("GET NOTE: User: %s try get access to note: %d without permission", user.LOGIN, note.ID)
			utils.ReturnJsonResponse(w, http.StatusForbidden, utils.ForbiddenMessage())
			return
		}
	}
}

func SetAccess(w http.ResponseWriter, r *http.Request, getNoteRe *regexp.Regexp) {
	isAuth, user := isAuth(r)
	if !isAuth {
		utils.ReturnJsonResponse(w, http.StatusUnauthorized, utils.UnauthorizedMessage())
		return
	}
	matches := getNoteRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		utils.InfoLog.Printf("SET ACCESS: Bad query: %s", r.URL.Path)
		utils.ReturnJsonResponse(w, http.StatusNotFound, utils.NotFoundMessage())
		return
	}
	noteID, err := strconv.Atoi(matches[1])
	if err != nil {
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.BadRequestMessage())
		return
	}
	if !db.IsUserNote(user.ID, noteID) {
		utils.InfoLog.Printf("SET ACCESS: User: %s try get access without permission.", user.LOGIN)
		utils.ReturnJsonResponse(w, http.StatusForbidden, utils.ForbiddenMessage())
		return
	}
	var accessDate models.AccessDate
	err = json.NewDecoder(r.Body).Decode(&accessDate)
	if err != nil {
		utils.InfoLog.Printf("SET ACCESS: Bad query: %s", matches)
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.BadRequestMessage())
		return
	}
	accessID := db.GetAccessID(accessDate.USERACCESSID, noteID)
	switch accessDate.MODE {
	case "ADD":
		if accessID != 0 {
			utils.InfoLog.Print("SET ACCESS: Try ADD existent note")
			utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.BadRequestMessage())
		} else {
			utils.InfoLog.Printf("SET ACCESS: User: %d  give access to note: %d for user: %d", user.ID, noteID, accessDate.USERACCESSID)
			_, err = db.AddAccess(accessDate.USERACCESSID, noteID)
			utils.ReturnJsonResponse(w, http.StatusOK, utils.SuccessMessage())
		}
	case "DELETE":
		if accessID == 0 {
			utils.InfoLog.Print("SET ACCESS: Try DELETE nonexistent note")
			utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.BadRequestMessage())
		} else {
			utils.InfoLog.Printf("SET ACCESS: user take away access to note: %d for user: %s", noteID, user.LOGIN)
			_, err = db.DeleteAccess(accessID)
			utils.ReturnJsonResponse(w, http.StatusOK, utils.SuccessMessage())
		}
	default:
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.BadRequestMessage())
	}
	return

}

func AddNote(w http.ResponseWriter, r *http.Request) {
	isAuth, user := isAuth(r)
	if isAuth != true {
		utils.ReturnJsonResponse(w, http.StatusUnauthorized, utils.UnauthorizedMessage())
		return
	}
	var note models.Note
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		utils.InfoLog.Printf("ADD NOTE: Bad JSON: %s", r.Body)
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.BadRequestMessage())
		return
	}
	var noteID int
	noteID, err = db.AddNote(note, user.ID)

	if err != nil || noteID == 0 {
		utils.InfoLog.Print("ADD NOTE: error")
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.BadRequestMessage())
		return
	}

	if note.EXPIRE > 0 {
		ch := make(chan bool)
		go db.DeleteNoteInTime(noteID, note.EXPIRE, ch)
	}
	utils.InfoLog.Printf("ADD NOTE: user: %s, added note: %d", user.LOGIN, noteID)
	utils.ReturnJsonResponse(w, http.StatusOK, utils.SuccessMessage())
	return
}

func UpdateNote(w http.ResponseWriter, r *http.Request) {
	isAuth, user := isAuth(r)
	if isAuth != true {
		utils.ReturnJsonResponse(w, http.StatusUnauthorized, utils.UnauthorizedMessage())
		return
	}
	var note models.Note
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		utils.InfoLog.Printf("UPDATE NOTE: Bad JSON: %d", r.Body)
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.BadRequestMessage())
		return
	}
	isUserNote := db.IsUserNote(user.ID, note.ID)
	if !isUserNote {
		utils.InfoLog.Printf("UPDATE NOTE: User: %s try get access without permission", user.LOGIN)
		utils.ReturnJsonResponse(w, http.StatusForbidden, utils.ForbiddenMessage())
		return
	}
	var id int
	id, err = db.UpdateNote(note, user.ID)
	if err != nil || id == 0 {
		utils.InfoLog.Printf("UPDATE NOTE: Update note error: %s", err)
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.BadRequestMessage())
		return
	}
	if note.EXPIRE > 0 {
		ch := make(chan bool)
		go db.DeleteNoteInTime(id, note.EXPIRE, ch)
	}
	utils.InfoLog.Printf("UPDATE NOTE: user: %s, updating note: %d", user.LOGIN, note.ID)
	utils.ReturnJsonResponse(w, http.StatusOK, utils.SuccessMessage())
	return
}

func DeleteNote(w http.ResponseWriter, r *http.Request, getNoteRe *regexp.Regexp) {
	isAuth, user := isAuth(r)
	if isAuth != true {
		utils.ReturnJsonResponse(w, http.StatusUnauthorized, utils.UnauthorizedMessage())
		return
	}
	matches := getNoteRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		utils.InfoLog.Printf("DELETE NOTE: Bad query: %s", matches)
		utils.ReturnJsonResponse(w, http.StatusNotFound, utils.NotFoundMessage())
		return
	}
	noteID, _ := strconv.Atoi(matches[1])
	isUserNote := db.IsUserNote(user.ID, noteID)
	if !isUserNote {
		utils.InfoLog.Printf("DELETE NOTE: User: %s try get access without permission", user.LOGIN)
		utils.ReturnJsonResponse(w, http.StatusForbidden, utils.ForbiddenMessage())
		return
	}
	utils.InfoLog.Printf("DELETE NOTE user: %s delete note: %d", user.LOGIN, noteID)
	_, _ = db.DeleteNote(noteID)
	utils.ReturnJsonResponse(w, http.StatusOK, utils.SuccessMessage())
	return
}
