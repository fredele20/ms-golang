package db

import (
	"context"

	"github.com/fredele20/microservice-practice/ms.users/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserStore interface {
	GetUserByField(ctx context.Context, field, value string) (*models.User, error)
	GetUserByPhone(ctx context.Context, phone string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserById(ctx context.Context, id string) (*models.User, error)
	ListUsers(ctx context.Context, filters models.ListUserFilter) (*models.UserList, error)
	CreateUser(ctx context.Context, payload *models.User) (*models.User, error)
	UpdateUser(ctx context.Context, payload *models.User) (*models.User, error)
	DeactivateUser(ctx context.Context, id string) (*models.User, error)
	ActivateUser(ctx context.Context, id string) (*models.User, error)
	DeleteUser(ctx context.Context, id string) error
	SessionCollection() *mongo.Collection
}

// func DBInstance() *mongo.Client {
// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	MongoDb := os.Getenv("DATABASE_URL")

// 	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	err = client.Connect(ctx)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("connected to MongoDB...")

// 	return client
// }
