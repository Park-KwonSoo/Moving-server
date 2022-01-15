package querybuilder

import (
	"fmt"
)

type query string

func Select(t string) query {
	s := query(fmt.Sprintf("SELECT %s ", t))
	return s
}

func (q query) From(t string) query {
	s := q + query(fmt.Sprintf("FROM %s ", t))
	return s
}

func (q query) Where(opt string, cond interface{}) query {
	s := q + query(fmt.Sprintf("WHERE %s='%s'", opt, fmt.Sprintf("%v", cond)))
	return s
}

func (q query) ToString() string {
	return string(q)
}
