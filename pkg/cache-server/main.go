package cacheserver

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	redis "github.com/go-redis/redis"

	errHandler "github.com/Park-Kwonsoo/moving-server/pkg/err-handler"
)

var Redis *redis.Client

//Cache를 관리할 Redis DB : Connect
func Connect() error {
	e := godotenv.Load(".env")
	errHandler.PanicErr(e)

	//1번 DB에 연결함
	Redis = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       1,
	})

	_, err := Redis.Ping().Result()
	if err != nil {
		return err
	}

	log.Println("Redis Connected!")

	return nil
}

//Cache Set : key, value => error
func SetCache(key string, val interface{}) error {
	expiration, err := strconv.Atoi(os.Getenv("REDIS_DB_EXPIRATION_HOUR"))
	if err != nil {
		return err
	}

	return Redis.Set(key, val, time.Hour*time.Duration(expiration)).Err()
}

//Cache Set For gRPC : key, value => error
func SetCacheProto(key string, val interface{}, e error) error {

	if e != nil {
		return e
	}

	expiration, err := strconv.Atoi(os.Getenv("REDIS_DB_EXPIRATION_HOUR"))
	if err != nil {
		return err
	}

	resp, _ := val.(protoreflect.ProtoMessage)
	marshal, err := proto.Marshal(resp)
	if err != nil {
		return err
	}

	return Redis.Set(key, marshal, time.Hour*time.Duration(expiration)).Err()
}

//Cache Get : key => value, error
func GetCache(key string) ([]byte, error) {
	val, err := Redis.Get(key).Bytes()
	if err == redis.Nil {
		return val, nil
	}
	return val, err
}

//Cache Get For gRPC : key, value => error
func GetCacheProto(key string, dst protoreflect.ProtoMessage) (protoreflect.ProtoMessage, error) {
	val, err := GetCache(key)
	if err == nil {
		err = proto.Unmarshal(val, dst)
		return dst, err
	}
	return nil, err
}

//Cache Remove : key => error
func DeleteCache(key string) error {
	return Redis.Del(key).Err()
}

/**
*	Cache Key를 만들기 위한 값들을 인자로 받아 key값을 리턴
* 	Parameter : service Name, Member Id, Request Message
 */
func KeyMake(model string, id interface{}) string {
	var b bytes.Buffer
	b.WriteString(model)
	b.WriteString(":")
	b.WriteString(fmt.Sprintf("%v", id))

	return b.String()
}
