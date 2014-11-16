package store

import (
	"encoding/base64"
	"fmt"

	"code.google.com/p/go.crypto/bcrypt"
)

func (scope UserScope) Authenticate(username, password string) *User {
	fmt.Println("AuthUser", username, password)
	user, err := scope.Email().Eq(username).Retrieve()
	if err != nil {
		fmt.Println("Couldn't retrieve", err)
		return nil
	}
	if err = bcrypt.CompareHashAndPassword(user.PasswordBytes(), []byte(password)); err != nil {
		fmt.Println("bcrypt error", err)
	} else {
		return &user
	}
	return nil
}

func (u User) PasswordBytes() []byte {
	if len(u.Password) == 0 {
		return []byte{}
	}
	bpw, err := base64.StdEncoding.DecodeString(u.Password)
	if err != nil {
		return []byte{}
	}
	return bpw
}

var IdiotPasswords = map[string]bool{
	"12345678":   true,
	"password":   true,
	"qwertyui":   true,
	"123456789":  true,
	"1234567890": true,
	"12341234":   true,
	"sunshine":   true,
	"password1":  true,
	"00000000":   true,
	"trustno1":   true,
	"abcd1234":   true,
	"iloveyou":   true,
	"yourself":   true,
	"princess":   true,
}

func (u *User) SetPassword(pw string) error {
	if IdiotPasswords[pw] {
		return fmt.Errorf("Password is too stupid")
	}
	bpw, err := bcrypt.GenerateFromPassword([]byte(pw), 10)
	if err != nil {
		return err
	}
	u.Password = base64.StdEncoding.EncodeToString(bpw)
	return nil
}
