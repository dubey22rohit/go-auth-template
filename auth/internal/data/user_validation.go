package data

import "github.com/dubey22rohit/heyyy_yo_backend/auth/internal/validator"

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "email must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "email must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "password must be provided")
	v.Check(len(password) >= 8, "password", "password must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "password must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Username != "", "username", "username must be provided")

	ValidateEmail(v, user.Email)
	// If the plaintext password is not nil, call the standalone // ValidatePasswordPlaintext() helper.
	if user.Password.plainText != nil {
		ValidatePasswordPlaintext(v, *user.Password.plainText)
	}
}
