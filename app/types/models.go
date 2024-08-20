package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id                 uuid.UUID
	FirstName          string
	LastName           string
	Email              string
	PasswordHash       string
	Avatar             sql.NullString
	LastLogin          sql.NullTime
	FailedLoginAttempt int
	EmailVerifiedAt    sql.NullTime
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          sql.NullTime
}
