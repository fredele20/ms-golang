package mongod

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/fredele20/microservice-practice/ms.products/database"
	"github.com/fredele20/microservice-practice/ms.products/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBStore struct {
	client *mongo.Client
	dbName string
	// collectionName string
}

// func NewDBStore(client *mongo.Client, dbName string) *DBStore {
// 	return &DBStore{
// 		client: client,
// 		dbName: dbName,
// 	}
// }

// var client *mongo.Client = db.DBInstance()

// var dbName = os.Getenv("DATABASE_NAME")

func (d DBStore) productColl() *mongo.Collection {
	return d.client.Database(d.dbName).Collection("products")
}

func (d DBStore) PurchasedCollection() *mongo.Collection {
	return d.client.Database(d.dbName).Collection("purchasedProduct")
}

func (d DBStore) CreateProduct(ctx context.Context, payload *models.Product) (*models.Product, error) {
	if _, err := d.productColl().InsertOne(ctx, payload); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return payload, nil
}

func (d DBStore) GetProducts(ctx context.Context, filters models.ProductFilter) (*models.ProductList, error) {

	opts := options.Find()

	if filters.Limit != 0 {
		opts.SetLimit(filters.Limit)
	}

	opts.SetSort(bson.M{"name": 1})

	filter := bson.M{}

	var products []*models.Product

	cursor, err := d.productColl().Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &products); err != nil {
		fmt.Println(err)
		return nil, err
	}

	count, err := d.productColl().CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &models.ProductList{
		Data:  products,
		Count: count,
	}, nil

}

var ErrDuplicate = errors.New("duplicate record")

func DBInstance(connectionUri, databaseName string) (database.DBInterface, error) {
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	// MongoDb := os.Getenv("DATABASE_URL")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.NewClient(options.Client().ApplyURI(connectionUri))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	fmt.Println("connected to MongoDB...")

	return &DBStore{client: client, dbName: databaseName}, nil
}
