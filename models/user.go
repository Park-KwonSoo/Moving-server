package models

import (
	"errors"
	"reflect"
	"strings"

	"database/sql"

	hashPassword "github.com/Park-Kwonsoo/moving-server/pkg/hashing-password"
	qb "github.com/Park-Kwonsoo/moving-server/pkg/query-builder"
)

type User struct {
	baseType
	ProfileId      uint           `db:"profile_id int references profile(id)"`
	UserPlaylistId uint           `db:"user_playlist_id int references userPlaylist(id)"`
	UserId         sql.NullString `db:"user_id varchar(255) unique not null"`
	UserType       string         `db:"user_type varchar(10) not null"`
	Password       string         `db:"password varchar(2000) not null"`
}

//테이블 생성
func userMigrate() error {

	u := User{}

	userId, _ := reflect.TypeOf(u).FieldByName("UserId")
	userType, _ := reflect.TypeOf(u).FieldByName("UserType")
	password, _ := reflect.TypeOf(u).FieldByName("Password")

	query := qb.CreateTable("user").TableComlumn([]string{
		strings.Join(getCreatedTableColumn(), ", "),
		userId.Tag.Get("db"),
		userType.Tag.Get("db"),
		password.Tag.Get("db"),
	}).ToString()

	_, err := psql.db.Exec(query)

	return err
}

//새 유저 등록 : C
func CreateNewUser(user User) error {

	existUser, _ := FindUserByUserId(user.UserId.String)
	if existUser != nil {
		return errors.New("Conflict")
	}

	//해쉬 비밀번호 생성
	hashed, _ := hashPassword.GenerateHashPassword(user.Password)

	//해쉬 비밀번호 생성된 데이터를 넣는다.
	query := qb.Insert("user", "user_id, user_type, password").Value([]string{
		user.UserId.String,
		user.UserType,
		hashed,
	}).ToString()

	//toDo : 새 유저를 생성하는 쿼리 필요
	_, err := psql.db.Exec(query)

	return err
}

//해쉬 패스워드 검증
func (u *User) ValidatePassword(pw string) (bool, error) {
	valid, err := hashPassword.CompareHashPassword(u.Password, pw)
	if err != nil {
		return false, err
	}
	return valid, nil
}

//아래부턴 쿼리
func FindUserById(id uint) (*User, error) {
	user := &User{}

	query := qb.Select("id, user_id, user_type, password").From("user").Where("id", id).ToString()
	psql.db.QueryRow(query).Scan(
		&user.ID, &user.UserId, &user.UserType, &user.Password,
	)

	if !user.UserId.Valid {
		err := errors.New("Not Found")
		return nil, err
	}

	return user, nil
}

func FindUserByUserId(userId string) (*User, error) {
	user := &User{}

	query := qb.Select("id, user_id, user_type, password").From("user").Where("user_id", userId).ToString()
	psql.db.QueryRow(query).Scan(
		&user.ID, &user.UserId, &user.UserType, &user.Password,
	)

	if !user.UserId.Valid {
		err := errors.New("Not Found")
		return nil, err
	}

	return user, nil
}
