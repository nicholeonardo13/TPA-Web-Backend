package model

import (
	"github.com/nicholeonardo13/gqlgen-todos/database"
	"time"
)

type User struct {
	ID            int        `json:"id"`
	Username      string     `json:"username"`
	Email         string     `json:"email"`
	Password      string     `json:"password"`
	CountryRegion int        `json:"country_region"`
	Money         int        `json:"money"`
	CreateAt      time.Time  `json:"create_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
}

//GetUserIdByUsername check if a user exists in database by given username
func GetUserIdByUsername(userId float64) (int, error) {
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}

	var user User
	db.First(&user, "ID=?", userId)


	return user.ID, nil
}
