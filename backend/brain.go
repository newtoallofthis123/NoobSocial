package main

type Store interface {
	preStart()
	CreateUser(req CreateUserRequest) error
	GetUser(username string) (User, error)
	GetLoginInfo(username string) (PrivateUser, error)
	CreateSession(username string) (string, error)
	GetSession(sessionId string) (string, error)
	SessionExists(username string) (Session, error)
	ValidSession(sessionId string) bool
	CreatePost(username string, req CreatePostRequest) (string, error)
	GetPost(postId string) (Post, error)
    GetAuthorPosts(username string) ([]Post, error)
    debugGetAllPosts() ([]Post, error)
}

func (pq *DbInstance) CreateUser(req CreateUserRequest) error {

	hashPassword, err := HashPassword(req.Password)
	if err != nil {
		return err
	}
	query := `
		INSERT INTO users (username, full_name, email, bio, password)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = pq.db.Exec(query, req.Username, req.FullName, req.Email, req.Bio, hashPassword)
	return err
}

func (pq *DbInstance) GetUser(username string) (User, error) {
	query := `
		SELECT username, full_name, bio, email, created_at FROM users
		WHERE username = $1
	`
	row := pq.db.QueryRow(query, username)
	var user User
	err := row.Scan(&user.Username, &user.FullName, &user.Bio, &user.Email, &user.CreatedAt)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (pq *DbInstance) GetLoginInfo(username string) (PrivateUser, error) {
	query := `
		SELECT username, password FROM users
		WHERE username = $1
	`
	row := pq.db.QueryRow(query, username)
	var user PrivateUser
	err := row.Scan(&user.Username, &user.Password)
	if err != nil {
		return PrivateUser{}, err
	}
	return user, nil
}

func (pq *DbInstance) CreateSession(username string) (string, error) {
	sessionId := GenerateSessionKey(username)
	query := `
		INSERT INTO sessions (session_id, username)
		VALUES ($1, $2)
	`
	_, err := pq.db.Exec(query, sessionId, username)
	if err != nil {
		return "", err
	}
	return sessionId, nil
}

func (pq *DbInstance) GetSession(sessionId string) (string, error) {
	query := `
		SELECT username FROM sessions
		WHERE session_id = $1
	`
	row := pq.db.QueryRow(query, sessionId)
	var username string
	err := row.Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}

func (pq *DbInstance) ValidSession(sessionId string) bool {
	query := `
		SELECT valid FROM sessions
		WHERE session_id = $1
	`
	row := pq.db.QueryRow(query, sessionId)
	var valid bool
	err := row.Scan(&valid)
	if err != nil {
		return false
	}
	return valid
}

func (pq *DbInstance) SessionExists(username string) (Session, error) {
	query := `
		SELECT session_id, username, created_at, valid FROM sessions
		WHERE username = $1
	`

	row := pq.db.QueryRow(query, username)
	var session Session
	err := row.Scan(&session.SessionId, &session.Username, &session.CreatedAt, &session.Valid)
	if err != nil {
		return Session{}, err
	}

	return session, nil
}

func (pq *DbInstance) CreatePost(username string, req CreatePostRequest) (string, error) {
	hash := RanHash(20)
	query := `
		INSERT INTO posts (post_id, username, content)
		VALUES ($1, $2, $3)
	`
	_, err := pq.db.Exec(query, hash, username, req.Content)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func (pq *DbInstance) GetPost(postId string) (Post, error) {
	query := `
		SELECT post_id, username, content, created_at FROM posts
		WHERE post_id = $1
	`
	row := pq.db.QueryRow(query, postId)
	var post Post
	err := row.Scan(&post.PostId, &post.Username, &post.Content, &post.CreatedAt)
	if err != nil {
		return Post{}, err
	}
	return post, nil
}

func (pq *DbInstance) GetAuthorPosts(username string)([]Post, error){
    query := `
    SELECT post_id, username, content, created_at FROM posts
    WHERE username = $1
    `

    rows,err := pq.db.Query(query, username)
    if err != nil{
        return nil, err
    }
    var posts []Post
    for rows.Next(){
        var post Post
        err := rows.Scan(&post.PostId, &post.Username, &post.Content, &post.CreatedAt)
        if err != nil {
            return nil, err
        }
        posts = append(posts, post)
    }

    return posts, nil
}

func (pq *DbInstance) debugGetAllPosts() ([]Post, error){
    query := `
    SELECT post_id, username, content, created_at FROM posts
    `

    rows,err := pq.db.Query(query)
    if err != nil{
        return nil, err
    }
    var posts []Post

    for rows.Next(){
        var post Post
        err := rows.Scan(&post.PostId, &post.Username, &post.Content, &post.CreatedAt)
        if err != nil {
            return nil, err
        }
        posts = append(posts, post)
    }

    return posts, nil
}
