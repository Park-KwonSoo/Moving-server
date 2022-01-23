package models

import (
	"errors"
	"strings"

	"database/sql"

	getTag "github.com/Park-Kwonsoo/moving-server/pkg/get-struct-info"
	hashPassword "github.com/Park-Kwonsoo/moving-server/pkg/hashing-password"
	qb "github.com/Park-Kwonsoo/moving-server/pkg/query-builder"
)

type Member struct {
	baseType
	MemId    sql.NullString `db:"mem_id varchar(255) unique not null"`
	MemType  string         `db:"mem_type varchar(10) not null"`
	Password string         `db:"password varchar(2000) not null"`
}

//member migrate
func memberMigrate() error {

	column := make([]string, 0)
	column = append(column, strings.Join(getCreatedTableColumn(), ", "))
	column = append(column, getTag.GetStructInfoByTag("db", &Member{})...)

	query := qb.CreateTable("member").TableComlumn(
		column...,
	).ToString()

	_, err := psql.Exec(query)
	if err != nil {
		return err
	}

	err = createUpdateTrigger("member")
	if err != nil {
		return err
	}

	err = tableMapping(&Member{})
	return err
}

//새 유저 등록 : C
func CreateNewMember(member *Member) error {

	existMember, _ := FindOneMemberByMemId(member.MemId.String)
	if existMember != nil {
		return errors.New("Conflict")
	}

	//해쉬 비밀번호 생성
	hashed, _ := hashPassword.GenerateHashPassword(member.Password)

	//해쉬 비밀번호 생성된 데이터를 넣는다.
	query := qb.Insert("member", "mem_id, mem_type, password").Value(
		member.MemId.String,
		member.MemType,
		hashed,
	).ToString()

	err := psql.QueryRow(query).Scan(&member.ID)

	return err
}

//해쉬 패스워드 검증
func (m *Member) ValidatePassword(pw string) (bool, error) {
	valid, err := hashPassword.CompareHashPassword(m.Password, pw)
	if err != nil {
		return false, err
	}
	return valid, nil
}

//패스워드 변경
func (m *Member) ChangePassword(pw string) error {

	hashed, _ := hashPassword.GenerateHashPassword(pw)
	m.Password = hashed

	query := qb.Update("member").Set("password", []string{
		m.Password,
	}).Where("id", m.ID).ToString()

	_, err := psql.Exec(query)

	return err
}

//아래부턴 쿼리
func FindOneMemberById(id uint) (*Member, error) {
	member := &Member{}

	query := qb.Select("id, created_at, updated_at, mem_id, mem_type, password").From("member").Where("id", id).ToString()
	psql.QueryRow(query).Scan(
		&member.ID, &member.CreatedAt, &member.UpdatedAt, &member.MemId, &member.MemType, &member.Password,
	)

	if !member.MemId.Valid {
		err := errors.New("Not Found")
		return nil, err
	}

	return member, nil
}

func FindOneMemberByMemId(memId string) (*Member, error) {
	member := &Member{}

	query := qb.Select("id, created_at, updated_at, mem_id, mem_type, password").From("member").Where("mem_id", memId).ToString()
	psql.QueryRow(query).Scan(
		&member.ID, &member.CreatedAt, &member.UpdatedAt, &member.MemId, &member.MemType, &member.Password,
	)

	if !member.MemId.Valid {
		err := errors.New("Not Found")
		return nil, err
	}

	return member, nil
}
