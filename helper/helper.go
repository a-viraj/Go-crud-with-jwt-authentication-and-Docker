package helper

import (
	"context"
	"log"
	"time"

	"github.com/a-viraj/project/database"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email    string
	Name     string
	UserId   string
	UserType string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var SecretKey string = "piddi"

func GenerateTokens(email string, name string, usertype string, userid string) (string, string) {
	claims := &SignedDetails{
		Email:    email,
		Name:     name,
		UserId:   userid,
		UserType: usertype,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SecretKey))
	if err != nil {
		log.Panic(err)

	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SecretKey))
	if err != nil {
		log.Panic(err)

	}
	return token, refreshToken
}
func UpdateAllTokens(signedToken string, refreshToken string, userId string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var update primitive.D
	defer cancel()
	update = append(update, primitive.E{Key: "token", Value: signedToken})
	update = append(update, primitive.E{Key: "refreshtoken", Value: refreshToken})
	upsert := true
	filter := bson.M{"userId": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{Key: "$set", Value: update},
		},
		&opt,
	)
	if err != nil {
		log.Panic(err)
	}
}
func ValidateTokens(token string) (claims *SignedDetails, msg string) {
	t, err := jwt.ParseWithClaims(token, &SignedDetails{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		msg = err.Error()
		return
	}
	claims, ok := t.Claims.(*SignedDetails)

	if !ok {
		msg = "couldn't parse claims"
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token is expired"
		return
	}
	return claims, msg
}
