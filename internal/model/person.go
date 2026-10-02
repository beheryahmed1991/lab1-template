package model

type Person struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
	Work    string `json:"work"`
}

// for new person creation request
type PersonRequest struct {
	Name    string `json:"name" binding:"required"`
	Age     int    `json:"age"`
	Address string `json:"address"`
	Work    string `json:"work"`
}
