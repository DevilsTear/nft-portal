package model

// APIStatus
type APIStatus struct {
	Status  bool   `json:"status" bson:"status"`
	Code    int    `json:"code" bson:"code"`
	Message string `json:"message" bson:"message"`
}

// Paging
type Paging struct {
	Offset int    `json:"offset" bson:"offset"`
	Limit  int    `json:"limit" bson:"limit"`
	Count  string `json:"count" bson:"count"`
}

// Result
type Result struct {
	Data      interface{} `json:"data" bson:"data"`
	Paging    Paging      `json:"paging" bson:"paging"`
	APIStatus APIStatus   `json:"apistatus" bson:"apistatus"`
}
