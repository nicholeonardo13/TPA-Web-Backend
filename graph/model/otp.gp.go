package model

import "time"

type Otp struct {
	ID            int       `json:"id"`
	Code          string    `json:"code"`
	Email         string    `json:"email"`
	CountryRegion int       `json:"country_region"`
	ValidTime     time.Time `json:"valid_time"`
}
