package models

import (
	"database/sql"

	qb "github.com/Park-Kwonsoo/moving-server/pkg/query-builder"
)

type User struct {
	basicType
	ProfileID      uint
	UserPlaylistID uint
	UserId         sql.NullString
	UserType       string
	Password       string
}

func (u *User) ValidatePassword(pw string) bool {
	return u.Password == pw
}

func FindUserById(id uint) *User {
	user := &User{}

	query := qb.Select("user_id, user_type, password").From("dev.user").Where("id", id).ToString()
	db.QueryRow(query).Scan(
		&user.UserId, &user.UserType, &user.Password)

	return user
}

func FindUserByUserId(userId string) *User {
	user := &User{}

	query := qb.Select("user_id, user_type, password").From("dev.user").Where("user_id", userId).ToString()
	db.QueryRow(query).Scan(
		&user.UserId, &user.UserType, &user.Password)

	if !user.UserId.Valid {
		return nil
	}

	return user
}
