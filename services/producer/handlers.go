package main

import (
	"github.com/troyjoachim/debezium-go-outbox-example/db/sqlc/DAL"

	"github.com/gin-gonic/gin"
)

func createUserHandler(c *gin.Context) {
	ctx := c.Request.Context()

	var user struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Get the connection from the context
	conn, err := GetDBConnFromContext(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	dal := DAL.New(conn)

	newUser, err := dal.CreateUser(ctx, DAL.CreateUserParams{
		Username: user.Username,
		Email:    user.Email,
	})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(200, gin.H{"message": "User updated successfully", "user": newUser})
}
