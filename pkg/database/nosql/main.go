package nosql

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	errHandler "github.com/Park-Kwonsoo/moving-server/pkg/err-handler"
)

var MongoDB *mongo.Database

func mongoConn() *mongo.Client {

	e := godotenv.Load(".env")
	errHandler.PanicErr(e)

	mongoUrl := os.Getenv("MONGO_URL")
	opt := options.Client().ApplyURI(mongoUrl)

	client, err := mongo.Connect(context.TODO(), opt)
	errHandler.PanicErr(err)

	//connection check
	err = client.Ping(context.TODO(), nil)
	errHandler.PanicErr(err)

	log.Println("Mongo DB Connected!")
	return client
}

func Connect() {
	conn := mongoConn()
	MongoDB = conn.Database(os.Getenv("MONGO_DATABASE"))
}

func GetCollection(collection string) *mongo.Collection {
	return MongoDB.Collection(collection)
}
