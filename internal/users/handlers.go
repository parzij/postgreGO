package users

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func Newhandler(db *gorm.DB) *Handler {
	return &Handler{DB: db}
}

// @Description Подключает все CRUD-хендлеры к роутеру chi
func (hand *Handler) RegisterRoutes(reg chi.Router) {
	reg.Post("/users", hand.CreateUser)
	reg.Get("/users", hand.ListUser)
	reg.Get("/users/{id}", hand.GetUser)
	reg.Delete("/users/{id}", hand.DeleteUser)
}

var allowedDomains = map[string]bool{
	"gmail.com": true,
	"yandex.ru": true,
	"mail.ru": true,
	"yahoo.com": true,
}

// CreateUser создаёт нового пользователя
// @Summary      Создать пользователя
// @Description  Регистрация нового пользователя с валидацией: возраст 1-99, пароль ≥8 символов, email с разрешённым доменом
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      CreateUserRequest  true  "Данные для регистрации"
// @Success      201   {object}  UserResponse       "Пользователь успешно создан"
// @Failure      400   {object}  ErrorResponse      "Неверный запрос: валидация или формат JSON"
// @Failure      500   {object}  ErrorResponse      "Ошибка базы данных"
// @Router       /users [post]
func (hand *Handler) CreateUser(write http.ResponseWriter, req *http.Request) {
	var input struct {
		Name string `json:"name"`
		Age int    `json:"age"`
		Email string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		http.Error(write, "invalid JSON", http.StatusBadRequest)
		return
	}

	if input.Age <= 0 || input.Age >= 100 {
		http.Error(write, "age must be real", http.StatusBadRequest)
		return
	}

	if len(input.Password) < 8 {
		http.Error(write, "password must be more stronge", http.StatusBadRequest)
		return
	}

	email := strings.ToLower(strings.TrimSpace(input.Email))

	if !strings.Contains(email, "@") {
		http.Error(write, "incorrect email", http.StatusBadRequest)
		return
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 || parts[1] == "" {
		http.Error(write, "incorrect email", http.StatusBadRequest)
		return
	}

	domain := parts[1]

	if !allowedDomains[domain] {
		http.Error(write, "this domain don`t support now", http.StatusBadRequest)
		return
	}

	user := User{
		Name: input.Name,
		Age: input.Age,
		Email: email,
		Password: input.Password,
	}

	if err := hand.DB.Create(&user).Error; err != nil {
		http.Error(write, "data base error", http.StatusBadRequest)
		return
	}

	writeJSON(write, http.StatusCreated, UserResponse{
		ID: user.ID,
		Name: user.Name,
		Age: user.Age,
		Email: user.Email,
	})
}

// ListUser возвращает список всех пользователей
// @Summary      Получить список пользователей
// @Description  Возвращает массив всех зарегистрированных пользователей
// @Tags         users
// @Produce      json
// @Success      200  {array}   UserResponse  "Список пользователей (пустой массив, если нет данных)"
// @Failure      500  {object}  ErrorResponse  "Ошибка базы данных"
// @Router       /users [get]
func (hand *Handler) ListUser(write http.ResponseWriter, _ *http.Request) {
	var users []User

	if err := hand.DB.Find(&users).Error; err != nil {
		http.Error(write, "data base error", http.StatusBadRequest)
		return
	}

	if users == nil {
		users = []User{}
	}

	response := make([]UserResponse, len(users))
	for i, u := range users {
		response[i] = UserResponse{
			ID: u.ID,
			Name: u.Name,
			Age: u.Age,
			Email: u.Email,
		}
	}

	writeJSON(write, http.StatusOK, response)
}

// GetUser возвращает пользователя по ID
// @Summary      Получить пользователя по ID
// @Description  Возвращает данные пользователя, если он существует в БД
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "ID пользователя"  minimum(1)
// @Success      200  {object}  UserResponse
// @Failure      400  {object}  ErrorResponse  "Неверный формат ID"
// @Failure      404  {object}  ErrorResponse  "Пользователь не найден"
// @Failure      500  {object}  ErrorResponse  "Ошибка базы данных"
// @Router       /users/{id} [get]
func (hand *Handler) GetUser(write http.ResponseWriter, req *http.Request) {
	idStr := chi.URLParam(req, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(write, "invalid ID", http.StatusBadRequest)
		return
	}

	var user User
	err = hand.DB.First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(write, "user is not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(write, "data base error", http.StatusInternalServerError)
		return
	}

	writeJSON(write, http.StatusOK, UserResponse{
		ID: user.ID,
		Name: user.Name,
		Age: user.Age,
		Email: user.Email,
	})
}

// DeleteUser удаляет пользователя по ID
// @Summary      Удалить пользователя
// @Description  Удаляет пользователя из базы данных по указанному ID
// @Tags         users
// @Produce      json
// @Param        id   path  int  true  "ID пользователя"  minimum(1)
// @Success      204  "Пользователь успешно удалён"
// @Failure      400  {object}  ErrorResponse  "Неверный формат ID"
// @Failure      404  {object}  ErrorResponse  "Пользователь не найден"
// @Failure      500  {object}  ErrorResponse  "Ошибка базы данных"
// @Router       /users/{id} [delete]
func (hand *Handler) DeleteUser(write http.ResponseWriter, req *http.Request) {
	idStr := chi.URLParam(req, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(write, "invalid id", http.StatusBadRequest)
		return
	}

	res := hand.DB.Delete(&User{}, id)
	if res.Error != nil {
		http.Error(write, "data base error", http.StatusInternalServerError)
		return
	}

	if res.RowsAffected == 0 {
		http.Error(write, "user id not found", http.StatusNotFound)
		return
	}

	write.WriteHeader(http.StatusNoContent)
}

func writeJSON(write http.ResponseWriter, status int, v any) {
	write.Header().Set("Content-Type", "application/json; charset=utf-8")
	write.WriteHeader(status)
	_ = json.NewEncoder(write).Encode(v)
}