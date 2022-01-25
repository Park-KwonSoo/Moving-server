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
	QUERY_TYPE_1       string //query 1
	QUERY_TYPE_1_VALUE string //query 1 value

	QUERY_TYPE_2       string //query 2
	QUERY_TYPE_2_VALUE string //query 2 value

	QUERY_TYPE_3       string
	QUERY_TYPE_3_VALUE string

	QUERY_TYPE_4        string //query 4
	QUERY_TYPE_4_OPTION string //query 4_1
	QUERY_TYPE_4_VALUE  string //query 4_1_2 value

	QUERY_TYPE_5        string
	QUERY_TYPE_5_OPTION string
	QUERY_TYPE_5_VALUE  string
}

//argument를 받아서 새로운 table을 생성한다.
func CreateTable(t string) *query {

	q := &query{
		QUERY_TYPE_1:        CREATE_TABLE,
		QUERY_TYPE_1_VALUE:  "",
		QUERY_TYPE_2:        t,
		QUERY_TYPE_2_VALUE:  "",
		QUERY_TYPE_3:        "",
		QUERY_TYPE_3_VALUE:  "",
		QUERY_TYPE_4:        "",
		QUERY_TYPE_4_OPTION: "",
		QUERY_TYPE_4_VALUE:  "",
		QUERY_TYPE_5:        "",
		QUERY_TYPE_5_OPTION: "",
		QUERY_TYPE_5_VALUE:  "",
	}

	return q
}

//table의 column 정보로 쿼리 생성
func (q *query) TableComlumn(opt ...string) *query {

	q.QUERY_TYPE_2_VALUE = fmt.Sprintf("(%s)", strings.Join(opt, ", "))

	return q
}

//Select할 column들 argument로 받는다.
func Select(c string) *query {

	q := &query{
		QUERY_TYPE_1:        SELECT,
		QUERY_TYPE_1_VALUE:  c,
		QUERY_TYPE_2:        "",
		QUERY_TYPE_2_VALUE:  "",
		QUERY_TYPE_3:        "",
		QUERY_TYPE_3_VALUE:  "",
		QUERY_TYPE_4:        "",
		QUERY_TYPE_4_OPTION: "",
		QUERY_TYPE_4_VALUE:  "",
		QUERY_TYPE_5:        "",
		QUERY_TYPE_5_OPTION: "",
		QUERY_TYPE_5_VALUE:  "",
	}

	return q
}

//Table을 argument로 받는다.
func (q *query) From(t string) *query {

	q.QUERY_TYPE_2 = FROM
	q.QUERY_TYPE_2_VALUE = t

	return q

}

//Join
func (q *query) Join(t_key string, t2 string, t2_key string) *query {

	q.QUERY_TYPE_3 = JOIN

	var b bytes.Buffer
	b.WriteString(t2)
	b.WriteString(" ")
	b.WriteString("ON")
	b.WriteString(" ")
	b.WriteString(q.QUERY_TYPE_2_VALUE)
	b.WriteString(".")
	b.WriteString(t_key)
	b.WriteString("=")
	b.WriteString(t2)
	b.WriteString(".")
	b.WriteString(t2_key)

	q.QUERY_TYPE_3_VALUE = b.String()

	return q
}

//where절의 column과 value를 argument로 받는다
func (q *query) Where(c string, v interface{}) *query {

	q.QUERY_TYPE_4 = WHERE
	q.QUERY_TYPE_4_OPTION = c
	q.QUERY_TYPE_4_VALUE = fmt.Sprintf("='%v'", v)

	return q
}

//And절 : where절 추가
func (q *query) And(c string, v interface{}) *query {

	q.QUERY_TYPE_5 = AND
	q.QUERY_TYPE_5_OPTION = c
	q.QUERY_TYPE_5_VALUE = fmt.Sprintf("='%v'", v)

	return q
}

//argument로 table과 colume을 받는다.
func Insert(t string, c string) *query {

	var b bytes.Buffer
	b.WriteString("(")
	b.WriteString(c)
	b.WriteString(")")

	q := &query{
		QUERY_TYPE_1:        INSERT,
		QUERY_TYPE_1_VALUE:  "",
		QUERY_TYPE_2:        t,
		QUERY_TYPE_2_VALUE:  b.String(),
		QUERY_TYPE_3:        "",
		QUERY_TYPE_3_VALUE:  "",
		QUERY_TYPE_4:        "",
		QUERY_TYPE_4_OPTION: "",
		QUERY_TYPE_4_VALUE:  "",
		QUERY_TYPE_5:        "",
		QUERY_TYPE_5_OPTION: "",
		QUERY_TYPE_5_VALUE:  "",
	}

	return q
}

//추가할 데이터의 column에 대한 value 값을 받는다.
func (q *query) Value(v ...interface{}) *query {

	q.QUERY_TYPE_4 = VALUE

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

	q.QUERY_TYPE_4_VALUE = b.String()

	return q
}

//행을 업데이트 하는 쿼리 : table
func Update(t string) *query {

	q := &query{
		QUERY_TYPE_1:        UPDATE,
		QUERY_TYPE_1_VALUE:  t,
		QUERY_TYPE_2:        "",
		QUERY_TYPE_2_VALUE:  "",
		QUERY_TYPE_3:        "",
		QUERY_TYPE_3_VALUE:  "",
		QUERY_TYPE_4:        "",
		QUERY_TYPE_4_OPTION: "",
		QUERY_TYPE_4_VALUE:  "",
		QUERY_TYPE_5:        "",
		QUERY_TYPE_5_OPTION: "",
		QUERY_TYPE_5_VALUE:  "",
	}

	return q

}

func (q *query) Set(c string, v []string) *query {

	column := strings.Split(c, ", ")

	if len(column) != len(v) {
		return nil
	}

	q.QUERY_TYPE_2 = SET

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

	q.QUERY_TYPE_2_VALUE = b.String()

	return q
}

//마지막으로 string으로 변환하는 receiver
func (q *query) ToString() string {
	var b bytes.Buffer
	b.WriteString(q.QUERY_TYPE_1)
	b.WriteString(" ")
	b.WriteString(q.QUERY_TYPE_1_VALUE)
	b.WriteString(" ")

	b.WriteString(q.QUERY_TYPE_2)
	b.WriteString(" ")
	b.WriteString(q.QUERY_TYPE_2_VALUE)
	b.WriteString(" ")

	b.WriteString(q.QUERY_TYPE_3)
	b.WriteString(" ")
	b.WriteString(q.QUERY_TYPE_3_VALUE)
	b.WriteString(" ")

	b.WriteString(q.QUERY_TYPE_4)
	b.WriteString(" ")
	b.WriteString(q.QUERY_TYPE_4_OPTION)
	b.WriteString(" ")
	b.WriteString(q.QUERY_TYPE_4_VALUE)
	b.WriteString(" ")

	b.WriteString(q.QUERY_TYPE_5)
	b.WriteString(" ")
	b.WriteString(q.QUERY_TYPE_5_OPTION)
	b.WriteString(" ")
	b.WriteString(q.QUERY_TYPE_5_VALUE)
	b.WriteString(" ")

	return b.String()
}
