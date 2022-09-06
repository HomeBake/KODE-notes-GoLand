package db

import (
	"database/sql"
	"fmt"
	"note-API/models"
	"os"
	"reflect"
)

func getType() string {
	return os.Getenv("bdType")
}

func OpenConnection() *sql.DB {
	dialect := os.Getenv("dialect")
	host := os.Getenv("host")
	port := os.Getenv("port")
	user := os.Getenv("user")
	password := os.Getenv("password")
	dbname := os.Getenv("dbname")
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
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

func IsUser(l string, p string) (bool, int) {
	switch getType() {
	case "postgres":
		{
			return isUser(l, p)
		}
	default:
		{
			return isUserDummy(l, p)

		}
	}
}

func IsUserNote(userID int, noteID int) bool {
	switch getType() {
	case "postgres":
		{
			return isUserNote(userID, noteID)
		}
	default:
		{
			return isUserNoteDummy(userID, noteID)
		}
	}
}

func IsUserHaveAccess(userID int, noteID int) bool {
	switch getType() {
	case "postgres":
		{
			return isUserHaveAccess(userID, noteID)
		}
	default:
		{
			return isUserHaveAccessDummy(userID, noteID)
		}
	}
}

func DeleteNoteInTime(noteID int, second int, ch chan bool) {
	switch getType() {
	case "postgres":
		{
			deleteNoteInTime(noteID, second, ch)
		}
	default:
		{
			deleteNoteInTimeDummy(noteID, second, ch)
		}
	}
}

func AddUser(log string, pass string) (ok bool, err error) {
	switch getType() {
	case "postgres":
		{
			return addUser(log, pass)
		}
	default:
		{
			return addUserDummy(log, pass)
		}
	}
}

func GetNotes(sortField string, userID int) ([]models.Note, error) {
	switch getType() {
	case "postgres":
		{
			return getNotes(sortField, userID)
		}
	default:
		{
			return getNotesDummy(sortField, userID)
		}
	}
}

func IsNote(noteID string) bool {
	switch getType() {
	case "postgres":
		{
			return isNote(noteID)
		}
	default:
		{
			return isNoteDummy(noteID)
		}
	}
}

func GetNote(noteID string) (models.Note, error) {
	switch getType() {
	case "postgres":
		{
			return getNote(noteID)
		}
	default:
		{
			return getNoteDummy(noteID)
		}
	}
}

func GetAccessID(userAccessID int, noteID int) int {
	switch getType() {
	case "postgres":
		{
			return GetAccessID(userAccessID, noteID)
		}
	default:
		{
			return getAccessIDDummy(userAccessID, noteID)
		}
	}
}

func AddAccess(userAccessID int, noteID int) (int64, error) {
	switch getType() {
	case "postgres":
		{
			return addAccess(userAccessID, noteID)
		}
	default:
		{
			return addAccessDummy(userAccessID, noteID)
		}
	}
}

func DeleteAccess(accessID int) (int, error) {
	switch getType() {
	case "postgres":
		{
			return deleteAccess(accessID)
		}
	default:
		{
			return deleteAccessDummy(accessID)
		}
	}
}

func AddNote(note models.Note, userID int) (int, error) {
	switch getType() {
	case "postgres":
		{
			return addNote(note, userID)
		}
	default:
		{
			return addNoteDummy(note, userID)
		}
	}
}

func UpdateNote(note models.Note, userID int) (int, error) {
	switch getType() {
	case "postgres":
		{
			return updateNote(note, userID)
		}
	default:
		{
			return updateNoteDummy(note, userID)
		}
	}
}

func DeleteNote(noteID int) (int, error) {
	switch getType() {
	case "postgres":
		{
			return deleteNote(noteID)
		}
	default:
		{
			return deleteNoteDummy(noteID)
		}
	}
}
