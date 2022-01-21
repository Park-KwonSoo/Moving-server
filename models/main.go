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

func migrate() error {

	//유저 정보 migrate
	err := memberMigrate()
	if err != nil {
		fmt.Println(err)
		return err
	}

	//Profile 정보 migrate
	err = profileMigrate()
	if err != nil {
		fmt.Println(err)
		return err
	}

	//멤버 플레이리스트 생성
	err = memberPlaylistMigrate()
	if err != nil {
		fmt.Println(err)
		return err
	}

	//음악 정보 migrate
	err = musicMigrate()
	if err != nil {
		fmt.Println(err)
		return err
	}

	//플레이리스트 migrate
	err = playlistMigrate()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

//return dbms
func GetDB() dbms {
	return psql
}

/*
* connect to PsqlDB
 */
func Connect() {

	//env 파일 로딩
	e := godotenv.Load(".env")
	checkError(e)

	//psql db
	psql.host = "localhost"
	psql.port = 5432
	psql.user = os.Getenv("DB_USER")
	psql.password = os.Getenv("DB_PASSWORD")
	psql.dbName = os.Getenv("DB_NAME")

	//postgresql 환경을 설정
	psqlconn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Seoul search_path=%s",
		psql.host, psql.port, psql.user, psql.password, psql.dbName, os.Getenv("DB_SCHEMA"))

	var err error
	psql.db, err = sql.Open("postgres", psqlconn)

	checkError(err)
	fmt.Println("db Connected!")

	migrate()
}
