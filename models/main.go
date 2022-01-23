package models

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

type dbms struct {
	port     uint
	host     string
	user     string
	password string
	dbName   string
	schema   string
	db       *sql.DB
}

var psql dbms

//db 실행 시 migrate
func migrate() error {

	//update Trigger 등록
	err := createUpdateFunction()
	errHandler.PanicErr(err)

	//유저 정보 migrate
	err = memberMigrate()
	errHandler.PanicErr(err)

	//Profile 정보 migrate
	err = profileMigrate()
	errHandler.PanicErr(err)

	//음악 정보 migrate
	err = musicMigrate()
	errHandler.PanicErr(err)

	//플레이리스트 migrate
	err = playlistMigrate()
	errHandler.PanicErr(err)

	return nil
}

/**
*	mapping : many2many 태그를 통해 테이블을 생성함
 */
func tableMapping(s interface{}) error {

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
			c = append(c, strings.Join(getCreatedTableColumn(), ", "))
			c = append(c, c1.String(), c2.String())

			query := qb.CreateTable(b.String()).TableComlumn(
				c...,
			).ToString()

			_, err := psql.Exec(query)
			errHandler.PanicErr(err)
		}
	}

	return nil
}

/**
*	dbms : Exec
 */
func (db *dbms) Exec(query string, args ...interface{}) (sql.Result, error) {
	//toDO : kogging
	// log.Println(query)
	return db.db.Exec(query, args...)
}

/**
*	dbms : QueryRow
 */
func (db *dbms) QueryRow(query string, args ...interface{}) *sql.Row {
	//toDo : logging
	// log.Println(query)
	return db.db.QueryRow(query, args...)
}

/**
*	dbms : Query
 */
func (db *dbms) Query(query string, args ...interface{}) (*sql.Rows, error) {
	//toDo : logging
	// log.Println(query)
	return db.db.Query(query, args...)
}

/*
* connect to PsqlDB
 */
func Connect() {

	//env 파일 로딩
	e := godotenv.Load(".env")
	errHandler.PanicErr(e)

	//psql db
	psql.host = "127.0.0.1"
	psql.port = 5432
	psql.user = os.Getenv("DB_USER")
	psql.password = os.Getenv("DB_PASSWORD")
	psql.dbName = os.Getenv("DB_NAME")
	psql.schema = os.Getenv("DB_SCHEMA")

	//postgresql 환경을 설정
	psqlconn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Seoul search_path=%s",
		psql.host, psql.port, psql.user, psql.password, psql.dbName, psql.schema)

	var err error
	psql.db, err = sql.Open("postgres", psqlconn)

	errHandler.PanicErr(err)
	log.Println("db Connected!")

	//db 연결 후 migrate
	migrate()
}

/*
* disconnect db
 */
func Disconnect() {
	psql.db.Close()
	log.Println("db disconnect!")
}
