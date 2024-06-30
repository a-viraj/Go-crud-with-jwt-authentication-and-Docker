package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/a-viraj/project/database"
	"github.com/a-viraj/project/helper"
	"github.com/a-viraj/project/models"
	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func HashPassword(pass string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), 14)
	if err != nil {
		log.Fatal(err.Error())
	}
	return string(bytes)
}

func VerifyPassword(pass, hash string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	if err != nil {
		return false, "incorrect password"
	}
	return true, ""

}

func Signup() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		defer cancel()
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		defer cancel()

		password := HashPassword(*user.Password)
		user.Password = &password

		count, err := userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the phone number"})
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone number already exists"})
		}

		user.ID = primitive.NewObjectID()
		user.UserId = user.ID.Hex()
		token, refreshToken := helper.GenerateTokens(*user.Email, *user.Name, *user.UserType, user.UserId)
		user.Token = &token
		user.RefreshToken = &refreshToken

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, resultInsertionNumber)
	}

}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		var foundUser models.User
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		res, err := VerifyPassword(*user.Password, *foundUser.Password)
		if err != "" {
			log.Fatal(err)
		}
		if !res {
			c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect password"})
			return
		}
		if foundUser.Email == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
			return
		}
		token, refreshtoken := helper.GenerateTokens(*user.Email, *user.Name, *user.UserType, user.UserId)
		user.Token = &token
		user.RefreshToken = &refreshtoken
		helper.UpdateAllTokens(token, refreshtoken, foundUser.UserId)
		if err := userCollection.FindOne(ctx, bson.M{"email": foundUser.Email}).Decode(&foundUser); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(200, foundUser)
	}
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}
		startindex := (page - 1) * recordPerPage
		matchStage := bson.D{
			{"$match", bson.D{{}}},
		}
		groupStage := bson.D{
			{"$group", bson.D{
				{"_id", bson.D{{"_id", "null"}}},
				{"totalCount", bson.D{{"$sum", 1}}},
				{"data", bson.D{{"$push", "$$ROOT"}}},
			}},
		}
		projectStage := bson.D{{"$project", bson.D{
			{"_id", 0},
			{"totalCount", 1},
			{Key: "user_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startindex, recordPerPage}}}}}}}

		res, err := userCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		var allUsers []bson.M
		if err = res.All(ctx, &allUsers); err != nil {
			log.Fatal(err)
		}
		c.JSON(200, allUsers[0])
	}

}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		userid := c.Param("userId")

		var user models.User

		if err := userCollection.FindOne(ctx, bson.M{"userid": userid}).Decode(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(200, user)
	}
}
