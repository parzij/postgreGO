# **postgreGO**

Небольшой учебный backend-сервис на **Go + GORM + PostgreSQL** для тренировки работы с базой данных:

- регистрирует пользователей в БД;
- валидирует почту (ограниченный список доменов);
- проверяет возраст и сложность пароля;
- умеет отдавать список пользователей, получать пользователя по ID и удалять его.

Проект заточен под консольное тестирование через `curl` и просмотр состояния БД через `psql`.

---

## Архитектура проекта

Проект организован по «взрослой» схеме: точка входа в `cmd/`, бизнес-логика и инфраструктура в `internal/`.

```text
.
├─ cmd/
│  └─ api/
│     └─ main.go           # Точка входа: загрузка .env, конфиг, БД, роутер
├─ internal/
│  ├─ config/
│  │  └─ config.go         # Загрузка конфигурации из переменных окружения
│  ├─ database/
│  │  └─ database.go       # Инициализация GORM + PostgreSQL
│  └─ users/
│     ├─ model.go          # Модель User (ID, Name, Age, Email, Password)
│     └─ handlers.go       # HTTP-хендлеры и роуты для работы с пользователями
├─ .env                    # Локальные переменные окружения (в git не коммитим)
├─ go.mod
├─ go.sum
└─ README.md
````

### Что делает каждый слой

* **`cmd/api/main.go`**

  * Загружает `.env` через `godotenv`.
  * Читает конфиг через `config.LoadConfig()`.
  * Создаёт подключение к БД через `database.New()`.
  * Запускает миграцию `AutoMigrate(&users.User{})`.
  * Поднимает HTTP-сервер на `:8080` c роутером `chi`.

* **`internal/config/config.go`**

  * Читает переменные окружения:

    * `HOST_DB`, `PORT_DB`, `USER_DB`, `PASSWORD_DB`, `NAME_DB`, `SSLMODE_DB`.
  * Проверяет, что все параметры заданы, иначе возвращает ошибку.

* **`internal/database/database.go`**

  * Собирает DSN для PostgreSQL:

    ```text
    host=... user=... password=... dbname=... port=... sslmode=...
    ```
  * Открывает соединение через GORM (`gorm.Open(postgres.Open(dsn))`).

* **`internal/users/model.go`**

  * Описывает структуру `User`:

    ```go
    type User struct {
        ID       uint   `json:"id" gorm:"primaryKey"`
        Name     string `json:"name"`
        Age      int    `json:"age"`
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    ```
  * Опционально задаёт имя таблицы `users` через `TableName()`.

* **`internal/users/handlers.go`**

  * `Handler` хранит `*gorm.DB`.
  * `RegisterRoutes` вешает маршруты:

    * `POST   /users`      — создать пользователя;
    * `GET    /users`      — список пользователей;
    * `GET    /users/{id}` — получить пользователя по ID;
    * `DELETE /users/{id}` — удалить пользователя.
  * В `CreateUser`:

    * нормализует email (`strings.ToLower`, `TrimSpace`);
    * проверяет формат email и домен (только `gmail.com`, `yandex.ru`, `mail.ru`, `yahoo.com`);
    * проверяет возраст и длину пароля;
    * сохраняет пользователя через `DB.Create(&user)`.

---



## Идеи для развития

* Обновление пользователя (`PUT /users/{id}`).
* Хеширование пароля вместо хранения в открытом виде.
* Пагинация и фильтрация списка пользователей.
* JWT-авторизация и простая регистрация/логин.
* Docker-compose для поднятия Go + PostgreSQL в контейнерах.


