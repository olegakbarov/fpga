package domain

import (
	"time"

	"github.com/guregu/null"
)

type PropertyMap map[string]interface{}

type ConferenceInput struct {
	Title             string      `json:"name"`
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
	Facebook          null.String `json:"facebook"`
	Youtube           null.String `json:"youtube"`
	Twitter           null.String `json:"twitter"`
	Details           PropertyMap `json:"details"`
}

type Conference struct {
	ConferenceInput
	Id       string `json:"id"`
	Added_by string `json:"added_by"`
}

type RawConference struct {
	Conference
	Verified   bool      `json:"verified"`
	Deleted    bool      `json:"deleted"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

func (r RawConference) PublicFields() Conference {
	return r.Conference
}
