package querybuilder

import (
	"bytes"
	"fmt"
	"strings"
)

const (
	SELECT       = "SELECT"
	INSERT       = "INSERT INTO"
	CREATE_TABLE = "CREATE TABLE IF NOT EXISTS"
	ALTER_TABLE  = "ALTER TABLE"

	FROM = "FROM"

	WHERE = "WHERE"
	VALUE = "VALUES"
)

var SCHEMA string = "dev"

type query struct {
	QUERY_TYPE       string //select | insert | create
	QUERY_TYPE_VALUE string //어떤 것을 select 할지, insert 할지 등

	QUERY_TABLE       string //테이블 이름
	QUERY_TABLE_VALUE string //테이블 값

	QUERY_WHERE        string //조건절
	QUERY_WHERE_OPTION string //조건 옵션
	QUERY_WHERE_VALUE  string //조건 값
}

//argument를 받아서 새로운 table을 생성한다.
func CreateTable(t string) *query {

	var b bytes.Buffer
	b.WriteString(SCHEMA)
	b.WriteString(".")
	b.WriteString(t)

	q := &query{
		QUERY_TYPE:         CREATE_TABLE,
		QUERY_TYPE_VALUE:   "",
		QUERY_TABLE:        b.String(),
		QUERY_TABLE_VALUE:  "",
		QUERY_WHERE:        "",
		QUERY_WHERE_OPTION: "",
		QUERY_WHERE_VALUE:  "",
	}

	return q
}

//table의 column 정보로 쿼리 생성
func (q *query) TableComlumn(opt []string) *query {

	q.QUERY_TABLE_VALUE = fmt.Sprintf("(%s)", strings.Join(opt, ", "))

	return q
}

//Select할 column들 argument로 받는다.
func Select(c string) *query {

	q := &query{
		QUERY_TYPE:         SELECT,
		QUERY_TYPE_VALUE:   c,
		QUERY_TABLE:        "",
		QUERY_TABLE_VALUE:  "",
		QUERY_WHERE:        "",
		QUERY_WHERE_OPTION: "",
		QUERY_WHERE_VALUE:  "",
	}

	return q
}

//Table을 argument로 받는다.
func (q *query) From(t string) *query {

	q.QUERY_TABLE = FROM

	var b bytes.Buffer
	b.WriteString(SCHEMA)
	b.WriteString(".")
	b.WriteString(t)

	q.QUERY_TABLE_VALUE = b.String()

	return q

}

//where절의 column과 value를 argument로 받는다
func (q *query) Where(c string, v interface{}) *query {

	q.QUERY_WHERE = WHERE
	q.QUERY_WHERE_OPTION = c
	q.QUERY_WHERE_VALUE = fmt.Sprintf("='%v'", v)

	return q
}

//argument로 table과 colume을 받는다.
func Insert(t string, c string) *query {

	var table bytes.Buffer
	table.WriteString(SCHEMA)
	table.WriteString(".")
	table.WriteString(t)

	var b bytes.Buffer
	b.WriteString("(")
	b.WriteString(c)
	b.WriteString(")")

	q := &query{
		QUERY_TYPE:         INSERT,
		QUERY_TYPE_VALUE:   "",
		QUERY_TABLE:        table.String(),
		QUERY_TABLE_VALUE:  b.String(),
		QUERY_WHERE:        "",
		QUERY_WHERE_OPTION: "",
		QUERY_WHERE_VALUE:  "",
	}

	return q
}

//추가할 데이터의 column에 대한 value 값을 받는다.
func (q *query) Value(v []string) *query {

	q.QUERY_WHERE = VALUE

	var b bytes.Buffer
	b.WriteString("(")

	r := make([]string, len(v))
	for i := 0; i < len(v); i++ {
		var t bytes.Buffer
		t.WriteString("'")
		t.WriteString(v[i])
		t.WriteString("'")

		r[i] = t.String()
	}

	b.WriteString(strings.Join(r, ", "))
	b.WriteString(")")

	q.QUERY_WHERE_VALUE = b.String()

	return q
}

//마지막으로 string으로 변환하는 receiver
func (q query) ToString() string {
	var b bytes.Buffer
	b.WriteString(q.QUERY_TYPE)
	b.WriteString(" ")
	b.WriteString(q.QUERY_TYPE_VALUE)
	b.WriteString(" ")

	b.WriteString(q.QUERY_TABLE)
	b.WriteString(" ")
	b.WriteString(q.QUERY_TABLE_VALUE)
	b.WriteString(" ")

	b.WriteString(q.QUERY_WHERE)
	b.WriteString(" ")
	b.WriteString(q.QUERY_WHERE_OPTION)
	b.WriteString(" ")
	b.WriteString(q.QUERY_WHERE_VALUE)
	b.WriteString(" ")

	return b.String()
}
