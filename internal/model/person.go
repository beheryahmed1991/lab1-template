package model

type Person struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Age     *int    `json:"age,omitempty"`
	Address *string `json:"address,omitempty"`
	Work    *string `json:"work,omitempty"`
}

// for new person creation request
type PersonRequest struct {
	Name    string `json:"name" binding:"required"`
	Age     *int    `json:"age"`
	Address *string `json:"address"`
	Work    *string `json:"work"`
}
