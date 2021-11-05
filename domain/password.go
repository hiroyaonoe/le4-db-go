package domain

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Password は受け取った文字列をハッシュ化して扱うためのstruct
type Password struct {
	value        string
	is_encrypted bool
}

func NewPassword(s string) (Password, error) {
	res := Password{value: s, is_encrypted: false}
	err := res.Encrypt()
	return res, err
}

func (t Password) String() string {
	return t.value
}

func (t *Password) Set(str string) {
	t.value = str
}

func (t *Password) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.value)
}

func (t *Password) UnmarshalJSON(data []byte) error {
	var str string
	json.Unmarshal(data, &str)
	t.Set(str)
	t.is_encrypted = false
	return nil
}

func (t Password) IsNull() bool {
	return t.value == ""
}

func (t Password) Equal(s Password) bool {
	return t.String() == s.String() && t.is_encrypted == s.is_encrypted
}

// Encrypt はトークンをハッシュ化する
func (t *Password) Encrypt() error {
	if t.is_encrypted {
		return fmt.Errorf("Already encrypted:%s", t)
	}
	if t.String() == "" {
		return nil
	}
	token := []byte(t.String())
	digest, err := bcrypt.GenerateFromPassword(token, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	t.Set(string(digest))
	t.is_encrypted = true

	return nil
}

/*
Authenticate は２つのトークンが同一のものか判定する
(tの元となる文字列がplainと等しいかどうか)
*/
func (t *Password) Authenticate(plain string) bool {
	return t.authenticate(plain) == nil
}

// authenticate はBoolではなくエラーを返す
func (t *Password) authenticate(plain string) error {
	// tはハッシュ化されていることが必要
	if !t.is_encrypted {
		return fmt.Errorf("Invalid tokens")
	}
	p1 := []byte(t.String())
	p2 := []byte(plain)

	return bcrypt.CompareHashAndPassword(p1, p2)
}

// Scan はデータベースの値をPasswordにマッピングする
func (t *Password) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("Invalid value:%s", value)
	}
	t.Set(string(str))
	t.is_encrypted = true
	return nil
}

// Value はPasswordのフィールドのうちデータベースに保存するものを指定する
func (t Password) Value() (driver.Value, error) {
	if !t.is_encrypted {
		return nil, fmt.Errorf("Not encrypted")
	}
	return t.value, nil
}
