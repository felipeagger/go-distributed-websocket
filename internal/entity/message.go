package entity

type Message struct {
	UserID      string `json:"userId"`
	Origin      string `json:"origin"`
	Data        string `json:"data"`
	ReceivedBy  string `json:"ReceivedBy"`
	ProcessedBy string `json:"processedBy"`
}
