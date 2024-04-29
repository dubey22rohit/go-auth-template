package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

func (u UserModel) Insert(user *User) (*UserID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := u.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	var userID uuid.UUID

	query_user := `
    INSERT INTO users (email, password, username) VALUES ($1, $2, $3) RETURNING id`
	args_user := []interface{}{user.Email, user.Password.hash, user.Username}

	if err := tx.QueryRowContext(ctx, query_user, args_user...).Scan(&userID); err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return nil, ErrDuplicateEmail
		default:
			return nil, err
		}
	}

	query_user_profile := `
    INSERT INTO user_profile (user_id) VALUES ($1) ON CONFLICT (user_id) DO NOTHING RETURNING user_id`

	_, err = tx.ExecContext(ctx, query_user_profile, userID)

	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	id := UserID{
		Id: userID,
	}

	return &id, nil
}

func (u UserModel) Get(id uuid.UUID) (*User, error) {
	query := `SELECT u.*, p.* FROM users u
	LEFT JOIN user_profile p ON p.user_id = u.id
	WHERE u.is_active = true AND u.id = $1`

	var user User
	var userProfile UserProfile
	//TODO: Get user posts as well

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, id).Scan(&user.ID,
		&user.Email,
		&user.Password,
		&user.Username,
		&user.IsActive,
		&user.IsStaff,
		&user.IsSuperuser,
		&user.DateJoined,
		&userProfile.ID,
		&userProfile.GeoLocation,
		&userProfile.Age,
		&userProfile.Thumbnail,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	user.Profile = userProfile
	return &user, nil
}

func (u UserModel) GetByEmail(email string) (*User, error) {
	query := `SELECT u.*, p.* FROM
	users u LEFT JOIN user_profile p on p.user_id = u.id
	WHERE u.email = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User

	err := u.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password.hash,
		&user.Username,
		&user.IsActive,
		&user.IsStaff,
		&user.IsSuperuser,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (u UserModel) ActivateUser(userID uuid.UUID) (*sql.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `UPDATE users SET is_active = true WHERE id = $1`

	result, err := u.DB.ExecContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
