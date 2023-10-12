package main

import (
	"github.com/gin-gonic/gin"
)

func (api *ApiServer) handleUserSignup(c *gin.Context) {
	var signupRequest CreateUserRequest
	err := c.BindJSON(&signupRequest)
	if err != nil {
		api.ErrorReturn(c, 400, "Invalid request")
		return
	}

	err = api.store.CreateUser(signupRequest)
	if err != nil {
		api.ErrorReturn(c, 400, "Username already exists")
		return
	}

	//create the session
	sessionId, err := api.store.CreateSession(signupRequest.Username)
	if err != nil {
		api.ErrorReturn(c, 500, "Internal error")
		return
	}

	c.SetCookie("session", sessionId, 3600, "/", "localhost", false, true)

	c.JSON(200, gin.H{
		"message": "User created",
		"session": sessionId,
	})
}

func (api *ApiServer) handleUserLogin(c *gin.Context) {
	var loginRequest LoginRequest
	err := c.BindJSON(&loginRequest)
	if err != nil {
		api.ErrorReturn(c, 400, "Invalid request")
		return
	}

	user, err := api.store.GetLoginInfo(loginRequest.Username)
	if err != nil {
		api.ErrorReturn(c, 400, "Invalid username/password")
		return
	}

	if !MatchPasswords(loginRequest.Password, user.Password) {
		api.ErrorReturn(c, 400, "Invalid username/password")
		return
	}

	//check if session exists
	session, err := api.store.SessionExists(loginRequest.Username)
	if err != nil && err.Error() != "sql: no rows in result set" {
		api.ErrorReturn(c, 500, "Internal error")
		return
	}

	if session.Valid {
		c.JSON(200, gin.H{
			"message": "User logged in",
			"session": session.SessionId,
		})
		return
	}

	//create the session
	sessionId, err := api.store.CreateSession(loginRequest.Username)
	if err != nil {
		api.ErrorReturn(c, 500, "Internal error")
		return
	}

	c.SetCookie("session", sessionId, 3600, "/", "localhost", false, true)

	c.JSON(200, gin.H{
		"message": len("User logged in"),
		"session": sessionId,
	})
}

func (api *ApiServer) handleAuth(c *gin.Context) {
	session := c.Request.Header.Get("session")
	if session == "" {
		api.ErrorReturn(c, 400, "No session provided")
		return
	}

	_, err := api.store.GetSession(session)
	if err != nil {
		api.ErrorReturn(c, 400, "Invalid session")
		return
	}

	c.JSON(200, gin.H{
		"message": "User authenticated",
	})
}

func (api *ApiServer) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// get the cookie
		session := c.GetHeader("session")
		if session == "" {
			api.ErrorReturn(c, 400, "No session provided")
			return
		}

		if !api.store.ValidSession(session) {
			api.ErrorReturn(c, 400, "Invalid session")
			return
		}

		// check if session exists
		decoded, err := api.store.GetSession(session)
		if err != nil {
			api.ErrorReturn(c, 400, "Invalid session")
			return
		}
		if decoded == "" {
			api.ErrorReturn(c, 400, "Invalid session")
			return
		}

		c.Set("username", decoded)

		c.Next()
	}
}
