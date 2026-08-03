package user

type Response struct {
	ID   ID     `json:"id"`
	Name string `json:"name"`
}

type IDResponse struct {
	ID ID `json:"id"`
}
