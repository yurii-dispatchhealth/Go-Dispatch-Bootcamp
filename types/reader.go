package types

type User struct {
	Id         int    `json:"id"`
	Username   string `json:"username"`
	Identifier string `json:"identifier"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
}
