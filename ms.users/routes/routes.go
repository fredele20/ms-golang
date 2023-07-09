package routes

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/fredele20/microservice-practice/ms.users/core"
	"github.com/fredele20/microservice-practice/ms.users/models"
	"github.com/gin-gonic/gin"
	"github.com/nyaruka/phonenumbers"
)

type UserRoutes struct {
	core *core.UserService
}

func NewUserRoute(core *core.UserService) *UserRoutes {
	return &UserRoutes{
		core: core,
	}
}

// var userCollection *mongo.Collection = mongod.UserCollection()
// var validate = validator.New()

// func checkExistingUser(ctx context.Context, field, value string) (int64, error) {
// 	count, err := userCollection.CountDocuments(ctx, bson.M{field: value})
// 	if err != nil {
// 		log.Panic(err)
// 		fmt.Printf("Error checking for %v", field)
// 		return count, err
// 	}

// 	return count, nil
// }

func parsePhone(phone, iso2 string) (string, error) {
	num, err := phonenumbers.Parse(phone, iso2)
	if err != nil {
		return "", err
	}

	switch phonenumbers.GetNumberType(num) {
	case phonenumbers.VOIP, phonenumbers.VOICEMAIL:
		return "", errors.New("Sorry, this number can not be used")
	}

	return phonenumbers.Format(num, phonenumbers.E164), nil
}

// // var con database.Store

func (u *UserRoutes) Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		newUser, err := u.core.CreateUser(ctx, user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		defer cancel()
		c.JSON(http.StatusOK, newUser)
	}
}

func (u *UserRoutes) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		foundUser, err := u.core.Login(ctx, user.Email, user.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, foundUser)
	}
}

func (u *UserRoutes) Logout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var context, cancel = context.WithTimeout(context.Background(), time.Second * 30)
		fmt.Println(context)
		defer cancel()
		// Get the token from the request header
		// token := ctx.GetHeader("token")
		err := u.core.Logout("token")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, "token")
	}
}

func (u *UserRoutes) ForgotPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx context.Context
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		foundUser, err := u.core.ForgotPassword(ctx, user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, foundUser)
	}
}

func (u UserRoutes) ResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), time.Second * 30)
		defer cancel()

		var user models.ConfirmPasswordRequest
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// token := c.GetHeader("token")

		reset, err := u.core.ResetPassword(ctx, user.Password, user.ConfirmPassword)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, reset)

	}
}

func (u UserRoutes) ListUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		var filter models.ListUserFilter
		if err := c.BindJSON(&filter); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		users, err := u.core.ListUsers(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

// func GetUsers() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		if err := helpers.CheckUserType(c, "ADMIN"); err != nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
// 			return
// 		}
// 		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

// 		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
// 		if err != nil || recordPerPage < 1 {
// 			recordPerPage = 10
// 		}
// 		page, err1 := strconv.Atoi(c.Query("page"))
// 		if err1 != nil || page < 1 {
// 			page = 1
// 		}

// 		startIndex := (page - 1) * recordPerPage
// 		startIndex, err = strconv.Atoi(c.Query("startIndex"))

// 		matchStage := bson.D{{"$match", bson.D{{}}}}
// 		groupStage := bson.D{{"$group", bson.D{
// 			{"_id", bson.D{{"_id", "null"}}},
// 			{"total_count", bson.D{{"$sum", 1}}},
// 			{"data", bson.D{{"$push", "$$ROOT"}}},
// 		}}}
// 		projectStage := bson.D{
// 			{"$project", bson.D{
// 				{"_id", 0},
// 				{"total_count", 1},
// 				{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
// 			}},
// 		}

// 		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
// 			matchStage, groupStage, projectStage,
// 		})

// 		defer cancel()
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing users"})
// 		}

// 		var allUsers []bson.M
// 		if err = result.All(ctx, &allUsers); err != nil {
// 			log.Fatal(err)
// 		}
// 		c.JSON(http.StatusOK, allUsers[0])

// 	}
// }

// only the admin has the acces to this request
// TODO: coming back here to implement this function better
// to allow users to get their informations but allows admin only to access other users info.
// func GetUser() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		userId := c.Param("user_id")
// 		userEmail := c.GetString("email")
// 		userType := c.GetString("user_type")
// 		fmt.Println(userEmail)
// 		fmt.Println(userId)

// 		if userType != "ADMIN" {
// 			c.JSON(http.StatusBadGateway, gin.H{"error": "you can not do this"})
// 			return
// 		}

// 		if err := helpers.MatchUserTypeToUid(c, userId); err != nil {
// 			fmt.Println(userId)
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			fmt.Println(err.Error())
// 			return
// 		}
// 		fmt.Println(userId)
// 		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

// 		var user models.User
// 		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
// 		defer cancel()
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.JSON(http.StatusOK, user)
// 	}
// }

// func GetUserById() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		userId := c.Param("user_id")
// 		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

// 		var user models.User
// 		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
// 		defer cancel()
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.JSON(http.StatusOK, user)
// 	}
// }
