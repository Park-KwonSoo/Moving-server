package sql

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"database/sql"

	_ "github.com/lib/pq"

	errHandler "github.com/Park-Kwonsoo/moving-server/pkg/err-handler"
	getTag "github.com/Park-Kwonsoo/moving-server/pkg/get-struct-info"
	qb "github.com/Park-Kwonsoo/moving-server/pkg/query-builder"
	"github.com/joho/godotenv"
)

type RDBMS struct {
	Port     uint
	Host     string
	User     string
	Password string
	DbName   string
	Schema   string
	Db       *sql.DB
}

var SQL RDBMS

/**
*	mapping : many2many 태그를 통해 테이블을 생성함
 */
func TableMapping(s interface{}) error {

	column := make([]string, 0)
	column = append(column, getTag.GetStructInfoByTag("mapping", s)...)
	for i := 0; i < len(column); i++ {
		if strings.HasPrefix(column[i], "many2many") {
			str := strings.Split(column[i], " ")
			sName := reflect.TypeOf(s).Elem().Name()

			var b bytes.Buffer
			b.WriteString(str[1])
			b.WriteString("_")
			b.WriteString(sName)

			var c1 bytes.Buffer
			c1.WriteString(str[1])
			c1.WriteString("_id integer references ")
			c1.WriteString(str[1])
			c1.WriteString("(id) on delete cascade")

			var c2 bytes.Buffer
			c2.WriteString(sName)
			c2.WriteString("_id integer references ")
			c2.WriteString(sName)
			c2.WriteString("(id) on delete cascade")

			c := make([]string, 0)
			c = append(c, strings.Join(GetCreatedTableColumn(), ", "))
			c = append(c, c1.String(), c2.String())

			query := qb.CreateTable(b.String()).TableComlumn(
				c...,
			).ToString()

			_, err := SQL.Exec(query)
			errHandler.PanicErr(err)
		}
	}

	return nil
}

/**
*	RDBMS : Exec
 */
func (db *RDBMS) Exec(query string, args ...interface{}) (sql.Result, error) {
	//toDO : kogging
	// log.Println(query)
	return db.Db.Exec(query, args...)
}

/**
*	RDBMS : QueryRow
 */
func (db *RDBMS) QueryRow(query string, args ...interface{}) *sql.Row {
	//toDo : logging
	// log.Println(query)
	return db.Db.QueryRow(query, args...)
}

/**
*	RDBMS : Query
 */
func (db *RDBMS) Query(query string, args ...interface{}) (*sql.Rows, error) {
	//toDo : logging
	// log.Println(query)
	return db.Db.Query(query, args...)
}

/*
* disconnect db
 */
func Disconnect() {
	SQL.Db.Close()
	log.Println("PostgreSQL disconnect!")
}

/*
* connect to SQLDB
 */
func init() {

	//env 파일 로딩
	e := godotenv.Load(".env")
	errHandler.PanicErr(e)

	//SQL db
	SQL.Host = "127.0.0.1"
	SQL.Port = 5432
	SQL.User = os.Getenv("SQL_DB_USER")
	SQL.Password = os.Getenv("SQL_DB_PASSWORD")
	SQL.DbName = os.Getenv("SQL_DB_NAME")
	SQL.Schema = os.Getenv("SQL_DB_SCHEMA")

	//postgresql 환경을 설정
	SQLconn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Seoul search_path=%s",
		SQL.Host, SQL.Port, SQL.User, SQL.Password, SQL.DbName, SQL.Schema)

	var err error
	SQL.Db, err = sql.Open("postgres", SQLconn)

	errHandler.PanicErr(err)
	log.Println("PostgreSQL Connected!")

	//update Trigger 등록
	err = createUpdateFunction()
	errHandler.PanicErr(err)
}
