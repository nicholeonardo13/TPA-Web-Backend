package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/nicholeonardo13/gqlgen-todos/database"
	"github.com/nicholeonardo13/gqlgen-todos/graph/generated"
	"github.com/nicholeonardo13/gqlgen-todos/graph/middleware"
	"github.com/nicholeonardo13/gqlgen-todos/graph/model"
	"golang.org/x/crypto/bcrypt"
)

func (r *mutationResolver) Login(ctx context.Context, username *string, password *string) (string, error) {
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}

	var user model.User
	db.First(&user, "username=?", username)

	hashedPassword := CheckPasswordHash(*password, user.Password)

	fmt.Println(user.Password)
	fmt.Println(hashedPassword)

	if hashedPassword == false {
		return "Wrong", nil
	}

	tokiens, _ := GenerateToken(user.ID)
	cookie := &http.Cookie{
		Name:     "CookieJar",
		Value:    tokiens,
		Expires:  time.Time{}.Add(time.Hour * 24),
		HttpOnly: true,
	}

	write := *middleware.ForWrite(ctx)
	http.SetCookie(write, cookie)
	return tokiens, nil
}

func (r *mutationResolver) Register(ctx context.Context, username *string, password *string, otpCode *string) (string, error) {
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}

	var otp model.Otp
	db.Debug().First(&otp, "code=?", otpCode)

	if otp.ID == 0 {
		return "OTP Salah", nil
	}

	fmt.Println(otp.ValidTime)

	var validTime = time.Since(otp.ValidTime)
	fmt.Println(validTime)

	if time.Since(otp.ValidTime).Minutes() >= 2 {
		db.Delete(&otp, "code=?", otpCode)
		return "OTP Tidak Valid", nil
	}

	db.Delete(&otp, "code=?", otpCode)

	hashedPassword := HashPassword(model.User{}.Password)

	user := model.User{Email: otp.Email, Username: *username, Password: hashedPassword, CountryRegion: otp.CountryRegion, Money: 0, CreateAt: time.Now(), UpdatedAt: time.Now(), DeletedAt: nil}
	db.Create(&user)

	return "Berhasil", nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}
	//defer db.Close()

	var users []*model.User
	db.Find(&users)

	return users, nil
}

func (r *queryResolver) User(ctx context.Context) (*model.User, error) {
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}

	userID := middleware.ForContext(ctx)
	fmt.Println(userID.ID)

	var user model.User
	db.Debug().First(&user, "id=?", userID.ID)

	return &user, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

var (
	SecretKey = []byte("rahasia")
)

func GenerateToken(userId int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	/* Create a map to store our claims */
	claims := token.Claims.(jwt.MapClaims)
	/* Set token claims */
	claims["userId"] = userId
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		log.Fatal("Error in Generating key")
		return "", err
	}
	return tokenString, nil
}
func ParseToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["username"].(string)
		return username, nil
	} else {
		return "", err
	}
}
