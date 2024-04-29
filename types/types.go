package types

import "github.com/google/uuid"

type User struct {
	User_id  int64
	Username string
	Online   bool
}

type UserID struct {
	Id uuid.UUID
}
