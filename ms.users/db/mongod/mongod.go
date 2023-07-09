package mongod

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/fredele20/microservice-practice/ms.users/db"
	"github.com/fredele20/microservice-practice/ms.users/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type dbStore struct {
	client         *mongo.Client
	dbName         string
}

func MongoConnection(connectionUri, databaseName string) (db.UserStore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	client, err := mongo.NewClient(options.Client().ApplyURI(connectionUri))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(ctx)

	fmt.Println("connected to mongodb successfully....")

	return &dbStore{client: client, dbName: databaseName}, nil
}

func (u dbStore) userCollection() *mongo.Collection {
	return u.client.Database(u.dbName).Collection("users")
}

func (u dbStore) SessionCollection() *mongo.Collection {
	return u.client.Database(u.dbName).Collection("session")
}

func (u dbStore) GetUserByField(ctx context.Context, field, value string) (*models.User, error) {
	var user models.User
	if err := u.userCollection().FindOne(ctx, bson.M{field: value}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (u dbStore) GetUserByPhone(ctx context.Context, phone string) (*models.User, error) {
	return u.GetUserByField(ctx, "phone", phone)
}

func (u dbStore) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return u.GetUserByField(ctx, "email", email)
}

func (u dbStore) GetUserById(ctx context.Context, id string) (*models.User, error) {
	return u.GetUserByField(ctx, "userid", id)
}

func (u dbStore) ListUsers(ctx context.Context, filters models.ListUserFilter) (*models.UserList, error) {
	opts := options.Find()
	opts.SetProjection(bson.M{
		"password": false,
		"token":    false,
	})

	if filters.Limit != 0 {
		opts.SetLimit(filters.Limit)
	}

	filter := bson.M{}

	if filters.Status != nil && filters.Status.IsValid() {
		filter["status"] = filters.Status.String()
	}

	var users []*models.User

	cursor, err := u.userCollection().Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &users); err != nil {
		print(err)
		return nil, err
	}

	delete(filter, "id")

	count, err := u.userCollection().CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &models.UserList{
		Count: count,
		Data:  users,
	}, nil

}

func (u dbStore) UpdateUser(ctx context.Context, payload *models.User) (*models.User, error) {
	updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	payload.UpdatedAt = updatedAt
	var user models.User
	if err := u.userCollection().FindOneAndUpdate(ctx, bson.M{"userid": payload.UserId}, bson.M{
		"$set": payload,
	}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u dbStore) DeactivateUser(ctx context.Context, id string) (*models.User, error) {
	return u.UpdateUser(ctx, &models.User{UserId: id, Status: models.StatusDeactivated})
}

func (u dbStore) ActivateUser(ctx context.Context, id string) (*models.User, error) {
	return u.UpdateUser(ctx, &models.User{UserId: id, Status: models.StatusActivated})
}

func (u dbStore) ResetPassword(ctx context.Context, id, password string) (*models.User, error) {
	return u.UpdateUser(ctx, &models.User{UserId: id, Password: password})
}

func (u dbStore) DeleteUser(ctx context.Context, id string) error {
	if _, err := u.userCollection().DeleteOne(ctx, bson.M{"userId": id}); err != nil {
		return err
	}

	return nil
}

func (u dbStore) CreateUser(ctx context.Context, payload *models.User) (*models.User, error) {
	filters := bson.M{
		"$or": []bson.M{
			{
				"email": payload.Email,
			},
			{
				"phone": payload.Phone,
			},
		},
	}

	var user models.User

	if err := u.userCollection().FindOne(ctx, filters).Decode(&user); err == nil {
		return nil, ErrDuplicate
	}

	if _, err := u.userCollection().InsertOne(ctx, payload); err != nil {
		return nil, err
	}

	return payload, nil
}

var ErrDuplicate = errors.New("duplicate record")
