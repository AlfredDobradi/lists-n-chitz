package auth

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/alfreddobradi/lists-n-chitz/internal/database"
)

type User struct {
	ID        uint   `json:"id" db:"id"`
	Status    uint   `json:"status" db:"status"`
	CreatedAt int64  `json:"created_at" db:"created_at"`
	Email     string `json:"email" db:"email"`
	Password  string `json:"password" db:"password"`
}

type UserResponse struct {
	ID        uint   `json:"id"`
	Status    uint   `json:"status"`
	CreatedAt int64  `json:"created_at"`
	Email     string `json:"email"`
}

type Token struct {
	userID  uint
	Token   string `json:"token"`
	Expires int64  `json:"expires"`
	address string
}

func Register(u User) (User, error) {
	db, err := database.New("postgres://postgres@localhost:5432/lists?sslmode=disable")
	if err != nil {
		return u, err
	}

	tx := db.MustBegin()

	createdAt := time.Now().Format(time.RFC3339Nano)
	pwd := createHash(u.Password)

	query := "INSERT INTO users (email, password, created_at) VALUES ($1, $2, $3)"

	res, err := tx.Exec(query, u.Email, pwd, createdAt)

	err = database.Try(err, tx)
	if err != nil || res == nil {
		return User{}, err
	}

	query = "SELECT id, status, created_at FROM users WHERE email = $1 AND password = $2"
	rows, err := db.Query(query, u.Email, pwd)
	if err != nil {
		return User{}, err
	}

	var id uint
	var status uint
	var c string
	for rows.Next() {
		rows.Scan(&id, &status, &c)
	}

	cInt, _ := time.Parse(time.RFC3339Nano, c)

	u.ID = id
	u.Status = status
	u.CreatedAt = cInt.Unix()

	return u, nil
}

func Authenticate(u User, address string) (Token, error) {
	var t Token
	db, err := database.New("postgres://postgres@localhost:5432/lists?sslmode=disable")
	if err != nil {
		return t, err
	}

	pwd := createHash(u.Password)
	query := "SELECT id FROM users WHERE email = $1 AND password = $2"
	res, err := db.Query(query, u.Email, pwd)
	if err != nil {
		return t, err
	}

	var i uint
	var id uint
	for res.Next() {
		i = i + 1
		res.Scan(&id)
	}

	if i != 1 {
		return t, errors.New("invalid credentials")
	}

	t.userID = id
	t.address = address

	tokenString := fmt.Sprintf("%s%d", pwd, time.Now().UnixNano())
	t.Token = createHash(tokenString)

	t, err = createToken(t)

	return t, err
}

func Authorize(address, token string) (u User, err error) {
	db, err := database.New("postgres://postgres@localhost:5432/lists?sslmode=disable")
	if err != nil {
		return
	}

	query := "SELECT u.* FROM user_tokens ut LEFT JOIN users u ON ut.iduser = u.id WHERE token = $1 AND address = $2 AND u.status = 1 AND NOW() <= expires LIMIT 1"
	res := db.QueryRowx(query, token, address)
	if res.Err() != nil {
		return
	}

	var cols struct {
		ID        uint   `db:"id"`
		Email     string `db:"email"`
		Status    uint   `db:"status"`
		Password  string `db:"password"`
		CreatedAt string `db:"created_at"`
	}
	err = res.StructScan(&cols)
	if err != nil {
		return
	}

	t, err := time.Parse(time.RFC3339, cols.CreatedAt)
	if err != nil {
		return
	}

	u.Email = cols.Email
	u.ID = cols.ID
	u.Status = cols.Status
	u.CreatedAt = t.Unix()

	return
}

func createToken(t Token) (Token, error) {
	db, err := database.New("postgres://postgres@localhost:5432/lists?sslmode=disable")
	if err != nil {
		return t, err
	}

	tx := db.MustBegin()
	query := "SELECT token FROM user_tokens WHERE iduser = $1 AND address = $2 AND NOW() <= expires"
	res, err := db.Query(query, t.userID, t.address)
	if err != nil {
		return t, err
	}
	var tTmp string
	var i uint
	for res.Next() {
		i = i + 1
		res.Scan(&tTmp)
	}

	expires := time.Now().Add(30 * time.Minute)
	t.Expires = expires.Unix()
	expiresTimestamp := expires.Format(time.RFC3339Nano)

	if i > 0 {
		t.Token = tTmp
		query = "UPDATE user_tokens SET expires = $1 WHERE token = $2"

		_, err = tx.Exec(query, expiresTimestamp, tTmp)
		err = database.Try(err, tx)
		if err != nil {
			return t, err
		}
	} else {
		query = "INSERT INTO user_tokens (iduser, token, expires, address) VALUES ($1, $2, $3, $4)"

		_, err = tx.Exec(query, t.userID, t.Token, expiresTimestamp, t.address)

		err = database.Try(err, tx)
		if err != nil {
			return t, err
		}
	}

	return t, err
}

func createHash(s string) (h string) {
	hh := sha256.New()
	hh.Write([]byte(s))
	h = fmt.Sprintf("%x", hh.Sum(nil))
	return
}
