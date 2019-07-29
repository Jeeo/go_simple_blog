package data

import (
	"log"
	"time"
)

type Session struct {
	Id        int
	Uuid      string
	UserId    int
	CreatedAt time.Time
}

// CreateSession Creates a session based on userId
func CreateSession(userID int) (session Session, err error) {
	statement := "insert into Sessions (uuid, user_id, created_at) values ($1, $2, $3) returning id, uuid, user_id, created_at"
	stmt, err := DB.Prepare(statement)
	defer stmt.Close()
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = stmt.QueryRow(createUUID(), userID, time.Now()).Scan(&session.Id, &session.Uuid, &session.UserId, &session.CreatedAt)

	return
}

// GetSessionByUserID Creates a session based on userId
func GetSessionByUserID(userID string) (session Session, err error) {
	err = DB.QueryRow("SELECT id, user_id, uuid, created_at from Sessions where user_id = $1", userID).Scan(&session.Id, &session.UserId, &session.Uuid, &session.CreatedAt)

	return
}

// Check checks a session based on UUID
func (session *Session) Check() (valid bool, err error) {
	err = DB.QueryRow("SELECT id from Sessions where uuid = $1", session.Uuid).Scan(&session.Id)
	if err != nil {
		return false, err
	}
	if session.Id != 0 {
		return true, nil
	}
	return
}

//DeleteSessionByUUID Deletes a session
func (session *Session) DeleteSessionByUUID() (err error) {
	statement := "delete from Sessions where uuid = $1"
	stmt, err := DB.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	stmt.Exec(session.Uuid)
	return
}
