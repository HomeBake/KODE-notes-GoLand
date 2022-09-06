package db

import (
	"database/sql"
	"errors"
	"fmt"
	"note-API/models"
	"note-API/utils"
	"os"
	"time"
)

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

func isUserNote(userID int, noteID int) bool {
	var note models.Note
	db := OpenConnection()
	query := "SELECT ID, USERID FROM NOTE WHERE ID = $1"
	row := db.QueryRow(query, noteID)
	err := row.Scan(&note.ID, &note.USERID)
	if err != nil || note.ID == 0 || note.USERID != userID {
		return false
	}
	return true
}

func isUserHaveAccess(userID int, noteID int) bool {
	if IsUserNote(userID, noteID) {
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

func deleteNoteInTime(noteID int, second int, ch chan bool) {
	time.Sleep(time.Duration(second) * time.Second)
	db := OpenConnection()
	query := "DELETE FROM NOTE WHERE ID = $1"
	_, err := db.Exec(query, noteID)
	if err != nil {
		utils.ErrorLog.Printf("BD ERROR: %s", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			utils.ErrorLog.Printf("BD ERROR: %s", err)
		}
	}(db)
	ch <- true
}

func addUser(log string, pass string) (ok bool, err error) {
	bd := OpenConnection()

	query := "SELECT ID FROM USERS WHERE LOGIN = $1"
	var id = 0
	rows := bd.QueryRow(query, log)
	err = rows.Scan(&id)
	if err == nil || id != 0 {
		utils.ErrorLog.Printf("BD ERROR: %s", err)
		defer func(bd *sql.DB) {
			err := bd.Close()
			if err != nil {
				utils.ErrorLog.Printf("BD ERROR: %s", err)
			}
		}(bd)
		return false, err
	}
	query = "INSERT INTO USERS (login, password) VALUES ($1, $2)"
	_, err = bd.Exec(query, log, pass)
	if err != nil {
		utils.ErrorLog.Printf("BD ERROR: %s", err)
		defer func(bd *sql.DB) {
			err := bd.Close()
			if err != nil {
				utils.ErrorLog.Printf("BD ERROR: %s", err)
			}
		}(bd)
		return false, err
	}
	defer func(bd *sql.DB) {
		err := bd.Close()
		if err != nil {
			utils.ErrorLog.Printf("BD ERROR: %s", err)
		}
	}(bd)
	return true, nil
}

func getNotes(sortField string, userID int) ([]models.Note, error) {
	bd := OpenConnection()

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
		if _, err := isItemExists(sortFields, sortField); err == nil {
			query = "SELECT * FROM NOTE WHERE USERID = $1 ORDER BY " + sortField
		} else {
			return nil, errors.New("sort field is not exist")
		}

	}
	rows, err := bd.Query(query, userID)
	if err != nil {
		return nil, err
	}
	var notes []models.Note

	for rows.Next() {
		var note models.Note
		err := rows.Scan(&note.ID, &note.BODY, &note.TITLE, &note.EXPIRE, &note.ISPRIVATE, &note.USERID)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			utils.ErrorLog.Printf("BD ERROR: %s", err)
		}
	}(rows)
	defer func(bd *sql.DB) {
		err := bd.Close()
		if err != nil {
			utils.ErrorLog.Printf("BD ERROR: %s", err)
		}
	}(bd)
	return notes, nil
}

func isNote(noteID string) bool {
	bd := OpenConnection()
	const query = "SELECT ID FROM NOTE WHERE ID = $1"
	row := bd.QueryRow(query, noteID)
	id := 0
	err := row.Scan(&id)
	if err != nil || id == 0 {
		defer func(bd *sql.DB) {
			err := bd.Close()
			if err != nil {
				utils.ErrorLog.Printf("BD ERROR: %s", err)
			}
		}(bd)
		return false
	}
	return true
}

func getNote(noteID string) (models.Note, error) {
	var note models.Note
	bd := OpenConnection()
	const query = "SELECT * FROM NOTE WHERE ID = $1"
	row := bd.QueryRow(query, noteID)
	err := row.Scan(&note.ID, &note.BODY, &note.TITLE, &note.EXPIRE, &note.ISPRIVATE, &note.USERID)
	if err != nil {
		defer func(bd *sql.DB) {
			err := bd.Close()
			if err != nil {
				utils.ErrorLog.Printf("BD ERROR: %s", err)
			}
		}(bd)
		return note, err
	}
	defer func(bd *sql.DB) {
		err := bd.Close()
		if err != nil {
			utils.ErrorLog.Printf("BD ERROR: %s", err)
		}
	}(bd)
	return note, nil
}

func getAccessID(userAccessID int, noteID int) int {
	db := OpenConnection()
	query := "SELECT ID FROM ACCESS WHERE USERID = $1 AND NOTEID = $2"
	row := db.QueryRow(query, userAccessID, noteID)
	id := 0
	_ = row.Scan(&id)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			utils.ErrorLog.Printf("BD ERROR: %s", err)
		}
	}(db)
	return id
}

func addAccess(userAccessID int, noteID int) (int64, error) {
	db := OpenConnection()
	query := "INSERT INTO ACCESS (userid, noteid) VALUES ($1, $2)"
	result, err := db.Exec(query, userAccessID, noteID)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			utils.ErrorLog.Printf("BD ERROR: %s", err)
		}
	}(db)
	var id int64
	id, err = result.LastInsertId()
	if err != nil {
		utils.ErrorLog.Printf("BD ERROR: %s", err)
	}
	return id, err
}

func deleteAccess(accessID int) (int, error) {
	db := OpenConnection()
	query := "DELETE FROM ACCESS WHERE ID = $1"
	_, err := db.Exec(query, accessID)
	if err != nil {
		utils.ErrorLog.Printf("BD ERROR: %s", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			utils.ErrorLog.Printf("BD ERROR: %s", err)
		}
	}(db)
	return accessID, err
}

func addNote(note models.Note, userID int) (int, error) {
	db := OpenConnection()
	const query = `INSERT INTO NOTE (body, title, expire, isPrivate, userid) VALUES ($1, $2, $3, $4, $5) RETURNING ID`
	raw := db.QueryRow(query, note.BODY, note.TITLE, note.EXPIRE, note.ISPRIVATE, userID)
	noteID := 0
	err := raw.Scan(&noteID)
	if err != nil {
		utils.ErrorLog.Printf("BD ERROR: %s", err)
		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {
				utils.ErrorLog.Printf("BD ERROR: %s", err)
			}
		}(db)
		return 0, err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			utils.ErrorLog.Printf("BD ERROR: %s", err)
		}
	}(db)
	return noteID, nil
}

func updateNote(note models.Note, userID int) (int, error) {
	db := OpenConnection()
	const query = `UPDATE NOTE SET BODY = $2, TITLE = $3, EXPIRE = $4, ISPRIVATE = $5 WHERE ID = $1 AND USERID = $6 RETURNING ID`
	row := db.QueryRow(query, note.ID, note.BODY, note.TITLE, note.EXPIRE, note.ISPRIVATE, userID)
	id := 0
	err := row.Scan(&id)
	if err != nil {
		utils.ErrorLog.Printf("BD ERROR: %s", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			utils.ErrorLog.Printf("BD ERROR: %s", err)
		}
	}(db)
	return id, err
}

func deleteNote(noteID int) (int, error) {
	db := OpenConnection()
	const query = "DELETE FROM NOTE WHERE ID = $1 RETURNING ID"
	row := db.QueryRow(query, noteID)
	id := 0
	err := row.Scan(&id)
	if err != nil {
		utils.ErrorLog.Printf("BD ERROR: %s", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			utils.ErrorLog.Printf("BD ERROR: %s", err)
		}
	}(db)
	return id, err
}
