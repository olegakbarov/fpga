package db

type PropertyMap map[string]interface{}

type Conf struct {
	Id                string      `json:"id"`
	Title             string      `json:"name"`
	Added_by          string      `json:"added_by"`
	Start_date        string      `json:"start_date"`
	End_date          string      `json:"end_date"`
	Description       string      `json:"description"`
	Picture           string      `json:"picture"`
	Country           string      `json:"country"`
	City              string      `json:"city"`
	Adress            string      `json:"adress"`
	Category          string      `json:"category"`
	Min_price         int         `json:"min_price"`
	Max_price         int         `json:"max_price"`
	Facebook_account  string      `json:"facebook_account"`
	Youtube_account   string      `json:"youtube_account"`
	Twitter_account   string      `json:"twitter_account"`
	Tickets_available bool        `json:"tickets_available"`
	Discount_program  bool        `json:"discount_program"`
	Details           PropertyMap `json:"details"`
	Speakers          PropertyMap `json:"speakers"`
	Sponsors          PropertyMap `json:"sponsors"`
	Verified          bool        `json:"verified"`
	Created_at        string      `json:"created_at"`
	Updated_at        string      `json:"updated_at"`
}
