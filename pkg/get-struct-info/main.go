package getstructinfo

import (
	"reflect"
)

//tag 이름과 struct 구조를 parameter로 받아 tag 정보를 리턴한다.
func GetStructInfoByTag(tag string, s interface{}) []string {

	e := reflect.ValueOf(s).Elem()
	fieldNum := e.NumField()

	rslt := make([]string, 0)

	for i := 0; i < fieldNum; i++ {
		t := e.Type().Field(i).Tag.Get(tag)
		if len(t) > 0 {
			rslt = append(rslt, t)
		}
	}

	return rslt
}
