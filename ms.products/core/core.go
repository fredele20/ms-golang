package core

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fredele20/microservice-practice/ms.products/cache"
	"github.com/fredele20/microservice-practice/ms.products/database"
	"github.com/fredele20/microservice-practice/ms.products/models"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductService struct {
	db database.DBInterface
	logger *logrus.Logger
	redis cache.RedisConnection
}

func NewProductService(db database.DBInterface, logger *logrus.Logger, redis cache.RedisConnection) *ProductService {
	return &ProductService{
		db: db,
		logger: logger,
		redis: redis,
	}
}

func (p ProductService) CreateProduct(ctx context.Context, payload models.Product) (*models.Product, error) {

	if err := payload.Validate(); err != nil {
		p.logger.WithError(err).Error("failed to validate product request body before persisting")
		return nil, err
	}

	payload.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	payload.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	payload.ID = primitive.NewObjectID()
	payload.ProductID = payload.ID.Hex()

	product, err := p.db.CreateProduct(ctx, &payload)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return product, nil
}

func (p ProductService) GetProducts(ctx context.Context, filter models.ProductFilter) (*models.ProductList, error) {

	var result models.ProductList
	cacheValue, err := p.redis.Get(ctx, "product_cache")
	if err == redis.Nil {

		result, err := p.db.GetProducts(ctx, filter)
		if err != nil {
			p.logger.WithError(err).Error(err.Error())
			return nil, err
		}
		
		cacheByte, err := json.Marshal(result)
		if err != nil {
			return nil, err
		}
		
		_, err = p.redis.Set(ctx, "product_cache", cacheByte, time.Second * 30)
		if err != nil {
			return nil, err
		}
		
		result.Source = models.DBData
		return result, nil

	} else if err != nil {
		return nil, err
		
	} else {
		err = json.Unmarshal(cacheValue, &result)
		if err != nil {
			return nil, err
		}
		
		result.Source = models.CacheData
		return &result, nil
	}
}
