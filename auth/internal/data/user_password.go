package data

import (
	"log"

	"github.com/alexedwards/argon2id"
)

func (p *password) Set(plainTextPassword string) error {
	hash, err := argon2id.CreateHash(plainTextPassword, argon2id.DefaultParams)
	if err != nil {
		return err
	}

	p.plainText = &plainTextPassword
	p.hash = hash

	return nil
}

func (p *password) Matches(plainTextPassword string) (bool, error) {
	matches, err := argon2id.ComparePasswordAndHash(plainTextPassword, p.hash)
	if err != nil {
		log.Fatal(err)
		return false, err
	}

	return matches, nil
}
