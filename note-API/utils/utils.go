package utils

import (
	"net/http"
	"regexp"
)

var (
	GetNotesRe     = regexp.MustCompile(`/notes/?$`)
	GetNoteRe      = regexp.MustCompile(`/notes/(\d+)*$`) // /notes/1
	SetAccessRe    = GetNoteRe
	AddNoteRe      = GetNotesRe
	UpdateNoteRe   = GetNotesRe
	DeleteNoteRe   = GetNoteRe
	RegisterUserRe = regexp.MustCompile(`/user/register/?$`)
	LoginUserRe    = regexp.MustCompile(`/user/login/?$`)
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
