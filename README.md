# mongo

high-level wrapper for [go.mongodb.org/mongo-driver/mongo](https://www.mongodb.com/docs/drivers/go/current/)

## quick start

0. create any file with a `.env` extension:

```
MONGOURI=<YOUR_MONGODB_CONNECTION_URI>
```

1. `go get github.com/samjtro/mongo`

```go
type Customer struct {
    Name string
    Address string
}

c := mongo.NewClient()

// Find
var customer Customer
bytes, err = c.FindOne("databaseName", "collectionName", bson.D{{"filterKey", "filterValue"}})
err = sonic.Unmarshal(bytes, &customer)

// Find All
var customers []Customer
bytes, err = c.FindAll("databaseName", "collectionName")
err = sonic.Unmarshal(bytes, &customers)

// Insert
customer := Customer{}
err = c.InsertOne("databaseName", "collectionName", customer)

// Update
err = c.UpdateOne("databaseName", "collectionName", bson.D{{"filterKey", "filterValue"}}, bson.D{{"updateKey", "updateValue"}})

// Replace
var customer Customer{}
err = c.ReplaceOne("databaseName", "collectionName", bson.D{{"filterKey", "filterValue"}}, customer)

// Delete
err = c.DeleteOne("databaseName", "collectionName", bson.D{{"filterKey", "filterValue"}})
```
