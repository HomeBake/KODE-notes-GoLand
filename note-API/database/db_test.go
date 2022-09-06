package db

import (
	"encoding/json"
	"note-API/models"
	"os"
	"reflect"
	"strconv"
	"testing"
)

func fillDB() {
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
			BODY:      "za",
			TITLE:     "az",
			ISPRIVATE: false,
			EXPIRE:    0,
			USERID:    3,
		},
		6: {
			ID:        6,
			BODY:      "az",
			TITLE:     "za",
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

func TestMain(m *testing.M) {
	fillDB()
	code := m.Run()
	os.Exit(code)
}

func TestIsItemExist(t *testing.T) {
	arr := [2]string{
		"first",
		"second",
	}
	t.Run("true", func(t *testing.T) {
		item := "first"
		result, _ := isItemExists(arr, item)
		if !result {
			t.Errorf("Expected true. Got false")
		}
	})
	t.Run("false", func(t *testing.T) {
		item := "third"
		result, _ := isItemExists(arr, item)
		if result {
			t.Errorf("Expected false. Got true")
		}
	})
	t.Run("Invalid data type", func(t *testing.T) {
		arr := map[string]string{
			"firstKey":  "firstValue",
			"secondKey": "secondValue",
		}
		item := "third"
		_, err := isItemExists(arr, item)
		if err == nil {
			t.Errorf("Expected error")
		}
	})
}

func TestGetLastNoteId(t *testing.T) {
	expected := 6
	result := getLastNoteId()
	if result != expected {
		t.Errorf("Expected %d, %d returned", expected, result)
	}
}

func TestGetLastAccessId(t *testing.T) {
	expected := 1
	result := getLastAccessId()
	if result != expected {
		t.Errorf("Expected %d, %d returned", expected, result)
	}
}

func TestIsUser(t *testing.T) {
	t.Run("True", func(t *testing.T) {
		login := "login"
		password := "password"
		_, _ = addUserDummy(login, password)
		result, _ := isUserDummy(login, password)
		if !result {
			t.Errorf("Expected true, false returned")
		}
	})
	t.Run("False", func(t *testing.T) {
		login := "UnExist"
		password := "UnExist"
		result, _ := isUserDummy(login, password)
		if result {
			t.Errorf("Expected false, true returned")
		}
	})
}

func TestAddUser(t *testing.T) {
	login := "login3"
	password := "password3"
	t.Run("Success", func(t *testing.T) {
		result, _ := addUserDummy(login, password)
		if !result {
			t.Errorf("Expected true, false returned")
		}
	})
	t.Run("unSuccess", func(t *testing.T) {
		result, _ := addUserDummy(login, password)
		if result {
			t.Errorf("Expected false, true returned")
		}
	})
}

func TestIsUserNote(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		userID := 1
		noteID := 1
		result := isUserNoteDummy(userID, noteID)
		if !result {
			t.Errorf("Expected true, false returned")
		}
	})
	t.Run("false", func(t *testing.T) {
		userID := 2
		noteID := 1
		result := isUserNoteDummy(userID, noteID)
		if result {
			t.Errorf("Expected false, true returned")
		}
	})
}

func TestIsUserHaveAccess(t *testing.T) {
	t.Run("user have access", func(t *testing.T) {
		userID := 2
		noteID := 1
		result := isUserHaveAccessDummy(userID, noteID)
		if !result {
			t.Errorf("Expected true, false returned")
		}
	})
	t.Run("user dont have access", func(t *testing.T) {
		userID := 2
		noteID := 2
		result := isUserHaveAccessDummy(userID, noteID)
		if result {
			t.Errorf("Expected false, true returned")
		}
	})
	t.Run("Its user note", func(t *testing.T) {
		userID := 1
		noteID := 2
		result := isUserHaveAccessDummy(userID, noteID)
		if !result {
			t.Errorf("Expected true, false returned")
		}
	})
}

func TestDeleteNoteInTimeDummy(t *testing.T) {
	t.Run("noteExistAndDelete", func(t *testing.T) {
		noteID := 3
		second := 1
		ch := make(chan bool)
		go deleteNoteInTimeDummy(noteID, second, ch)
		if !isNoteDummy(strconv.Itoa(noteID)) {
			t.Errorf("Note deleted before time")
		}
		<-ch
		if isNoteDummy(strconv.Itoa(noteID)) {
			t.Errorf("Note is not deleted")
		}
	})
	t.Run("noteDoesntExist", func(t *testing.T) {
		noteID := 7
		second := 1
		ch := make(chan bool)
		go deleteNoteInTimeDummy(noteID, second, ch)
		result := <-ch
		if result {
			t.Errorf("Note is not deleted")
		}
	})
}

func TestGetNotes(t *testing.T) {
	t.Run("User have notes and sortField empty", func(t *testing.T) {
		userID := 2
		sortField := ""
		notes, err := getNotesDummy(sortField, userID)
		expectedNotes := []models.Note{
			{
				ID:        4,
				BODY:      "1",
				TITLE:     "2",
				ISPRIVATE: false,
				EXPIRE:    0,
				USERID:    2,
			},
		}
		if err != nil || !reflect.DeepEqual(notes, expectedNotes) {
			expectedBytes, _ := json.MarshalIndent(expectedNotes, "", "\t")
			notesByte, _ := json.MarshalIndent(notes, "", "\t")
			t.Errorf("Wrong notes return, expect :%s, returned: %s ", expectedBytes, notesByte)
		}
	})
	t.Run("Bad sort field ", func(t *testing.T) {
		userID := 2
		sortField := "bad"
		_, err := getNotesDummy(sortField, userID)
		if err == nil {
			t.Errorf("Error expect")
		}
	})
	t.Run("User have notes and sortField +", func(t *testing.T) {
		userID := 3
		sortField := "body"
		notes, err := getNotesDummy(sortField, userID)
		expectedNotes := []models.Note{
			{
				ID:        6,
				BODY:      "az",
				TITLE:     "za",
				ISPRIVATE: false,
				EXPIRE:    0,
				USERID:    3,
			},
			{
				ID:        5,
				BODY:      "za",
				TITLE:     "az",
				ISPRIVATE: false,
				EXPIRE:    0,
				USERID:    3,
			},
		}
		if err != nil || !reflect.DeepEqual(notes, expectedNotes) {
			expectedBytes, _ := json.MarshalIndent(expectedNotes, "", "\t")
			notesByte, _ := json.MarshalIndent(notes, "", "\t")
			t.Errorf("Wrong notes return, expect :%s, returned: %s ", expectedBytes, notesByte)
		}
	})
	t.Run("User have notes and sortField -", func(t *testing.T) {
		userID := 3
		sortField := "-body"
		notes, err := getNotesDummy(sortField, userID)
		expectedNotes := []models.Note{
			{
				ID:        5,
				BODY:      "za",
				TITLE:     "az",
				ISPRIVATE: false,
				EXPIRE:    0,
				USERID:    3,
			},
			{
				ID:        6,
				BODY:      "az",
				TITLE:     "za",
				ISPRIVATE: false,
				EXPIRE:    0,
				USERID:    3,
			},
		}
		if err != nil || !reflect.DeepEqual(notes, expectedNotes) {
			expectedBytes, _ := json.MarshalIndent(expectedNotes, "", "\t")
			notesByte, _ := json.MarshalIndent(notes, "", "\t")
			t.Errorf("Wrong notes return, expect :%s, returned: %s ", expectedBytes, notesByte)
		}
	})
}

func TestGetNote(t *testing.T) {
	noteId := "1"
	expectedNote := models.Note{
		ID:        1,
		BODY:      "1",
		TITLE:     "2",
		ISPRIVATE: true,
		EXPIRE:    0,
		USERID:    1,
	}
	note, _ := getNoteDummy(noteId)
	if !reflect.DeepEqual(expectedNote, note) {
		expectedBytes, _ := json.MarshalIndent(expectedNote, "", "\t")
		notesByte, _ := json.MarshalIndent(note, "", "\t")
		t.Errorf("Wrong notes return, expect :%s, returned: %s ", expectedBytes, notesByte)
	}
}

func TestGetAccessID(t *testing.T) {
	t.Run("access id exist", func(t *testing.T) {
		userAccessID := 2
		noteID := 1
		expectedID := 1
		id := getAccessIDDummy(userAccessID, noteID)
		if id != expectedID {
			t.Errorf("Wrong id return, expect :%d, returned: %d ", expectedID, id)
		}
	})
	t.Run("access id exist", func(t *testing.T) {
		userAccessID := 2
		noteID := 2
		expectedID := 0
		id := getAccessIDDummy(userAccessID, noteID)
		if id != expectedID {
			t.Errorf("Wrong id return, expect :%d, returned: %d ", expectedID, id)
		}
	})
}

func TestAddAccess(t *testing.T) {
	t.Run("access already exist", func(t *testing.T) {
		userAccessID := 2
		noteID := 1
		expectedID := 0
		id, _ := addAccessDummy(userAccessID, noteID)
		if id != int64(expectedID) {
			t.Errorf("Wrong id return, expect :%d, returned: %d ", expectedID, id)
		}
	})
	t.Run("success add", func(t *testing.T) {
		userAccessID := 1
		noteID := 6
		expectedID := 2
		id, _ := addAccessDummy(userAccessID, noteID)
		if id != int64(expectedID) {
			t.Errorf("Wrong id return, expect :%d, returned: %d ", expectedID, id)
		}
	})
}

func TestDeleteAccess(t *testing.T) {
	t.Run("access not exist", func(t *testing.T) {
		accessID := 3
		expectedID := 0
		id, _ := deleteAccessDummy(accessID)
		if id != expectedID {
			t.Errorf("Wrong id return, expect :%d, returned: %d ", expectedID, id)
		}
	})
	t.Run("success delete", func(t *testing.T) {
		accessID := 1
		expectedID := 0
		_, found := AccessDateDB[accessID]
		if !found {
			t.Errorf("success already delete")
		}
		id, _ := deleteAccessDummy(accessID)
		_, found = AccessDateDB[accessID]
		if found {
			t.Errorf("Wrong id return, expect :%d, returned: %d ", expectedID, id)
		}
	})
}

func TestAddNote(t *testing.T) {
	note := models.Note{
		BODY:      "dsfdsf",
		TITLE:     "dsfsdf",
		ISPRIVATE: false,
		EXPIRE:    0,
	}
	userID := 1
	expectedID := getLastNoteId() + 1
	id, _ := addNoteDummy(note, userID)
	note.USERID = userID
	note.ID = expectedID
	expectedNote := note
	addedNote, _ := getNoteDummy(strconv.Itoa(id))
	if !reflect.DeepEqual(addedNote, expectedNote) || expectedID != id {
		expectedBytes, _ := json.MarshalIndent(expectedNote, "", "\t")
		notesByte, _ := json.MarshalIndent(addedNote, "", "\t")
		t.Errorf("Wrong notes return, expect :%s, returned: %s ", expectedBytes, notesByte)
	}
}

func TestUpdateNote(t *testing.T) {
	note := models.Note{
		ID:        1,
		BODY:      "sadfdsf",
		TITLE:     "dsds",
		ISPRIVATE: false,
		EXPIRE:    0,
	}
	userID := 1

	id, _ := updateNoteDummy(note, userID)
	note.USERID = userID
	expectedNote := note
	addedNote, _ := getNoteDummy(strconv.Itoa(id))
	if !reflect.DeepEqual(addedNote, expectedNote) || note.ID != id {
		expectedBytes, _ := json.MarshalIndent(expectedNote, "", "\t")
		notesByte, _ := json.MarshalIndent(addedNote, "", "\t")
		t.Errorf("Wrong notes return, expect :%s, returned: %s ", expectedBytes, notesByte)
	}
}

func TestDeleteNote(t *testing.T) {
	noteID := 1
	id, _ := deleteNoteDummy(noteID)
	isDeleted := isNoteDummy(strconv.Itoa(noteID))
	if isDeleted || id != noteID {
		t.Error("Note is not deleted")
	}
}
