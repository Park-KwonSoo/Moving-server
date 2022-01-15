package models

import (
	"reflect"
	"time"
)

type baseType struct {
	ID        uint      `model:"int identity(1,1) primary key"`
	CreatedAt time.Time `model:"timestamp not null default now()"`
	UpdatedAt time.Time `model:"timestamp"`
	DeletedAt uint      `model:"index"`
}

//BaseType의 create table시 사용할 컬럼 옵션 리턴
func getCreatedTableColumn() []string {
	b := baseType{}

	id, _ := reflect.TypeOf(b).FieldByName("ID")
	createdAt, _ := reflect.TypeOf(b).FieldByName("CreatedAt")
	updatedAt, _ := reflect.TypeOf(b).FieldByName("UpdatedAt")
	deletedAt, _ := reflect.TypeOf(b).FieldByName("DeletedAt")

	return []string{
		"id " + id.Tag.Get("model"),
		"createdAt " + createdAt.Tag.Get("model"),
		"updatedAt " + updatedAt.Tag.Get("model"),
		"deletedAt " + deletedAt.Tag.Get("model"),
	}
}
