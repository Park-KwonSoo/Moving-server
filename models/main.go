package models

import (
	"fmt"
	"os"

	"database/sql"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

type dbms struct {
	port     uint
	host     string
	user     string
	password string
	dbName   string
	db       *sql.DB
}

var psql dbms

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func migrate() bool {
	userMigrate() //유저 정보 migrate

	return true
}

func GetDB() dbms {
	return psql
}

func Connect() {

	//env 파일 로딩
	e := godotenv.Load(".env")
	checkError(e)

	//psql db
	psql.host = "localhost"
	psql.port = 5432
	psql.user = os.Getenv("USER")
	psql.password = os.Getenv("PASSWORD")
	psql.dbName = os.Getenv("DB_NAME")

	//postgresql 환경을 설정
	psqlconn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Seoul",
		psql.host, psql.port, psql.user, psql.password, psql.dbName)

	//gorm을 이용해 postgres를 오픈한다. schema를 지정하는 옵션은 gorm.Config에 존재
	// dbSchema := os.Getenv("DB_SCHEMA")

	//글로벌 변수를 초기화 하기 위해서는 := 대신 = 를 사용해야 한다!
	var err error
	psql.db, err = sql.Open("postgres", psqlconn)
	checkError(err)
	fmt.Println("db Connected!")

	migrate()
}
