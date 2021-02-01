package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	mailjet "github.com/mailjet/mailjet-apiv3-go"
	"github.com/nicholeonardo13/gqlgen-todos/database"
	"github.com/nicholeonardo13/gqlgen-todos/graph/model"
)

func (r *mutationResolver) CreateOtp(ctx context.Context, email *string, countryRegion *int) (string, error) {
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}

	otp := model.Otp{Code: StringRandom(5), Email: *email, ValidTime: time.Now(), CountryRegion: *countryRegion}
	db.Create(&otp)

	SendOTP(*email, otp.Code)

	return otp.Code, nil
}

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
const charset = "abcdeABCDE12345"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
func StringRandom(length int) string {
	return StringWithCharset(length, charset)
}
func SendOTP(email, code string) {
	mailjetClient := mailjet.NewMailjetClient("6c57474f8caa1ec637a5de594e1d9620", "0e3ba056e9806c9dbbb2e9f5116e10e7")
	messagesInfo := []mailjet.InfoMessagesV31{
		mailjet.InfoMessagesV31{
			From: &mailjet.RecipientV31{
				Email: "nicholeonardo13@gmail.com",
				Name:  "Staem Admin",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: email,
					Name:  email,
				},
			},
			Subject:  "Welcome to Staem.",
			TextPart: "Welcome New User",
			HTMLPart: "<h3>Dear New User , welcome to <a href='https://store.steampowered.com'>Staem</a>!</h3><br/> <h1>This is your OTP : " + code + "</h1>",
			CustomID: "AppGettingStartedTest",
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}
	res, err := mailjetClient.SendMailV31(&messages)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Data: %+v\n", res)
}
