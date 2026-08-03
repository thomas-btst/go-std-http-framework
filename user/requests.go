package user

type Request struct {
	Name     string `json:"name" validate:"required,min=2"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}
