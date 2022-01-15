package models

import (
	"fmt"
	"reflect"
	"strings"

	"database/sql"

	hashPassword "github.com/Park-Kwonsoo/moving-server/pkg/hashing-password"
	qb "github.com/Park-Kwonsoo/moving-server/pkg/query-builder"
)

type User struct {
	baseType
	ProfileID      uint           `model:"foreign key"`
	UserPlaylistID uint           `model:"foreign key"`
	UserId         sql.NullString `model:"varchar(255) unique not null"`
	UserType       string         `model:"varchar(10) not null"`
	Password       string         `model:"varchar(2000) not null"`
}

//테이블 생성
func userMigrate() bool {

	u := User{}

	userId, _ := reflect.TypeOf(u).FieldByName("UserId")
	userType, _ := reflect.TypeOf(u).FieldByName("UserType")
	password, _ := reflect.TypeOf(u).FieldByName("Password")

	query := qb.CreateTable("dev.user").TableOption([]string{
		strings.Join(getCreatedTableColumn(), ", "),
		"userId " + userId.Tag.Get("model"),
		"userType " + userType.Tag.Get("model"),
		"password " + password.Tag.Get("model"),
	}).ToString()

	fmt.Println(query)

	return true
}

//새 유저 등록 : C
func CreateNewUser(user User) bool {
	//해쉬 비밀번호 생성
	hashed, _ := hashPassword.GenerateHashPassword(user.Password)

	//해쉬 비밀번호 생성된 데이터를 넣는다.
	query := qb.Insert("dev.user", "user_id, user_type, password").Value([]string{
		user.UserId.String,
		user.UserType,
		hashed,
	}).ToString()

	//toDo : 새 유저를 생성하는 쿼리 필요
	fmt.Println(query)

	return true
}

//해쉬 패스워드 검증
func (u *User) ValidatePassword(pw string) bool {
	valid, _ := hashPassword.CompareHashPassword(u.Password, pw)
	return valid
}

//아래부턴 쿼리
func FindUserById(id uint) *User {
	user := &User{}

	query := qb.Select("user_id, user_type, password").From("dev.user").Where("id", id).ToString()
	psql.db.QueryRow(query).Scan(
		&user.UserId, &user.UserType, &user.Password)

	return user
}

func FindUserByUserId(userId string) *User {
	user := &User{}

	query := qb.Select("user_id, user_type, password").From("dev.user").Where("user_id", userId).ToString()
	psql.db.QueryRow(query).Scan(
		&user.UserId, &user.UserType, &user.Password)

	if !user.UserId.Valid {
		return nil
	}

	return user
}
