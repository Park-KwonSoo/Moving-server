package models

import (
	qb "github.com/Park-Kwonsoo/moving-server/pkg/query-builder"
)

type Profile struct {
	baseType
	User         User
	Name         string
	Birth        string
	Gender       string
	ProfileImage string
}

func FindProfileById(id uint) *Profile {
	profile := &Profile{}

	query := qb.Select("user, name, birth, gender, profile_image").From("dev.profile").Where("id", id).ToString()
	psql.db.QueryRow(query).Scan(
		&profile.User, &profile.Name, &profile.Birth, &profile.Gender, &profile.ProfileImage)

	return profile
}

func FindProfileByUserId(userId string) *Profile {
	profile := &Profile{}

	query := qb.Select("user, name, birth, gender, profile_image").From("dev.profile").Where("user_id", userId).ToString()
	psql.db.QueryRow(query).Scan(
		&profile.User, &profile.Name, &profile.Birth, &profile.Gender, &profile.ProfileImage)

	return profile
}
