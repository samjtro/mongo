package mongo

import (
	"context"
	"io/fs"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MONGOURI string
)

type Client struct{ *mongo.Client }

func init() {
	MONGOURI = os.Getenv("MONGOURI")
}

// find all env files
func findAllEnvFiles() []string {
	var files []string
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		split := strings.Split(d.Name(), ".")
		if len(split) > 1 {
			if split[1] == "env" {
				files = append(files, d.Name())
			}
		}
		return err
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	return files
}

// Returns a working Client
func NewClient() Client {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(MONGOURI))
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Connection established")
	return Client{client}
}

// Finds an item from the given MongoDB Collection & Database
func (c Client) FindOne(databaseName, collectionName string, filter bson.D) ([]byte, error) {
	coll := c.Database(databaseName).Collection(collectionName)
	resp, err := coll.FindOne(context.TODO(), filter).Raw()
	return resp, err
}

// Get all entries from the given MongoDB Collection & Database
func (c Client) FindAll(databaseName, collectionName string) ([]byte, error) {
	coll := c.Database(databaseName).Collection(collectionName)
	resp, err := coll.FindOne(context.TODO(), bson.D{{}}).Raw()
	return resp, err
}

// Inserts an item into the given MongoDB Collection & Database
func (c Client) InsertOne(databaseName, collectionName string, item interface{}) error {
	coll := c.Database(databaseName).Collection(collectionName)
	_, err := coll.InsertOne(context.TODO(), item)
	return err
}

// Updates the specified item from the given MongoDB Core Collection with the given update
func (c Client) UpdateOne(databaseName, collectionName string, filter, update bson.D) error {
	coll := c.Database(databaseName).Collection(collectionName)
	// TODO: Figure out a way to handle this result for errors
	resp, err := coll.UpdateOne(context.TODO(), filter, bson.D{{"$set", update}})
	if err != nil {
		return err
	} else {
		log.Printf("%s Updated in Database: %s, Collection: %s\n", resp.UpsertedID.(string), databaseName, collectionName)
	}
	return nil
}

// Replaces the specified item from the given MongoDB Core Collection with the given replacement
func (c Client) ReplaceOne(collectionName, databaseName string, filter bson.D, replacement interface{}) error {
	coll := c.Database("test").Collection(collectionName)
	resp, err := coll.ReplaceOne(context.TODO(), filter, replacement)
	if err != nil {
		return err
	} else {
		log.Printf("%s Updated in Database: %s, Collection: %s\n", resp.UpsertedID.(string), databaseName, collectionName)
	}
	return nil
}

// Replaces the specified item from the given MongoDB Core Collection with the given replacement
func (c Client) DeleteOne(databaseName, collectionName string, filter bson.D) error {
	coll := c.Database(databaseName).Collection(collectionName)
	_, err := coll.DeleteOne(context.TODO(), filter)
	return err
}

func mongoUpdateMultiple() {
	//TODO
}

func mongoInsertMultiple() {
	//TODO
}

func mongoDeleteMultiple() {
	//TODO
}

// TODO: Broken []byte conversion of net.IP
func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}
