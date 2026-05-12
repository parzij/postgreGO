package users

// @Description Полная модель пользователя, хранящаяся в БД
type User struct {
	ID uint `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
	Age int `json:"age"`
	Email string `json:"email"`
	Password string `json:"-"`
}

func (User) TableName() string {
	return "users"
}

// @Description Данные, которые клиент отправляет при регистрации
type CreateUserRequest struct {
	Name string `json:"name" minLength:"1"`
	Age int `json:"age" minimum:"1" maximum:"99"`
	Email string `json:"email" format:"email"`
	Password string `json:"password" minLength:"8"`
}

// @Description Безопасная версия пользователя для отправки клиенту
type UserResponse struct {
	ID uint `json:"id"`
	Name string `json:"name"`
	Age int `json:"age"`
	Email string `json:"email"`
}

// @Description Возвращается при ошибках валидации или сервера
type ErrorResponse struct {
	Error string `json:"error"`
}