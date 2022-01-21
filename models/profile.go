package models

import (
	"errors"
	"strings"

	getTag "github.com/Park-Kwonsoo/moving-server/pkg/get-struct-info"
	qb "github.com/Park-Kwonsoo/moving-server/pkg/query-builder"
)

type Profile struct {
	baseType
	User         User   `db:"user_user_id varchar(255) references dev.user(user_id)"`
	Name         string `db:"name varchar(255)"`
	Birth        string `db:"birth varchar(10)"`
	Gender       string `db:"gender varchar(6)"`
	ProfileImage string `db:"profile_img varchar(2000)"`
}

func profileMigrate() error {

	column := make([]string, 0)
	column = append(column, strings.Join(getCreatedTableColumn(), ", "))
	column = append(column, getTag.GetStructInfoByTag("db", &Profile{})...)

	query := qb.CreateTable("profile").TableComlumn(
		column...,
	).ToString()

	_, err := psql.db.Exec(query)

	return err
}

//새 프로필 생성
func CreateNewProfile(profile *Profile) error {

	existProfile, _ := FindProfileByUserUserId(profile.User.UserId.String)
	if existProfile != nil {
		return errors.New("Conflict")
	}

	query := qb.Insert("profile", "name, birth, gender, profile_img, user_user_id").Value(
		profile.Name,
		profile.Birth,
		profile.Gender,
		profile.ProfileImage,
		profile.User.UserId.String,
	).ToString()

	//toDo : make profile -> join with User
	_, err := psql.db.Exec(query)

	return err
}

func FindProfileById(id uint) (*Profile, error) {

	profile := &Profile{}
	var userId string //유저 아이디 저장

	query := qb.Select("id, creatd_at, updated_at, name, birth, gender, profile_img, user_user_id").From("profile").Where("id", id).ToString()
	psql.db.QueryRow(query).Scan(
		&profile.ID, &profile.CreatedAt, &profile.UpdatedAt, &profile.Name, &profile.Birth, &profile.Gender, &profile.ProfileImage, &userId,
	)

	user, _ := FindUserByUserId(userId)
	profile.User = *user

	if !profile.User.UserId.Valid {
		err := errors.New("Not Found")
		return nil, err
	}

	return profile, nil
}

func FindProfileByUserUserId(userId string) (*Profile, error) {

	user, err := FindUserByUserId(userId)
	if err != nil {
		return nil, err
	}
	if !user.UserId.Valid {
		err := errors.New("Not Found")
		return nil, err
	}

	profile := &Profile{}
	profile.User = *user

	query := qb.Select("id, created_at, updated_at, name, birth, gender, profile_img").From("profile").Where("user_user_id", userId).ToString()
	psql.db.QueryRow(query).Scan(
		&profile.ID, &profile.CreatedAt, &profile.UpdatedAt, &profile.Name, &profile.Birth, &profile.Gender, &profile.ProfileImage,
	)

	if profile.ID == 0 {
		err := errors.New("Not Found")
		return nil, err
	}

	return profile, nil
}
