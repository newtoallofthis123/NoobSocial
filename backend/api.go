package main

import (
	"github.com/gin-gonic/gin"
)

type ApiServer struct {
	client string
	store  Store
}

func NewApiServer() *ApiServer {
	env := GetEnv()
	return &ApiServer{
		client: env.Client,
		store:  NewDbInstance(),
	}
}

func (api *ApiServer) ErrorReturn(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"message": message,
	})
	c.Abort()
}

func (api *ApiServer) debugMiddleWare() gin.HandlerFunc {
    return func(c *gin.Context) {
        debug_key := c.GetHeader("debug")
        env := GetEnv()
        if debug_key != env.SessionKey {
            api.ErrorReturn(c, 401, "Unauthorized")
        }
        c.Next()
    }
}

func (api *ApiServer) Start() {
	api.store.preStart()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//Auth routes
	r.POST("/signup", api.handleUserSignup)
	r.POST("/login", api.handleUserLogin)
	r.GET("/auth", api.handleAuth)

	protected := r.Group("/", api.AuthMiddleware())
	protected.POST("/new-post", api.handleCreatePost)

	// Post routes
	r.GET("/post/:postId", api.handleGetPost)
    r.GET("/posts/:author", api.handleGetAuthorPosts)

    debug := r.Group("/debug", api.debugMiddleWare())
    debug.GET("/posts", api.handleGetAllPosts)

	err := r.Run(api.client)
	if err != nil {
		panic(err)
	}
}
