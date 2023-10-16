package main

import (
	"github.com/gin-gonic/gin"
)

func (api *ApiServer) handleCreatePost(c *gin.Context) {
	var createPostRequest CreatePostRequest
	err := c.BindJSON(&createPostRequest)
	if err != nil {
		api.ErrorReturn(c, 400, "Invalid request")
		return
	}

	hash, err := api.store.CreatePost(c.GetString("username"), createPostRequest)
	if err != nil {
		api.ErrorReturn(c, 500, "Internal error")
		return
	}

	c.JSON(200, gin.H{
		"message": "Post created",
		"post_id": hash,
	})
}

func (api *ApiServer) handleGetPost(c *gin.Context) {
	postId := c.Param("postId")

	post, err := api.store.GetPost(postId)
	if err != nil {
		api.ErrorReturn(c, 404, "Post not found")
		return
	}

	c.JSON(200, post)
}

func (api *ApiServer) handleGetAuthorPosts(c *gin.Context){
    author := c.Param("author")

    posts, err := api.store.GetAuthorPosts(author)
    if err != nil {
        api.ErrorReturn(c, 404, "Author not found")
        return
    }

    c.JSON(200, posts)
}

func (api *ApiServer) handleGetAllPosts(c *gin.Context){
    posts, err := api.store.debugGetAllPosts()
    if err != nil {
        api.ErrorReturn(c, 404, "No posts found")
        return
    }

    c.JSON(200, posts)
}
