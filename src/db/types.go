package db

import (
	"time"

	"github.com/guregu/null"
)

type PropertyMap map[string]interface{}

type Conf struct {
	Id                string      `json:"id"`
	Title             string      `json:"name"`
	Added_by          string      `json:"added_by"`
	Start_date        time.Time   `json:"start_date"`
	End_date          time.Time   `json:"end_date"`
	Description       string      `json:"description"`
	Picture           null.String `json:"picture"`
	Country           string      `json:"country"`
	City              string      `json:"city"`
	Address           string      `json:"address"`
	Category          string      `json:"category"`
	Tickets_available null.Bool   `json:"tickets_available"`
	Discount_program  null.Bool   `json:"discount_program"`
	Min_price         null.Int    `json:"min_price"`
	Max_price         null.Int    `json:"max_price"`
	Facebook_account  null.String `json:"facebook_account"`
	Youtube_account   null.String `json:"youtube_account"`
	Twitter_account   null.String `json:"twitter_account"`
	Details           PropertyMap `json:"details"`
	Speakers          PropertyMap `json:"speakers"`
	Sponsors          PropertyMap `json:"sponsors"`
	Verified          bool        `json:"verified"`
	Deleted           bool        `json:"deleted"`
	Created_at        time.Time   `json:"created_at"`
	Updated_at        time.Time   `json:"updated_at"`
}

// func (c *Conf) MarshalJSON() ([]byte, error) {
// 	middleNameValue, err := u.MiddleName.Value()
//
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	var middleNameJsonString string
//
// 	if middleNameValue == nil {
// 		middleNameJsonString = "null"
// 	} else {
// 		middleNameJsonString = fmt.Sprintf("\"%s\"", titleValue)
// 	}
//
// 	jsonString := fmt.Sprintf("{\"id\":%d,\"first_name\":%s,\"middle_name\":%s,\"last_name\":%s}", u.Id, u.FirstName, middleNameJsonString, u.LastName)
//
// 	return []byte(jsonString), nil
// }
