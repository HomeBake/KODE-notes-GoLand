package models

type Note struct {
	ID        int    `json:"id"`
	BODY      string `json:"body"`
	TITLE     string `json:"title"`
	ISPRIVATE bool   `json:"isPrivate"`
	EXPIRE    int    `json:"expire"`
	USERID    int    `json:"userid"`
}

type UserData struct {
	ID       int    `json:"ID"`
	LOGIN    string `json:"login"`
	PASSWORD string `json:"password"`
}

type AccessDate struct {
	ID           int    `json:"ID"`
	USERACCESSID int    `json:"userAccessId"`
	NOTEID       int    `json:"noteId"`
	MODE         string `json:"mode"`
}
