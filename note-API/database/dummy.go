package db

import (
	"errors"
	"note-API/models"
	"note-API/utils"
	"sort"
	"strconv"
	"time"
)

type SortField struct {
	field string
	order string
}

var (
	UserDB = make(map[int]models.UserData)
)
var (
	NoteDB = make(map[int]models.Note)
)
var (
	AccessDateDB = make(map[int]models.AccessDate)
)

func getLastNoteId() int {
	max := 0
	for _, element := range NoteDB {
		if element.ID > max {
			max = element.ID
		}
	}
	return max
}

func getLastAccessId() int {
	max := 0
	for _, element := range AccessDateDB {
		if element.ID > max {
			max = element.ID
		}
	}
	return max
}

func isUserDummy(l string, p string) (bool, int) {
	for _, element := range UserDB {
		if element.LOGIN == l && element.PASSWORD == p {
			return true, element.ID
		}
	}
	return false, 0
}

func isUserNoteDummy(userID int, noteID int) bool {
	isUserNotes := false
	for _, element := range NoteDB {
		if element.ID == noteID && element.USERID == userID {
			isUserNotes = true
		}
	}
	return isUserNotes

}

func isUserHaveAccessDummy(userID int, noteID int) bool {
	if isUserNoteDummy(userID, noteID) {
		return true
	}
	isAccess := false
	for _, element := range AccessDateDB {
		if element.USERACCESSID == userID && element.NOTEID == noteID {
			isAccess = true
		}
	}
	return isAccess
}

func deleteNoteInTimeDummy(noteID int, second int, ch chan bool) {
	time.Sleep(time.Duration(second) * time.Second)
	ok, err := deleteNoteDummy(noteID)
	if err != nil || ok == 0 {
		//log
		ch <- false
		return
	}
	ch <- true
	return
}

func addUserDummy(log string, pass string) (ok bool, err error) {
	for _, element := range UserDB {
		if element.LOGIN == log {
			return false, errors.New("user exist")
		}
	}
	var user models.UserData
	user.ID = len(UserDB) + 1
	user.LOGIN = log
	user.PASSWORD = pass
	UserDB[user.ID] = user
	return true, nil

}

func getNotesDummy(sortField string, userID int) ([]models.Note, error) {
	var noteBD []models.Note
	for _, element := range NoteDB {
		if element.USERID == userID {
			noteBD = append(noteBD, element)
		}
	}

	if sortField != "" {
		sortFields := map[string]SortField{
			"-id":        {"ID", "-"},
			"id":         {"ID", "+"},
			"-title":     {"TITLE", "-"},
			"title":      {"TITLE", "+"},
			"-body":      {"BODY", "-"},
			"body":       {"BODY", "+"},
			"-expire":    {"EXPIRE", "-"},
			"expire":     {"EXPIRE", "+"},
			"-isprivate": {"ISPRIVATE", "-"},
			"isprivate":  {"ISPRIVATE", "+"},
		}
		if sortInfo, ok := sortFields[sortField]; ok {
			if sortInfo.order == "+" {
				sort.SliceStable(noteBD, func(i, j int) (less bool) {
					noteI := map[string]string{
						"ID":        strconv.Itoa(noteBD[i].ID),
						"BODY":      noteBD[i].BODY,
						"TITLE":     noteBD[i].TITLE,
						"ISPRIVATE": strconv.FormatBool(noteBD[i].ISPRIVATE),
						"EXPIRE":    strconv.Itoa(noteBD[i].EXPIRE),
						"USERID":    strconv.Itoa(noteBD[i].USERID),
					}
					noteJ := map[string]string{
						"ID":        strconv.Itoa(noteBD[j].ID),
						"BODY":      noteBD[j].BODY,
						"TITLE":     noteBD[j].TITLE,
						"ISPRIVATE": strconv.FormatBool(noteBD[j].ISPRIVATE),
						"EXPIRE":    strconv.Itoa(noteBD[j].EXPIRE),
						"USERID":    strconv.Itoa(noteBD[j].USERID),
					}
					return noteI[sortInfo.field] < noteJ[sortInfo.field]
				})
			} else {
				sort.SliceStable(noteBD, func(i, j int) (less bool) {
					noteI := map[string]string{
						"ID":        strconv.Itoa(noteBD[i].ID),
						"BODY":      noteBD[i].BODY,
						"TITLE":     noteBD[i].TITLE,
						"ISPRIVATE": strconv.FormatBool(noteBD[i].ISPRIVATE),
						"EXPIRE":    strconv.Itoa(noteBD[i].EXPIRE),
						"USERID":    strconv.Itoa(noteBD[i].USERID),
					}
					noteJ := map[string]string{
						"ID":        strconv.Itoa(noteBD[j].ID),
						"BODY":      noteBD[j].BODY,
						"TITLE":     noteBD[j].TITLE,
						"ISPRIVATE": strconv.FormatBool(noteBD[j].ISPRIVATE),
						"EXPIRE":    strconv.Itoa(noteBD[j].EXPIRE),
						"USERID":    strconv.Itoa(noteBD[j].USERID),
					}
					return noteI[sortInfo.field] > noteJ[sortInfo.field]
				})
			}
		} else {
			return nil, errors.New("sort field is not exist")
		}
	}
	return noteBD, nil
}

func isNoteDummy(noteID string) bool {
	id, err := strconv.Atoi(noteID)
	_, found := NoteDB[id]
	if err != nil || found == false {
		return false
	}
	return true
}

func getNoteDummy(noteID string) (models.Note, error) {
	id, err := strconv.Atoi(noteID)
	utils.ErrorLog.Printf("BD ERROR: %s", err)
	return NoteDB[id], err
}

func getAccessIDDummy(userAccessID int, noteID int) int {
	for _, element := range AccessDateDB {
		if element.USERACCESSID == userAccessID && element.NOTEID == noteID {
			return element.ID
		}
	}
	return 0
}

func addAccessDummy(userAccessID int, noteID int) (int64, error) {
	if getAccessIDDummy(userAccessID, noteID) != 0 {
		return 0, errors.New("already exist")
	}
	accessItem := models.AccessDate{
		ID:           getLastAccessId() + 1,
		USERACCESSID: userAccessID,
		NOTEID:       noteID,
	}
	AccessDateDB[accessItem.ID] = accessItem
	return int64(accessItem.ID), nil
}

func deleteAccessDummy(accessID int) (int, error) {
	_, found := AccessDateDB[accessID]
	if !found {
		return 0, errors.New("not exist")
	}
	delete(AccessDateDB, accessID)
	return accessID, nil
}

func addNoteDummy(note models.Note, userID int) (int, error) {
	note.USERID = userID
	note.ID = getLastNoteId() + 1
	NoteDB[note.ID] = note
	return note.ID, nil
}

func updateNoteDummy(note models.Note, userID int) (int, error) {
	note.USERID = userID
	NoteDB[note.ID] = note
	return note.ID, nil
}

func deleteNoteDummy(noteID int) (int, error) {
	_, found := NoteDB[noteID]
	if !found {
		return 0, errors.New("not exist")
	}
	delete(NoteDB, noteID)
	return noteID, nil
}

func FillDB() {
	notes := map[int]models.Note{
		1: {
			ID:        1,
			BODY:      "1",
			TITLE:     "2",
			ISPRIVATE: true,
			EXPIRE:    0,
			USERID:    1,
		},
		2: {
			ID:        2,
			BODY:      "1",
			TITLE:     "2",
			ISPRIVATE: true,
			EXPIRE:    0,
			USERID:    1,
		},
		3: {
			ID:        3,
			BODY:      "1",
			TITLE:     "2",
			ISPRIVATE: false,
			EXPIRE:    0,
			USERID:    1,
		},
		4: {
			ID:        4,
			BODY:      "1",
			TITLE:     "2",
			ISPRIVATE: false,
			EXPIRE:    0,
			USERID:    2,
		},
		5: {
			ID:        5,
			BODY:      "1",
			TITLE:     "2",
			ISPRIVATE: false,
			EXPIRE:    0,
			USERID:    3,
		},
		6: {
			ID:        6,
			BODY:      "1",
			TITLE:     "2",
			ISPRIVATE: false,
			EXPIRE:    0,
			USERID:    3,
		},
	}
	NoteDB = notes
	users := map[int]models.UserData{
		1: {
			ID:       1,
			LOGIN:    "1",
			PASSWORD: "1",
		},
		2: {
			ID:       2,
			LOGIN:    "2",
			PASSWORD: "1",
		},
		3: {
			ID:       3,
			LOGIN:    "3",
			PASSWORD: "1",
		},
	}
	UserDB = users
	access := models.AccessDate{
		ID:           1,
		USERACCESSID: 2,
		NOTEID:       1,
	}
	AccessDateDB[access.ID] = access
}
