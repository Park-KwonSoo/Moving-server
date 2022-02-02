package querybuilder

import (
	"bytes"
	"fmt"
	"strings"
)

const (
	SELECT       = "SELECT"
	INSERT       = "INSERT INTO"
	UPDATE       = "UPDATE"
	CREATE_TABLE = "CREATE TABLE IF NOT EXISTS"
	ALTER_TABLE  = "ALTER TABLE"

	SET  = "SET"
	FROM = "FROM"

	JOIN = "JOIN"

	WHERE = "WHERE"
	VALUE = "VALUES"

	AND = "AND"
)

type query struct {
	bytes.Buffer
}

//argument를 받아서 새로운 table을 생성한다.
func CreateTable(t string) *query {

	q := &query{}

	q.WriteString(CREATE_TABLE)
	q.WriteString(" ")
	q.WriteString(t)
	q.WriteString(" ")

	return q
}

//table의 column 정보로 쿼리 생성
func (q *query) TableComlumn(opt ...string) *query {

	q.WriteString("(")
	q.WriteString(strings.Join(opt, ", "))
	q.WriteString(") ")

	return q
}

//Select할 column들 argument로 받는다.
func Select(c string) *query {

	q := &query{}

	q.WriteString(SELECT)
	q.WriteString(" ")
	q.WriteString(c)
	q.WriteString(" ")

	return q
}

//Table을 argument로 받는다.
func (q *query) From(t string) *query {

	q.WriteString(FROM)
	q.WriteString(" ")
	q.WriteString(t)
	q.WriteString(" ")

	return q

}

//Join
func (q *query) Join(t string, t_key string, t2 string, t2_key string) *query {

	q.WriteString(JOIN)
	q.WriteString(" ")

	var b bytes.Buffer
	b.WriteString(t2)
	b.WriteString(" ")
	b.WriteString("ON")
	b.WriteString(" ")
	b.WriteString(t)
	b.WriteString(".")
	b.WriteString(t_key)
	b.WriteString("=")
	b.WriteString(t2)
	b.WriteString(".")
	b.WriteString(t2_key)

	q.WriteString(b.String())
	q.WriteString(" ")

	return q
}

//where절의 column과 value를 argument로 받는다
func (q *query) Where(c string, v interface{}) *query {

	q.WriteString(WHERE)
	q.WriteString(" ")
	q.WriteString(c)
	q.WriteString(" ")
	q.WriteString("='")
	q.WriteString(fmt.Sprintf("%v", v))
	q.WriteString("' ")

	return q
}

//And절 : where절 추가
func (q *query) And(c string, v interface{}) *query {

	q.WriteString(AND)
	q.WriteString(" ")
	q.WriteString(c)
	q.WriteString(" ")
	q.WriteString("='")
	q.WriteString(fmt.Sprintf("%v", v))
	q.WriteString("' ")

	return q
}

//argument로 table과 colume을 받는다.
func Insert(t string, c string) *query {

	q := &query{}

	q.WriteString(INSERT)
	q.WriteString(" ")
	q.WriteString(t)
	q.WriteString(" ")
	q.WriteString("(")
	q.WriteString(c)
	q.WriteString(") ")

	return q
}

//추가할 데이터의 column에 대한 value 값을 받는다.
func (q *query) Value(v ...interface{}) *query {

	q.WriteString(VALUE)
	q.WriteString(" ")

	var b bytes.Buffer
	b.WriteString("(")

	r := make([]string, len(v))
	for i := 0; i < len(v); i++ {
		var t bytes.Buffer
		t.WriteString("'")
		t.WriteString(fmt.Sprintf("%v", v[i]))
		t.WriteString("'")

		r[i] = t.String()
	}

	b.WriteString(strings.Join(r, ", "))
	b.WriteString(")")

	b.WriteString(" ")
	b.WriteString("RETURNING id")

	q.WriteString(b.String())
	q.WriteString(" ")

	return q
}

//행을 업데이트 하는 쿼리 : table
func Update(t string) *query {

	q := &query{}

	q.WriteString(UPDATE)
	q.WriteString(" ")
	q.WriteString(t)
	q.WriteString(" ")

	return q

}

func (q *query) Set(c string, v []string) *query {

	column := strings.Split(c, ", ")

	if len(column) != len(v) {
		return nil
	}

	q.WriteString(SET)
	q.WriteString(" ")

	var b bytes.Buffer

	ql := len(column)
	for i := 0; i < ql-1; i++ {
		b.WriteString(column[i])
		b.WriteString("=")
		b.WriteString("'")
		b.WriteString(v[i])
		b.WriteString("'")
		b.WriteString(", ")
	}
	b.WriteString(column[ql-1])
	b.WriteString("=")
	b.WriteString("'")
	b.WriteString(v[ql-1])
	b.WriteString("'")

	q.WriteString(b.String())

	return q
}

//마지막으로 string으로 변환하는 receiver
func (q *query) ToString() string {
	return q.String()
}
