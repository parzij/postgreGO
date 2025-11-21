package users

type User struct {
	ID uint `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
	Age int `json:"age"`
	Email string `json:"email"`
	Password string `json:"password"`
}

func (User) TableName() string {
	return "users"
}
