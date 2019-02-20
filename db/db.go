package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/mongodb/mongo-go-driver/x/network/connstring"
	"reflect"
	"strings"
	"time"
)

var DBClient *mongoClient

const (
	DB_TEST  = "Test"
	COLLECTION_PATIENTS = "Patients"
	ID = "_id"
)

type mongoClient struct {
	client mongo.Client
}

func init() {
	opts := &options.ClientOptions{
		ConnString: connstring.ConnString{
			Username: "mongoadmin",
			Password: "mongopassword",
			Original: "mongodb://localhost:27017",
		},
	}
	customClient, err := mongo.NewClientWithOptions("mongodb://localhost:27017", opts)
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

func ShutDownClient(ctx context.Context) {
	err := DBClient.client.Disconnect(ctx)
	if err != nil {
		fmt.Println("Error while shutting down mongo client", err)
	}
}
func InitialiseDummyPatients(ctx context.Context) error {
	p1 := Patient{Name: "Ash", Age: 20, Disease: "Cough"}
	p2 := Patient{Name: "Misty", Age: 22, Disease: "Cold"}
	p3 := Patient{Name: "Brock", Age: 34, Disease: "Fever"}
	p4 := Patient{Name: "Rahul", Age: 44, Disease: "Rashes"}
	p5 := Patient{Name: "Raj", Age: 55, Disease: "Measles"}
	documents := []interface{}{p1, p2, p3, p4, p5}
	_, err := DBClient.client.Database(DB_TEST).Collection(COLLECTION_PATIENTS).InsertMany(ctx, documents)
	return err
}

func GetAllPatients(ctx context.Context) ([]Patient, error)  {
	patients := make([]Patient, 0)
	c, err := DBClient.client.Database(DB_TEST).Collection(COLLECTION_PATIENTS).Find(ctx, bson.M{})
	if err != nil{
		return nil, err
	}
	for c.Next(ctx) {
		p := Patient{}
		err = c.Decode(&p)
		if err != nil{
			fmt.Println("Error while decoding patients", err)
		}
		patients = append(patients, p)
	}
	return patients, nil
}

func PostPatient(ctx context.Context, patients []interface{}) error  {
	err := checkPatients(patients)
	if err != nil{
		fmt.Println("Error while parsing patients", err)
		return err
	}
	cursor, err := DBClient.client.Database(DB_TEST).Collection(COLLECTION_PATIENTS).InsertMany(ctx, patients)
	if err != nil {
		fmt.Println("Error while creating patients", err)
		return err
	}
	fmt.Println("Created ", len(cursor.InsertedIDs), " patients")
	return nil
}

func checkPatients(patients []interface{}) error {
	p := Patient{}
	t := reflect.TypeOf(p)
	for _, patient := range patients {
		for i :=0; i <t.NumField(); i++ {
			f := t.Field(i)
			if f.Name == "ID" {
				continue
			}
			tagName := f.Tag.Get("json")
			if tagName == "" {
				tagName = f.Name
			}
			v := patient.(map[string]interface{})
			value, found := v[tagName]
			if !found {
				return errors.New(tagName + " not mentioned for Patient")
			}
			switch vv := value.(type) {
			case string:
				if strings.TrimSpace(vv) == "" {
					return errors.New(tagName + " is empty")
				}
			case float64:
				if vv == 0.0 {
					return errors.New(tagName + " is 0")
				}
			case int:
				if vv == 0 {
					return errors.New(tagName + " is 0")
				}
			}
		}
	}
	return nil
}

func UpdatePatient(ctx context.Context, id string, patient *Patient) (int64, error)  {
	err := checkNilValues(patient)
	if err != nil{
		return 0, err
	}
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, err
	}
	cursor, err := DBClient.client.Database(DB_TEST).Collection(COLLECTION_PATIENTS).UpdateOne(ctx, bson.M{ID: objectId},
	bson.D{
		{"$set", bson.D{
			{"name", patient.Name},
			{"age", patient.Age},
			{"disease", patient.Disease},
		}},
	})
	if err != nil {
		fmt.Println("Error while updating patient", err)
		return 0, err
	}
	fmt.Println("Updated", cursor.ModifiedCount, " patients")
	return cursor.ModifiedCount, err
}

func checkNilValues(p *Patient) error{
	if strings.TrimSpace(p.Name) == "" {
		return errors.New("name is empty")
	}
	if p.Age == 0 {
		return errors.New("age is 0")
	}
	if strings.TrimSpace(p.Disease) == "" {
		return errors.New("disease is empty")
	}
	return nil
}

func DeletePatient(ctx context.Context, id string) (int64, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, err
	}
	cursor, err := DBClient.client.Database(DB_TEST).Collection(COLLECTION_PATIENTS).DeleteOne(ctx,
		bson.D{{ID, objectId}})
	if err != nil {
		fmt.Println("Error while deleting patient", err)
		return 0, err
	}
	fmt.Println("Deleted", cursor.DeletedCount," patients")
	return cursor.DeletedCount, nil
}
