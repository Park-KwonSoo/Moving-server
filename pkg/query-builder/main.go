package querybuilder

import (
	"bytes"
	"fmt"
	"strings"
)

type query string

//argument를 받아서 새로운 table을 생성한다.
func CreateTable(t string) query {
	var b bytes.Buffer
	b.WriteString("CREATE TABLE IF NOT EXISTS ")
	b.WriteString(t)
	b.WriteString(" ")

	s := query(b.String())
	return s
}

//table의 column 정보로 쿼리 생성
func (q query) TableOption(opt []string) query {
	var b bytes.Buffer
	b.WriteString(string(q))
	b.WriteString("(")
	b.WriteString(strings.Join(opt, ", "))
	b.WriteString(")")

	s := query(b.String())
	return s
}

//Select할 column들 argument로 받는다.
func Select(c string) query {
	var b bytes.Buffer
	b.WriteString("SELECT ")
	b.WriteString(c)
	b.WriteString(" ")

	s := query(b.String())
	return s
}

//Table을 argument로 받는다.
func (q query) From(t string) query {
	var b bytes.Buffer
	b.WriteString(string(q))
	b.WriteString("FROM ")
	b.WriteString(t)
	b.WriteString(" ")

	s := query(b.String())
	return s
}

//where절의 column과 value를 argument로 받는다
func (q query) Where(c string, v interface{}) query {
	var b bytes.Buffer
	b.WriteString(string(q))
	b.WriteString("WHERE ")
	b.WriteString(c)
	b.WriteString("=")
	b.WriteString("'")
	b.WriteString(fmt.Sprint(v))
	b.WriteString("'")

	s := query(b.String())
	return s
}

//argument로 table과 colume을 받는다.
func Insert(t string, c string) query {
	var b bytes.Buffer
	b.WriteString("INSERT INTO ")
	b.WriteString(t)
	b.WriteString(" ")
	b.WriteString("(")
	b.WriteString(c)
	b.WriteString(") ")

	s := query(b.String())
	return s
}

func (q query) Value(v []string) query {
	var b bytes.Buffer
	b.WriteString(string(q))
	b.WriteString("VALUE (")

	value := strings.Join(v, ", ")
	b.WriteString(value)
	b.WriteString(")")

	s := query(b.String())
	return s
}

//마지막으로 string으로 변환하는 receiver
func (q query) ToString() string {
	return string(q)
}
