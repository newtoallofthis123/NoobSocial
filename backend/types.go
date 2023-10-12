package main

type Env struct {
	Db         string
	User       string
	Password   string
	Host       string
	Port       string
	Client     string
	SessionKey string
}

type CreateUserRequest struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Bio      string `json:"bio"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	Username  string `json:"username"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	Bio       string `json:"bio"`
	CreatedAt string `json:"created_at"`
}

type PrivateUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Session struct {
	SessionId string `json:"session_id"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
	Valid     bool   `json:"valid"`
}

type CreatePostRequest struct {
	Content string `json:"content"`
}

type Post struct {
	PostId    string `json:"post_id"`
	Username  string `json:"username"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}
