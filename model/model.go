package model

// APIStatus
type APIStatus struct {
	Status  bool   `json:"status" bson:"status"`
	Code    int    `json:"code" bson:"code"`
	Message string `json:"message" bson:"message"`
}

// PageObject
type PageObject struct {
	Offset *int `json:"offset" bson:"offset"`
	Limit  *int `json:"limit" bson:"limit"`
	Count  *int `json:"count" bson:"count"`
}

// Result
type Result struct {
	Data       interface{} `json:"data" bson:"data"`
	Status     bool        `json:"status" bson:"status"`
	Code       int         `json:"code" bson:"code"`
	Message    string      `json:"message" bson:"message"`
	PageObject PageObject  `json:"page-object" bson:"page-object"`
}
