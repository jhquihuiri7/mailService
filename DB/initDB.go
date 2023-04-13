package DB

import (
	"context"
	"database/sql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var (
	SQliteDB *sql.DB
	Mongodb  *mongo.Database
	err      error
	MailDB   *mongo.Collection
)

func InitSQLite() {
	SQliteDB, err = sql.Open("sqlite3", "./data/clients.db")
	if err != nil {
		log.Fatal(err)
	}
}
func InitMongo() {
	options := options.Client().ApplyURI("mongodb+srv://doadmin:Z3d87ni4E91g05aX@logiciel-applab-dab57134.mongo.ondigitalocean.com/admin?authSource=admin&replicaSet=logiciel-applab&tls=true&tlsCAFile=DB/ca-certificate.crt")
	client, err := mongo.Connect(context.TODO(), options)
	if err != nil {
		log.Fatal(err)
	}
	Mongodb = client.Database("LogicielDB")
}
func InitColletion() {
	MailDB = Mongodb.Collection("MailDB")
}
