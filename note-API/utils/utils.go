package utils

import (
	"net/http"
	"note-API/database"
	"note-API/models"
)

func ReturnJsonResponse(res http.ResponseWriter, httpCode int, resMessage []byte) {
	res.Header().Set("Content-type", "application/json")
	res.WriteHeader(httpCode)
	res.Write(resMessage)
}

func UnauthorizedMessage() []byte {
	HandlerMessage := []byte(`{
		"success": false,
   		"message": "You unauthorized ",
	}`)
	return HandlerMessage
}

func BadRequestMessage() []byte {
	HandlerMessage := []byte(`{
		"success": false,
   		"message": "Bad request",
	}`)
	return HandlerMessage
}

func SuccessMessage() []byte {
	HandlerMessage := []byte(`{
		"success": true,
   		"message": "Success",
	}`)
	return HandlerMessage
}

func ErrorMessage(err error) []byte {
	HandlerMessage := []byte(`{
		"success": false,
   		"message":` + err.Error() + `,
	}`)
	return HandlerMessage
}

func ForbiddenMessage() []byte {
	HandlerMessage := []byte(`{
		"success": false,
   		"message": "You dont have access",
	}`)
	return HandlerMessage
}

func NotFoundMessage() []byte {
	HandlerMessage := []byte(`{
		"success": false,
   		"message": "Not found",
	}`)
	return HandlerMessage
}

func FillDB() {
	note := models.Note{
		ID:        1,
		BODY:      "1",
		TITLE:     "2",
		ISPRIVATE: true,
		EXPIRE:    0,
		USERID:    1,
	}
	db.NoteDB[note.ID] = note
	note = models.Note{
		ID:        2,
		BODY:      "1",
		TITLE:     "2",
		ISPRIVATE: false,
		EXPIRE:    0,
		USERID:    1,
	}
	db.NoteDB[note.ID] = note
	note = models.Note{
		ID:        3,
		BODY:      "1",
		TITLE:     "2",
		ISPRIVATE: false,
		EXPIRE:    0,
		USERID:    1,
	}
	db.NoteDB[note.ID] = note
	note = models.Note{
		ID:        4,
		BODY:      "1",
		TITLE:     "2",
		ISPRIVATE: false,
		EXPIRE:    0,
		USERID:    1,
	}
	db.NoteDB[note.ID] = note
	note = models.Note{
		ID:        5,
		BODY:      "1",
		TITLE:     "2",
		ISPRIVATE: false,
		EXPIRE:    0,
		USERID:    1,
	}
	db.NoteDB[note.ID] = note
	note = models.Note{
		ID:        6,
		BODY:      "1",
		TITLE:     "2",
		ISPRIVATE: false,
		EXPIRE:    0,
		USERID:    1,
	}
	db.NoteDB[note.ID] = note
	user := models.UserData{
		ID:       1,
		LOGIN:    "1",
		PASSWORD: "1",
	}
	db.UserDB[user.ID] = user
	user = models.UserData{
		ID:       2,
		LOGIN:    "2",
		PASSWORD: "1",
	}
	db.UserDB[user.ID] = user
	user = models.UserData{
		ID:       3,
		LOGIN:    "3",
		PASSWORD: "1",
	}
	db.UserDB[user.ID] = user
	access := models.AccessDate{
		ID:           1,
		USERACCESSID: 2,
		NOTEID:       1,
	}
	db.AccessDateDB[access.ID] = access
}
