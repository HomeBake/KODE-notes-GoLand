package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	db "note-API/database"
	"note-API/models"
	"note-API/utils"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	db.FillDB()
	code := m.Run()
	os.Exit(code)
}

func TestIsAuth(t *testing.T) {
	t.Run("user auth", func(t *testing.T) {
		appHost := "localhost:8080"
		url := fmt.Sprintf("http://%s/api/login/", appHost)
		jsonStr := []byte(`{}`)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		auth := []byte("1:1")
		basic := base64.StdEncoding.EncodeToString(auth)
		authHeader := fmt.Sprintf("Basic %s", basic)
		req.Header.Set("Authorization", authHeader)

		ok, _ := isAuth(req)
		if !ok {
			t.Errorf("Expected true. Got false")
		}
	})
	t.Run("user not auth", func(t *testing.T) {
		appHost := "localhost:8080"
		url := fmt.Sprintf("http://%s/api/login/", appHost)
		authHeader := ""
		jsonStr := []byte(`{}`)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Authorization", authHeader)

		ok, _ := isAuth(req)
		if ok {
			t.Errorf("Expected false. Got true")
		}
	})
}

func TestRegisterUser(t *testing.T) {
	t.Run("user register", func(t *testing.T) {
		appHost := "localhost:8080"
		url := fmt.Sprintf("http://%s/api/register/", appHost)
		jsonStr := []byte(`{"login" : "10" , "password" : "10"}`)
		r, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		RegisterUser(w, r)
		expectedStatus := 200
		expectedJSON := utils.SuccessMessage()

		if expectedStatus != w.Code {
			t.Errorf("Invalid code %d expected %d", w.Code, expectedStatus)
		}
		if !bytes.Equal(expectedJSON, w.Body.Bytes()) {
			t.Errorf("Invalid JSONbody %s expected %s", w.Body.Bytes(), expectedJSON)
		}
	})
	t.Run("user already register", func(t *testing.T) {
		appHost := "localhost:8080"
		url := fmt.Sprintf("http://%s/api/register/", appHost)
		jsonStr := []byte(`{"login" : "1" , "password" : "1"}`)
		r, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		RegisterUser(w, r)
		expectedStatus := 200
		expectedJSON := []byte(`"Такой логин занят"`)

		if expectedStatus != w.Code {
			t.Errorf("Invalid code %d expected %d", w.Code, expectedStatus)
		}
		if !bytes.Equal(expectedJSON, w.Body.Bytes()) {
			t.Errorf("Invalid JSONbody %s expected %s", w.Body.Bytes(), expectedJSON)
		}
	})
	t.Run("Bad json", func(t *testing.T) {
		appHost := "localhost:8080"
		url := fmt.Sprintf("http://%s/api/register/", appHost)
		jsonStr := []byte(`{login : "1" , password : "1"}`)
		r, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		RegisterUser(w, r)
		expectedStatus := 400
		expectedJSON := utils.BadRequestMessage()

		if expectedStatus != w.Code {
			t.Errorf("Invalid code %d expected %d", w.Code, expectedStatus)
		}
		if !bytes.Equal(expectedJSON, w.Body.Bytes()) {
			t.Errorf("Invalid JSONbody %s expected %s", w.Body.Bytes(), expectedJSON)
		}
	})

}

func TestLoginUser(t *testing.T) {
	t.Run("user login success", func(t *testing.T) {
		appHost := "localhost:8080"
		url := fmt.Sprintf("http://%s/api/login/", appHost)
		jsonStr := []byte(`{"login" : "1" , "password" : "1"}`)
		r, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		LoginUser(w, r)
		expectedStatus := 200
		expectedUser := models.UserData{
			ID:       1,
			LOGIN:    "1",
			PASSWORD: "1",
		}
		expectedJSON, _ := json.MarshalIndent(expectedUser, "", "\t")
		if expectedStatus != w.Code {
			t.Errorf("Invalid code %d expected %d", w.Code, expectedStatus)
		}
		if !bytes.Equal(expectedJSON, w.Body.Bytes()) {
			t.Errorf("Invalid JSONbody %s expected %s", w.Body.Bytes(), expectedJSON)
		}
	})
	t.Run("user bad data", func(t *testing.T) {
		appHost := "localhost:8080"
		url := fmt.Sprintf("http://%s/api/login/", appHost)
		jsonStr := []byte(`{"login" : "bad" , "password" : "bad"}`)
		r, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		LoginUser(w, r)
		expectedStatus := 200
		expectedJSON := []byte(`"Неверный логин или пароль"`)
		if expectedStatus != w.Code {
			t.Errorf("Invalid code %d expected %d", w.Code, expectedStatus)
		}
		if !bytes.Equal(expectedJSON, w.Body.Bytes()) {
			t.Errorf("Invalid JSONbody %s expected %s", w.Body.Bytes(), expectedJSON)
		}
	})
	t.Run("Bad json", func(t *testing.T) {
		appHost := "localhost:8080"
		url := fmt.Sprintf("http://%s/api/login/", appHost)
		jsonStr := []byte(`{login : "1" , password : "1"}`)
		r, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		LoginUser(w, r)
		expectedStatus := 400
		expectedJSON := utils.BadRequestMessage()

		if expectedStatus != w.Code {
			t.Errorf("Invalid code %d expected %d", w.Code, expectedStatus)
		}
		if !bytes.Equal(expectedJSON, w.Body.Bytes()) {
			t.Errorf("Invalid JSONbody %s expected %s", w.Body.Bytes(), expectedJSON)
		}
	})
}
func TestGetNotes(t *testing.T) {
	t.Run("Notes got success without sort field", func(t *testing.T) {
		appHost := "localhost:8080"
		url := fmt.Sprintf("http://%s/api/notes/", appHost)
		jsonStr := []byte(``)
		r, _ := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
		auth := []byte("2:1")
		basic := base64.StdEncoding.EncodeToString(auth)
		authHeader := fmt.Sprintf("Basic %s", basic)
		r.Header.Set("Authorization", authHeader)
		w := httptest.NewRecorder()
		GetNotes(w, r)
		expectedStatus := 200
		expectedNotes := []models.Note{
			{ID: 4,
				BODY:      "1",
				TITLE:     "2",
				ISPRIVATE: false,
				EXPIRE:    0,
				USERID:    2},
		}
		expectedJSON, _ := json.MarshalIndent(expectedNotes, "", "\t")
		if expectedStatus != w.Code {
			t.Errorf("Invalid code %d expected %d", w.Code, expectedStatus)
		}
		if !bytes.Equal(expectedJSON, w.Body.Bytes()) {
			t.Errorf("Invalid JSONbody %s expected %s", w.Body.Bytes(), expectedJSON)
		}
	})
	t.Run("Notes got success with sort field", func(t *testing.T) {
		appHost := "localhost:8080"
		url := fmt.Sprintf("http://%s/api/notes?sort=id", appHost)
		jsonStr := []byte(``)
		r, _ := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
		auth := []byte("1:1")
		basic := base64.StdEncoding.EncodeToString(auth)
		authHeader := fmt.Sprintf("Basic %s", basic)
		r.Header.Set("Authorization", authHeader)
		w := httptest.NewRecorder()
		GetNotes(w, r)
		expectedStatus := 200
		expectedNotes := []models.Note{
			{
				ID:        1,
				BODY:      "1",
				TITLE:     "2",
				ISPRIVATE: true,
				EXPIRE:    0,
				USERID:    1,
			},
			{
				ID:        2,
				BODY:      "1",
				TITLE:     "2",
				ISPRIVATE: true,
				EXPIRE:    0,
				USERID:    1,
			},
			{
				ID:        3,
				BODY:      "1",
				TITLE:     "2",
				ISPRIVATE: false,
				EXPIRE:    0,
				USERID:    1,
			},
		}
		expectedJSON, _ := json.MarshalIndent(expectedNotes, "", "\t")
		if expectedStatus != w.Code {
			t.Errorf("Invalid code %d expected %d", w.Code, expectedStatus)
		}
		if !bytes.Equal(expectedJSON, w.Body.Bytes()) {
			t.Errorf("Invalid JSONbody %s expected %s", w.Body.Bytes(), expectedJSON)
		}
	})
	t.Run("Notes got unSuccess with bad sort field", func(t *testing.T) {
		appHost := "localhost:8080"
		url := fmt.Sprintf("http://%s/api/notes?sort=bad", appHost)
		jsonStr := []byte(``)
		r, _ := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
		auth := []byte("1:1")
		basic := base64.StdEncoding.EncodeToString(auth)
		authHeader := fmt.Sprintf("Basic %s", basic)
		r.Header.Set("Authorization", authHeader)
		w := httptest.NewRecorder()
		GetNotes(w, r)
		expectedStatus := 400
		err := errors.New("sort field is not exist")
		expectedMessage := utils.ErrorMessage(err)
		if expectedStatus != w.Code {
			t.Errorf("Invalid code %d expected %d", w.Code, expectedStatus)
		}
		if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
			t.Errorf("Invalid Body %s expected %s", w.Body.Bytes(), expectedMessage)
		}
	})
	t.Run("Notes got unSuccess with no auth", func(t *testing.T) {
		appHost := "localhost:8080"
		url := fmt.Sprintf("http://%s/api/notes?sort=bad", appHost)
		jsonStr := []byte(``)
		r, _ := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		GetNotes(w, r)
		expectedStatus := 401
		expectedMessage := utils.UnauthorizedMessage()
		if expectedStatus != w.Code {
			t.Errorf("Invalid code %d expected %d", w.Code, expectedStatus)
		}
		if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
			t.Errorf("Invalid Body %s expected %s", w.Body.Bytes(), expectedMessage)
		}
	})
}

func TestGetNote(t *testing.T) {
	t.Run("User get own note", func(t *testing.T) {
		appHost := "localhost:8080"
		url := fmt.Sprintf("http://%s/api/notes/1", appHost)
		jsonStr := []byte(``)
		r, _ := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
		auth := []byte("1:1")
		basic := base64.StdEncoding.EncodeToString(auth)
		authHeader := fmt.Sprintf("Basic %s", basic)
		r.Header.Set("Authorization", authHeader)
		w := httptest.NewRecorder()
		GetNote(w, r, utils.GetNoteRe)
		expectedStatus := 200
		expectedNote := models.Note{
			ID:        1,
			BODY:      "1",
			TITLE:     "2",
			ISPRIVATE: true,
			EXPIRE:    0,
			USERID:    1,
		}
		expectedJSON, _ := json.MarshalIndent(expectedNote, "", "\t")
		if expectedStatus != w.Code {
			t.Errorf("Invalid code %d expected %d", w.Code, expectedStatus)
		}
		if !bytes.Equal(expectedJSON, w.Body.Bytes()) {
			t.Errorf("Invalid JSONbody %s expected %s", w.Body.Bytes(), expectedJSON)
		}
	})
	t.Run("User get not private note", func(t *testing.T) {
		appHost := "localhost:8080"
		url := fmt.Sprintf("http://%s/api/notes/3", appHost)
		jsonStr := []byte(``)
		r, _ := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
		auth := []byte("2:1")
		basic := base64.StdEncoding.EncodeToString(auth)
		authHeader := fmt.Sprintf("Basic %s", basic)
		r.Header.Set("Authorization", authHeader)
		w := httptest.NewRecorder()
		GetNote(w, r, utils.GetNoteRe)
		expectedStatus := 200
		expectedNote := models.Note{
			ID:        3,
			BODY:      "1",
			TITLE:     "2",
			ISPRIVATE: false,
			EXPIRE:    0,
			USERID:    1,
		}
		expectedJSON, _ := json.MarshalIndent(expectedNote, "", "\t")
		if expectedStatus != w.Code {
			t.Errorf("Invalid code %d expected %d", w.Code, expectedStatus)
		}
		if !bytes.Equal(expectedJSON, w.Body.Bytes()) {
			t.Errorf("Invalid JSONbody %s expected %s", w.Body.Bytes(), expectedJSON)
		}
	})
	t.Run("User get private note with access", func(t *testing.T) {
		appHost := "localhost:8080"
		url := fmt.Sprintf("http://%s/api/notes/1", appHost)
		jsonStr := []byte(``)
		r, _ := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
		auth := []byte("2:1")
		basic := base64.StdEncoding.EncodeToString(auth)
		authHeader := fmt.Sprintf("Basic %s", basic)
		r.Header.Set("Authorization", authHeader)
		w := httptest.NewRecorder()
		GetNote(w, r, utils.GetNoteRe)
		expectedStatus := 200
		expectedNote := models.Note{
			ID:        1,
			BODY:      "1",
			TITLE:     "2",
			ISPRIVATE: true,
			EXPIRE:    0,
			USERID:    1,
		}
		expectedJSON, _ := json.MarshalIndent(expectedNote, "", "\t")
		if expectedStatus != w.Code {
			t.Errorf("Invalid code %d expected %d", w.Code, expectedStatus)
		}
		if !bytes.Equal(expectedJSON, w.Body.Bytes()) {
			t.Errorf("Invalid JSONbody %s expected %s", w.Body.Bytes(), expectedJSON)
		}
	})
	t.Run("User try get private note without access", func(t *testing.T) {
		appHost := "localhost:8080"
		url := fmt.Sprintf("http://%s/api/notes/2", appHost)
		jsonStr := []byte(``)
		r, _ := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
		auth := []byte("2:1")
		basic := base64.StdEncoding.EncodeToString(auth)
		authHeader := fmt.Sprintf("Basic %s", basic)
		r.Header.Set("Authorization", authHeader)
		w := httptest.NewRecorder()
		GetNote(w, r, utils.GetNoteRe)
		expectedStatus := 403
		expectedMessage := utils.ForbiddenMessage()
		if expectedStatus != w.Code {
			t.Errorf("Invalid code %d expected %d", w.Code, expectedStatus)
		}
		if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
			t.Errorf("Invalid body %s expected %s", w.Body.Bytes(), expectedMessage)
		}
	})
	t.Run("User try get not exist note", func(t *testing.T) {
		appHost := "localhost:8080"
		url := fmt.Sprintf("http://%s/api/notes/22", appHost)
		jsonStr := []byte(``)
		r, _ := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
		auth := []byte("2:1")
		basic := base64.StdEncoding.EncodeToString(auth)
		authHeader := fmt.Sprintf("Basic %s", basic)
		r.Header.Set("Authorization", authHeader)
		w := httptest.NewRecorder()
		GetNote(w, r, utils.GetNoteRe)
		expectedStatus := 400
		expectedMessage := utils.BadRequestMessage()
		if expectedStatus != w.Code {
			t.Errorf("Invalid code %d expected %d", w.Code, expectedStatus)
		}
		if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
			t.Errorf("Invalid body %s expected %s", w.Body.Bytes(), expectedMessage)
		}
	})
	t.Run("Bad link", func(t *testing.T) {
		appHost := "localhost:8080"
		url := fmt.Sprintf("http://%s/api/notes/h", appHost)
		jsonStr := []byte(``)
		r, _ := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
		auth := []byte("2:1")
		basic := base64.StdEncoding.EncodeToString(auth)
		authHeader := fmt.Sprintf("Basic %s", basic)
		r.Header.Set("Authorization", authHeader)
		w := httptest.NewRecorder()
		GetNote(w, r, utils.GetNoteRe)
		expectedStatus := 404
		expectedMessage := utils.NotFoundMessage()
		if expectedStatus != w.Code {
			t.Errorf("Invalid code %d expected %d", w.Code, expectedStatus)
		}
		if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
			t.Errorf("Invalid body %s expected %s", w.Body.Bytes(), expectedMessage)
		}
	})
	t.Run("User unauthorized", func(t *testing.T) {
		appHost := "localhost:8080"
		url := fmt.Sprintf("http://%s/api/notes/h", appHost)
		jsonStr := []byte(``)
		r, _ := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
		w := httptest.NewRecorder()
		GetNote(w, r, utils.GetNoteRe)
		expectedStatus := 401
		expectedMessage := utils.UnauthorizedMessage()
		if expectedStatus != w.Code {
			t.Errorf("Invalid code %d expected %d", w.Code, expectedStatus)
		}
		if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
			t.Errorf("Invalid body %s expected %s", w.Body.Bytes(), expectedMessage)
		}
	})
}

func TestSetAccess(t *testing.T) {
	appHost := "localhost:8080"
	auth := []byte("1:1")
	basic := base64.StdEncoding.EncodeToString(auth)
	authHeader := fmt.Sprintf("Basic %s", basic)
	t.Run("User add access success", func(t *testing.T) {
		url := fmt.Sprintf("http://%s/api/notes/2", appHost)
		w := httptest.NewRecorder()
		access := models.AccessDate{
			USERACCESSID: 2,
			MODE:         "ADD",
		}
		jsonStr, _ := json.MarshalIndent(access, "", "\t")
		r, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		r.Header.Set("Authorization", authHeader)
		SetAccess(w, r, utils.GetNoteRe)
		expectedStatus := 200
		expectedMessage := utils.SuccessMessage()
		if expectedStatus != w.Code {
			t.Errorf("Invalid code %d expected %d", w.Code, expectedStatus)
		}
		if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
			t.Errorf("Invalid JSONbody %s expected %s", w.Body.Bytes(), expectedMessage)
		}
	})
	t.Run("User add exist access", func(t *testing.T) {
		url := fmt.Sprintf("http://%s/api/notes/2", appHost)
		w := httptest.NewRecorder()
		access := models.AccessDate{
			USERACCESSID: 2,
			MODE:         "ADD",
		}
		jsonStr, _ := json.MarshalIndent(access, "", "\t")
		r, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		r.Header.Set("Authorization", authHeader)
		SetAccess(w, r, utils.GetNoteRe)
		expectedStatus := 400
		expectedMessage := utils.BadRequestMessage()
		if expectedStatus != w.Code {
			t.Errorf("Invalid code %d expected %d", w.Code, expectedStatus)
		}
		if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
			t.Errorf("Invalid JSONbody %s expected %s", w.Body.Bytes(), expectedMessage)
		}
	})
	t.Run("User delete access success", func(t *testing.T) {
		url := fmt.Sprintf("http://%s/api/notes/1", appHost)
		w := httptest.NewRecorder()
		access := models.AccessDate{
			USERACCESSID: 2,
			MODE:         "DELETE",
		}
		jsonStr, _ := json.MarshalIndent(access, "", "\t")
		r, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		r.Header.Set("Authorization", authHeader)
		SetAccess(w, r, utils.GetNoteRe)
		expectedStatus := 200
		expectedMessage := utils.SuccessMessage()
		if expectedStatus != w.Code {
			t.Errorf("Invalid code %d expected %d", w.Code, expectedStatus)
		}
		if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
			t.Errorf("Invalid JSONbody %s expected %s", w.Body.Bytes(), expectedMessage)
		}
	})
	t.Run("User delete non-exist access", func(t *testing.T) {
		url := fmt.Sprintf("http://%s/api/notes/2", appHost)
		w := httptest.NewRecorder()

		access := models.AccessDate{
			USERACCESSID: 3,
			MODE:         "DELETE",
		}
		jsonStr, _ := json.MarshalIndent(access, "", "\t")
		r, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		r.Header.Set("Authorization", authHeader)
		SetAccess(w, r, utils.GetNoteRe)
		expectedStatus := 400
		expectedMessage := utils.BadRequestMessage()
		if expectedStatus != w.Code {
			t.Errorf("Invalid code %d expected %d", w.Code, expectedStatus)
		}
		if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
			t.Errorf("Invalid JSONbody %s expected %s", w.Body.Bytes(), expectedMessage)
		}
	})
	t.Run("User bad mode ", func(t *testing.T) {
		url := fmt.Sprintf("http://%s/api/notes/2", appHost)
		w := httptest.NewRecorder()
		access := models.AccessDate{
			USERACCESSID: 2,
			MODE:         "BAD",
		}
		jsonStr, _ := json.MarshalIndent(access, "", "\t")
		r, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		r.Header.Set("Authorization", authHeader)
		SetAccess(w, r, utils.GetNoteRe)
		expectedStatus := 400
		expectedMessage := utils.BadRequestMessage()
		if expectedStatus != w.Code {
			t.Errorf("Invalid code %d expected %d", w.Code, expectedStatus)
		}
		if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
			t.Errorf("Invalid JSONbody %s expected %s", w.Body.Bytes(), expectedMessage)
		}
	})
	t.Run("User bad JSON ", func(t *testing.T) {
		url := fmt.Sprintf("http://%s/api/notes/2", appHost)
		w := httptest.NewRecorder()
		access := []byte(`{ bad : bad}`)
		r, _ := http.NewRequest("POST", url, bytes.NewBuffer(access))
		r.Header.Set("Authorization", authHeader)
		SetAccess(w, r, utils.GetNoteRe)
		expectedStatus := 400
		expectedMessage := utils.BadRequestMessage()
		if expectedStatus != w.Code {
			t.Errorf("Invalid code %d expected %d", w.Code, expectedStatus)
		}
		if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
			t.Errorf("Invalid JSONbody %s expected %s", w.Body.Bytes(), expectedMessage)
		}
	})

	t.Run("User doesnt have access ", func(t *testing.T) {
		url := fmt.Sprintf("http://%s/api/notes/6", appHost)
		w := httptest.NewRecorder()
		access := models.AccessDate{
			USERACCESSID: 2,
			MODE:         "ADD",
		}
		jsonStr, _ := json.MarshalIndent(access, "", "\t")
		r, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		r.Header.Set("Authorization", authHeader)
		SetAccess(w, r, utils.GetNoteRe)
		expectedStatus := 403
		expectedMessage := utils.ForbiddenMessage()
		if expectedStatus != w.Code {
			t.Errorf("Invalid code %d expected %d", w.Code, expectedStatus)
		}
		if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
			t.Errorf("Invalid JSONbody %s expected %s", w.Body.Bytes(), expectedMessage)
		}
	})

	t.Run("Bad note ID ", func(t *testing.T) {
		url := fmt.Sprintf("http://%s/api/notes/", appHost)
		w := httptest.NewRecorder()
		access := models.AccessDate{
			USERACCESSID: 2,
			MODE:         "ADD",
		}
		jsonStr, _ := json.MarshalIndent(access, "", "\t")
		r, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		r.Header.Set("Authorization", authHeader)
		SetAccess(w, r, utils.GetNoteRe)
		expectedStatus := 400
		expectedMessage := utils.BadRequestMessage()
		if expectedStatus != w.Code {
			t.Errorf("Invalid code %d expected %d", w.Code, expectedStatus)
		}
		if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
			t.Errorf("Invalid JSONbody %s expected %s", w.Body.Bytes(), expectedMessage)
		}
	})
	t.Run("Bad query ", func(t *testing.T) {
		url := fmt.Sprintf("http://%s/api/notes/dsa34", appHost)
		w := httptest.NewRecorder()
		access := models.AccessDate{
			USERACCESSID: 2,
			MODE:         "ADD",
		}
		jsonStr, _ := json.MarshalIndent(access, "", "\t")
		r, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		r.Header.Set("Authorization", authHeader)
		SetAccess(w, r, utils.GetNoteRe)
		expectedStatus := 404
		expectedMessage := utils.NotFoundMessage()
		if expectedStatus != w.Code {
			t.Errorf("Invalid code %d expected %d", w.Code, expectedStatus)
		}
		if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
			t.Errorf("Invalid JSONbody %s expected %s", w.Body.Bytes(), expectedMessage)
		}
	})
	t.Run("Not auth ", func(t *testing.T) {
		url := fmt.Sprintf("http://%s/api/notes/dsa34", appHost)
		auth := []byte("bad:bad")
		basic := base64.StdEncoding.EncodeToString(auth)
		authHeader := fmt.Sprintf("Basic %s", basic)
		w := httptest.NewRecorder()
		access := models.AccessDate{
			USERACCESSID: 2,
			MODE:         "ADD",
		}
		jsonStr, _ := json.MarshalIndent(access, "", "\t")
		r, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		r.Header.Set("Authorization", authHeader)
		SetAccess(w, r, utils.GetNoteRe)
		expectedStatus := 401
		expectedMessage := utils.UnauthorizedMessage()
		if expectedStatus != w.Code {
			t.Errorf("Invalid code %d expected %d", w.Code, expectedStatus)
		}
		if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
			t.Errorf("Invalid JSONbody %s expected %s", w.Body.Bytes(), expectedMessage)
		}
	})
}

func TestAddNote(t *testing.T) {
	t.Run("User add note successfully", func(t *testing.T) {
		appHost := "localhost:8080"
		url := fmt.Sprintf("http://%s/api/notes", appHost)
		note := models.Note{
			BODY:      "new",
			TITLE:     "new",
			ISPRIVATE: false,
			EXPIRE:    0,
		}
		jsonStr, _ := json.MarshalIndent(note, "", "\t")
		r, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		auth := []byte("1:1")
		basic := base64.StdEncoding.EncodeToString(auth)
		authHeader := fmt.Sprintf("Basic %s", basic)
		r.Header.Set("Authorization", authHeader)
		w := httptest.NewRecorder()
		AddNote(w, r)
		expectedStatus := 200
		expectedMessage := utils.SuccessMessage()
		if expectedStatus != w.Code {
			t.Errorf("Invalid code %d expected %d", w.Code, expectedStatus)
		}
		if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
			t.Errorf("Invalid JSONbody %s expected %s", w.Body.Bytes(), expectedMessage)
		}
	})
	t.Run("User add note successfully with expire time", func(t *testing.T) {
		appHost := "localhost:8080"
		url := fmt.Sprintf("http://%s/api/notes", appHost)
		note := models.Note{
			BODY:      "new",
			TITLE:     "new",
			ISPRIVATE: false,
			EXPIRE:    1,
		}
		jsonStr, _ := json.MarshalIndent(note, "", "\t")
		r, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		auth := []byte("1:1")
		basic := base64.StdEncoding.EncodeToString(auth)
		authHeader := fmt.Sprintf("Basic %s", basic)
		r.Header.Set("Authorization", authHeader)
		w := httptest.NewRecorder()
		AddNote(w, r)
		expectedStatus := 200
		expectedMessage := utils.SuccessMessage()
		if expectedStatus != w.Code {
			t.Errorf("Invalid code %d expected %d", w.Code, expectedStatus)
		}
		if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
			t.Errorf("Invalid JSONbody %s expected %s", w.Body.Bytes(), expectedMessage)
		}
	})
	t.Run("Bad body", func(t *testing.T) {
		appHost := "localhost:8080"
		url := fmt.Sprintf("http://%s/api/notes", appHost)
		note := []byte(`{bad : bad'}`)
		r, _ := http.NewRequest("POST", url, bytes.NewBuffer(note))
		auth := []byte("1:1")
		basic := base64.StdEncoding.EncodeToString(auth)
		authHeader := fmt.Sprintf("Basic %s", basic)
		r.Header.Set("Authorization", authHeader)
		w := httptest.NewRecorder()
		AddNote(w, r)
		expectedStatus := 400
		expectedMessage := utils.BadRequestMessage()
		if expectedStatus != w.Code {
			t.Errorf("Invalid code %d expected %d", w.Code, expectedStatus)
		}
		if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
			t.Errorf("Invalid JSONbody %s expected %s", w.Body.Bytes(), expectedMessage)
		}
	})
	t.Run("User unauthorized", func(t *testing.T) {
		appHost := "localhost:8080"
		url := fmt.Sprintf("http://%s/api/notes", appHost)
		note := []byte(`{bad : bad'}`)
		r, _ := http.NewRequest("POST", url, bytes.NewBuffer(note))
		auth := []byte("bad:bad")
		basic := base64.StdEncoding.EncodeToString(auth)
		authHeader := fmt.Sprintf("Basic %s", basic)
		r.Header.Set("Authorization", authHeader)
		w := httptest.NewRecorder()
		AddNote(w, r)
		expectedStatus := 401
		expectedMessage := utils.UnauthorizedMessage()
		if expectedStatus != w.Code {
			t.Errorf("Invalid code %d expected %d", w.Code, expectedStatus)
		}
		if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
			t.Errorf("Invalid JSONbody %s expected %s", w.Body.Bytes(), expectedMessage)
		}
	})
}

func TestUpdateNote(t *testing.T) {
	t.Run("User update note successfully", func(t *testing.T) {
		appHost := "localhost:8080"
		url := fmt.Sprintf("http://%s/api/notes/1", appHost)
		note := models.Note{
			ID:        1,
			BODY:      "update",
			TITLE:     "update",
			ISPRIVATE: false,
			EXPIRE:    0,
		}
		jsonStr, _ := json.MarshalIndent(note, "", "\t")
		r, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		auth := []byte("1:1")
		basic := base64.StdEncoding.EncodeToString(auth)
		authHeader := fmt.Sprintf("Basic %s", basic)
		r.Header.Set("Authorization", authHeader)
		w := httptest.NewRecorder()
		UpdateNote(w, r)
		expectedStatus := 200
		expectedMessage := utils.SuccessMessage()
		if expectedStatus != w.Code {
			t.Errorf("Invalid code %d expected %d", w.Code, expectedStatus)
		}
		if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
			t.Errorf("Invalid JSONbody %s expected %s", w.Body.Bytes(), expectedMessage)
		}
	})
	t.Run("Delete note after success update", func(t *testing.T) {
		appHost := "localhost:8080"
		url := fmt.Sprintf("http://%s/api/notes/1", appHost)
		note := models.Note{
			ID:        1,
			BODY:      "update",
			TITLE:     "update",
			ISPRIVATE: false,
			EXPIRE:    1,
		}
		jsonStr, _ := json.MarshalIndent(note, "", "\t")
		r, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		auth := []byte("1:1")
		basic := base64.StdEncoding.EncodeToString(auth)
		authHeader := fmt.Sprintf("Basic %s", basic)
		r.Header.Set("Authorization", authHeader)
		w := httptest.NewRecorder()
		UpdateNote(w, r)
		expectedStatus := 200
		expectedMessage := utils.SuccessMessage()
		if expectedStatus != w.Code {
			t.Errorf("Invalid code %d expected %d", w.Code, expectedStatus)
		}
		if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
			t.Errorf("Invalid JSONbody %s expected %s", w.Body.Bytes(), expectedMessage)
		}
	})
	t.Run("User try update someone else note", func(t *testing.T) {
		appHost := "localhost:8080"
		url := fmt.Sprintf("http://%s/api/notes/1", appHost)
		note := models.Note{
			ID:        1,
			BODY:      "update",
			TITLE:     "update",
			ISPRIVATE: false,
			EXPIRE:    1,
		}
		jsonStr, _ := json.MarshalIndent(note, "", "\t")
		r, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		auth := []byte("2:1")
		basic := base64.StdEncoding.EncodeToString(auth)
		authHeader := fmt.Sprintf("Basic %s", basic)
		r.Header.Set("Authorization", authHeader)
		w := httptest.NewRecorder()
		UpdateNote(w, r)
		expectedStatus := 403
		expectedMessage := utils.ForbiddenMessage()
		if expectedStatus != w.Code {
			t.Errorf("Invalid code %d expected %d", w.Code, expectedStatus)
		}
		if !bytes.Equal(expectedMessage, w.Body.Bytes()) {
			t.Errorf("Invalid JSONbody %s expected %s", w.Body.Bytes(), expectedMessage)
		}
	})
}
