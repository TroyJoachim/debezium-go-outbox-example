package main

import (
	"encoding/json"

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

	type UserCreatedEvent struct {
		Username string `json:"username"`
	}

	event := UserCreatedEvent{
		Username: newUser.Username,
	}
	eventBytes, err := json.Marshal(event)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to marshal event"})
		return
	}

	// Publish the user creation event
	_, err = dal.CreateOutbox(ctx, DAL.CreateOutboxParams{
		AggregateType: "user",
		AggregateID:   newUser.ID.String(),
		Type:    "user_created",
		Payload: json.RawMessage(eventBytes),
	})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create outbox entry"})
		return
	}

	c.JSON(200, gin.H{"message": "User updated successfully", "user": newUser})
}
