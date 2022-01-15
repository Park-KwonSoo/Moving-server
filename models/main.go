package models

import (
	"fmt"
	"os"

	"database/sql"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

var db *sql.DB

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func GetDB() *sql.DB {
	return db
}

//db에 데이터를 추가한다.
func AddData(value interface{}) {
	// db.Create(value)
}

func Connect() {

	//env 파일 로딩
	e := godotenv.Load(".env")
	checkError(e)

	//psql host & name
	const (
		host = "localhost"
		port = 5432
	)
	//psql db
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	dbName := os.Getenv("DB_NAME")

	//postgresql 환경을 설정
	psqlconn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Seoul",
		host, port, user, password, dbName)

	//gorm을 이용해 postgres를 오픈한다. schema를 지정하는 옵션은 gorm.Config에 존재
	// dbSchema := os.Getenv("DB_SCHEMA")

	//글로벌 변수를 초기화 하기 위해서는 := 대신 = 를 사용해야 한다!
	var err error
	db, err = sql.Open("postgres", psqlconn)
	checkError(err)

	fmt.Println("db Connected!")

	//새로 테이블을 생성 : migrate한다. : 순서 중요

	// 	UserType: "LOCAL",
	// 	UserId:   "test312",
	// 	Password: "1234",
	// }
	// db.Create(&newUser)

	// db.Create(&Profile{
	// 	User:         newUser,
	// 	Name:         "박권수",
	// 	Birth:        "1996-07-18",
	// 	Gender:       "MALE",
	// 	ProfileImage: "NULL_IMAGE",
	// })
	// db.Create(&UserPlaylist{
	// 	User: newUser,
	// 	Playlist: []Playlist{
	// 		{PlaylistName: "가요", Music: []Music{}},
	// 		{PlaylistName: "팝송", Music: []Music{}},
	// 		{PlaylistName: "힙합", Music: []Music{}},
	// 	},
	// })

}
