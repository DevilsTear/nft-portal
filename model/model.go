package model

// APIStatus
type APIStatus struct {
	Status  bool
	Code    int
	Message string
}

// Paging
type Paging struct {
	Offset int
	Limit  int
	Count  string
}

// Result
type Result struct {
	Data      interface{}
	Paging    Paging
	APIStatus APIStatus
}
