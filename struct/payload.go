package _struct

type CustomerToBusiness struct {
	Shortcode     string `json:"shortcode"`
	CommandId     string `json:"command_id"`
	Amount        int    `json:"amount"`
	Msisdn        string `json:"msisdn"`
	BillRefNumber string `json:"account"`
}

type Config struct {
	ConsumerKey    string `json:"consumer_key"`
	ConsumerSecret string `json:"consumer_secret"`
	Environment    string `json:"environment"`
}
