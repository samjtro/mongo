package mongo

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MONGOURI string
)

type Client struct{ *mongo.Client }

func init() {
	err := godotenv.Load(findAllEnvFiles()...)
	if err != nil {
		log.Fatal(err.Error())
	}
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

// Returns a workable *mongo.Client connection to the Turba cluster
func NewClient() Client {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(MONGOURI))

	if err != nil {
		fmt.Printf("[err] %s", err.Error())
	} else {
		fmt.Println("[log] conn established")
	}

	return Client{client}
}

// Finds an item from the given MongoDB Core Collection
func (c Client) FindOne(collectionName string, filter bson.D) []byte {
	coll := c.Database("test").Collection(collectionName)
	b, err := coll.FindOne(context.TODO(), filter).Raw()

	if err != nil {
		fmt.Printf("[ERROR]: %s", err.Error())
	}

	return b
}

// TODO Get all entries from the given MongoDB Core Collection
func (c Client) FindAll(collectionName string) interface{} {
	coll := c.Database("test").Collection(collectionName)
	var result interface{}
	err := coll.FindOne(context.TODO(), bson.D{{}}).Decode(&result)

	if err != nil {
		fmt.Printf("[ERROR]: %s", err.Error())
	}

	return result
}

// Inserts an item into the given MongoDB Core Collection
func (c Client) InsertOne(collectionName string, item interface{}) {
	coll := c.Database("test").Collection(collectionName)
	_, err := coll.InsertOne(context.TODO(), item)

	if err != nil {
		fmt.Printf("[ERROR]: %s", err.Error())
	} else {
		fmt.Printf("[LOG]: Item Inserted into Collection %s.\n", collectionName)
	}
}

// Updates the specified item from the given MongoDB Core Collection with the given update
func (c Client) UpdateOne(collectionName, filterKey, filterValue string, update bson.D) {
	coll := c.Database("test").Collection(collectionName)

	// TODO: Figure out a way to handle this result for errors
	_, err := coll.UpdateOne(context.TODO(), bson.D{{filterKey, filterValue}}, bson.D{{"$set", update}})

	if err != nil {
		fmt.Printf("[ERROR]: %s", err.Error())
	} else {
		fmt.Printf("[LOG]: Item Updated in Collection %s.\n", collectionName)
	}
}

// Replaces the specified item from the given MongoDB Core Collection with the given replacement
func (c Client) MongoReplaceOne(collectionName, filterKey, filterValue string, replacement interface{}) {
	coll := c.Database("test").Collection(collectionName)
	_, err := coll.ReplaceOne(context.TODO(), bson.D{{filterKey, filterValue}}, replacement)

	if err != nil {
		fmt.Printf("[ERROR]: %s", err.Error())
	} else {
		fmt.Printf("[LOG]: Item Replaced in Collection %s.\n", collectionName)
	}
}

func (c Client) MongoDeleteOne(collectionName, filterKey, filterValue string) {
	coll := c.Database("test").Collection(collectionName)
	_, err := coll.DeleteOne(context.TODO(), bson.D{{filterKey, filterValue}})

	if err != nil {
		log.Println(err)
	}
}

func MongoUpdateMultiple() {
	//TODO
}

func MongoInsertMultiple() {
	//TODO
}

func MongoDeleteMultiple() {
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
