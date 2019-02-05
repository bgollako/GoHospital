package db

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/mongodb/mongo-go-driver/x/network/connstring"
	"time"
)

var DBClient *mongoClient

type mongoClient struct {
	client mongo.Client
}

func init() {
	opts := &options.ClientOptions{
		ConnString: connstring.ConnString{
			Username: "mongoadmin",
			Password: "mongopassword",
			Original: "mongodb://db:27017",
		},
	}
	customClient, err := mongo.NewClientWithOptions("mongodb://db:27017", opts)
	if err != nil {
		fmt.Println("Error while initializing mongo client " + err.Error())
		return
	}
	fmt.Println("Successfully created Mongo client")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = customClient.Connect(ctx)
	if err != nil {
		fmt.Println("Unable to connect to client " + err.Error())
		return
	}
	fmt.Println("Successfully connected to Mongo database")
	DBClient = &mongoClient{}
	DBClient.client = *customClient
}

func TestConnection(ctx context.Context) error {
	err := DBClient.client.Ping(ctx, nil)
	return err
}

func InitialiseDummyPatients(ctx context.Context) error {
	p1 := Patient{Name: "Ash", Age: 20, Disease: "Cough"}
	p2 := Patient{Name: "Misty", Age: 22, Disease: "Cold"}
	p3 := Patient{Name: "Brock", Age: 34, Disease: "Fever"}
	p4 := Patient{Name: "Rahul", Age: 44, Disease: "Rashes"}
	p5 := Patient{Name: "Raj", Age: 55, Disease: "Measles"}
	documents := []interface{}{p1, p2, p3, p4, p5}
	_, err := DBClient.client.Database("Test").Collection("Patients").InsertMany(ctx, documents)
	return err
}

func GetAllPatients(ctx context.Context) ([]Patient, error)  {
	patients := make([]Patient, 0)
	c, err := DBClient.client.Database("Test").Collection("Patients").Find(ctx, bson.M{})
	if err != nil{
		return nil, err
	}
	for c.Next(ctx) {
		p := Patient{}
		c.Decode(&p)
		patients = append(patients, p)
	}
	return patients, nil
}
