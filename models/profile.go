package models

import (
	"errors"
	"fmt"

	qb "github.com/Park-Kwonsoo/moving-server/pkg/query-builder"
)

type Profile struct {
	baseType
	User         User   `db:"user_id int references user(id)"`
	Name         string `db:"name varchar(255)"`
	Birth        string `db:"birth varchar(10)"`
	Gender       string `db:"gender varchar(6)"`
	ProfileImage string `db:"profile_img varchar(2000)"`
}

//새 프로필 생성
func CreateNewProfile(profile Profile) error {

	existProfile, _ := FindProfileByUserId(profile.User.UserId.String)
	if existProfile != nil {
		return errors.New("Conflict")
	}

	query := qb.Insert("profile", "name, birth, gender, profile_image").Value([]string{
		profile.Name,
		profile.Birth,
		profile.Gender,
		profile.ProfileImage,
	}).ToString()

	//toDo : make profile -> join with User
	fmt.Println(query)
	_, err := psql.db.Exec(query)
	fmt.Println(err)

	return err
}

func FindProfileById(id uint) (*Profile, error) {

	user := &User{}
	query := qb.Select("id, user_id, user_type").From("user").Where("profile_id", id).ToString()
	psql.db.QueryRow(query).Scan(
		&user.ID, &user.UserId, &user.UserType,
	)

	if !user.UserId.Valid {
		err := errors.New("Not Found")
		return nil, err
	}

	profile := &Profile{}
	profile.User = *user

	query = qb.Select("id, user, name, birth, gender, profile_image").From("profile").Where("id", id).ToString()
	psql.db.QueryRow(query).Scan(
		&profile.ID, &profile.User, &profile.Name, &profile.Birth, &profile.Gender, &profile.ProfileImage,
	)

	if !profile.User.UserId.Valid {
		err := errors.New("Not Found")
		return nil, err
	}

	return profile, nil
}

func FindProfileByUserId(userId string) (*Profile, error) {

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

	query := qb.Select("id, name, birth, gender, profile_image").From("profile").Where("user_id", userId).ToString()
	psql.db.QueryRow(query).Scan(
		&profile.ID, &profile.Name, &profile.Birth, &profile.Gender, &profile.ProfileImage,
	)

	if !profile.User.UserId.Valid {
		err := errors.New("Not Found")
		return nil, err
	}

	return profile, nil
}
