package db

import (
	"database/sql"
	"time"
)

type PropertyMap map[string]interface{}

type Conf struct {
	Id                string          `json:"id"`
	Title             string          `json:"name"`
	Added_by          sql.NullString  `json:"added_by"`
	Start_date        time.Time       `json:"start_date"`
	End_date          time.Time       `json:"end_date"`
	Description       string          `json:"description"`
	Picture           sql.NullString  `json:"picture"`
	Country           string          `json:"country"`
	City              string          `json:"city"`
	Address           string          `json:"address"`
	Category          string          `json:"category"`
	Tickets_available sql.NullBoolean `json:"tickets_available"`
	Discount_program  sql.NullBoolean `json:"discount_program"`
	Min_price         sql.NullInt64   `json:"min_price"`
	Max_price         sql.NullInt64   `json:"max_price"`
	Facebook_account  sql.NullString  `json:"facebook_account"`
	Youtube_account   sql.NullString  `json:"youtube_account"`
	Twitter_account   sql.NullString  `json:"twitter_account"`
	Details           PropertyMap     `json:"details"`
	Speakers          PropertyMap     `json:"speakers"`
	Sponsors          PropertyMap     `json:"sponsors"`
	Verified          bool            `json:"verified"`
	Deleted           bool            `json:"deleted"`
	Created_at        time.Time       `json:"created_at"`
	Updated_at        time.Time       `json:"updated_at"`
}
