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
	var userDate models.UserData
	isOK := false
	userDate.LOGIN, userDate.PASSWORD, isOK = r.BasicAuth()
	if isOK != true {
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
		notesBytes, _ := json.MarshalIndent("Такой логин занят", "", "\t")
		utils.ReturnJsonResponse(w, http.StatusOK, notesBytes)
		return
	}
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
		return
	} else {
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
		utils.ReturnJsonResponse(w, http.StatusUnauthorized, utils.ErrorMessage(err))
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
		utils.ReturnJsonResponse(w, http.StatusNotFound, utils.NotFoundMessage())
		return
	}
	id := matches[1]
	if !db.IsNote(id) {
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.BadRequestMessage())
		return
	}
	note, err := db.GetNote(id)
	if err != nil {
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.BadRequestMessage())
		return
	}
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
		utils.ReturnJsonResponse(w, http.StatusNotFound, utils.NotFoundMessage())
		return
	}
	noteID, err := strconv.Atoi(matches[1])
	if err != nil {
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.BadRequestMessage())
		return
	}
	if !db.IsUserNote(user.ID, noteID) {
		utils.ReturnJsonResponse(w, http.StatusForbidden, utils.ForbiddenMessage())
		return
	}
	var accessDate models.AccessDate
	err = json.NewDecoder(r.Body).Decode(&accessDate)
	if err != nil {
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.BadRequestMessage())
		return
	}
	accessID := db.GetAccessID(accessDate.USERACCESSID, noteID)
	switch accessDate.MODE {
	case "ADD":
		if accessID != 0 {
			utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.BadRequestMessage())
			return
		}
		_, err = db.AddAccess(accessDate.USERACCESSID, noteID)
	case "DELETE":
		if accessID == 0 {
			utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.BadRequestMessage())
			return
		}
		_, err = db.DeleteAccess(accessID)
	default:
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.BadRequestMessage())
		return
	}
	if err != nil {
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.BadRequestMessage())
		return
	}
	utils.ReturnJsonResponse(w, http.StatusOK, utils.SuccessMessage())
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
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.BadRequestMessage())
		return
	}
	var noteID int
	noteID, err = db.AddNote(note, user.ID)

	if err != nil {
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.BadRequestMessage())
		return
	}

	if note.EXPIRE > 0 {
		ch := make(chan bool)
		go db.DeleteNoteInTime(noteID, note.EXPIRE, ch)
	}
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
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.BadRequestMessage())
		return
	}
	isUserNote := db.IsUserNote(user.ID, note.ID)
	if !isUserNote {
		utils.ReturnJsonResponse(w, http.StatusForbidden, utils.ForbiddenMessage())
		return
	}
	var id int
	id, err = db.UpdateNote(note, user.ID)
	if err != nil || id == 0 {
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.BadRequestMessage())
		return
	}
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
		http.NotFound(w, r)
		return
	}
	noteID, _ := strconv.Atoi(matches[1])
	isUserNote := db.IsUserNote(user.ID, noteID)
	if !isUserNote {
		utils.ReturnJsonResponse(w, http.StatusForbidden, utils.ForbiddenMessage())
		return
	}
	id, err := db.DeleteNote(noteID)
	if err != nil || id == 0 {
		utils.ReturnJsonResponse(w, http.StatusBadRequest, utils.BadRequestMessage())
		return
	}
	utils.ReturnJsonResponse(w, http.StatusOK, utils.SuccessMessage())
	return
}
