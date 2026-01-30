package request

type UserRegistration struct {
	Email     string    `json:"email" binding:"required,email"`
	Password  string    `json:"password" binding:"required,min=6,max=100"`
	FirstName string    `json:"first_name" binding:"required"`
	LastName  string    `json:"last_name"`
	Addresses []Address `json:"addresses" binding:"omitempty,dive"`
	Phone     string    `json:"phone" binding:"required,min=10,max=10"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type Address struct {
	Title      string `json:"title" binding:"omitempty"`
	IsDefault  bool   `json:"is_default"`
	Line1      string `json:"line_1" binding:"required"`
	Line2      string `json:"line_2"`
	PostalCode string `json:"postal_code" binding:"required,min=5,max=5"`
}
