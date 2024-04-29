package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type password struct {
	plainText *string
	hash      string
}

type User struct {
	ID          uuid.UUID   `json:"id"`
	Email       string      `json:"email"`
	Password    password    `json:"-"`
	Username    string      `json:"username"`
	IsActive    bool        `json:"is_active"`
	IsStaff     bool        `json:"is_staff"`
	IsSuperuser bool        `json:"is_superuser"`
	DateJoined  time.Time   `json:"date_joined"`
	Profile     UserProfile `json:"profile"`
}

type UserProfile struct {
	ID          *uuid.UUID  `json:"id"`
	UserID      *uuid.UUID  `json:"user_id"`
	Age         uint8       `json:"age"`
	GeoLocation string      `json:"geo_location"`
	Thumbnail   string      `json:"thumbnail"`
	Posts       []UserPosts `json:"user_posts"`
}

type UserPosts struct {
	ID            *uuid.UUID `json:"id"`
	UserProfileID *uuid.UUID `json:"user_profile_id"`
	PostContent   string     `json:"post_content"`
	PostImageURL  string     `json:"post_image_URL"`
	PostDate      time.Time  `json:"post_date"`
}

type UserID struct {
	Id uuid.UUID
}

type UserModel struct {
	DB *sql.DB
}

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)
