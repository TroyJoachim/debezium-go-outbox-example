package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// AppContext holds shared resources for your application
type AppContext struct {
	DB *pgxpool.Pool
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Create a context for pool connection (e.g., for initial connection timeout)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	connString := os.Getenv("DATABASE_URL")
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Unable to reach database: %v\n", err)
	}
	log.Println("Connected to the database successfully!")

	appCtx := &AppContext{
		DB: pool,
	}

	router := gin.Default()

	// Use the middleware to provide a per-request database connection
	router.Use(DBConnectionMiddleware(appCtx.DB))

	// Routes
	router.POST("/users", createUserHandler)

	router.Run("localhost:5010")
}

// DBConnectionMiddleware acquires a database connection from the pool
// and stores it in the Gin context for the duration of the request.
func DBConnectionMiddleware(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := pool.Acquire(c.Request.Context()) // Use request context for acquire timeout
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "Failed to acquire database connection"})
			return
		}
		defer conn.Release() // Release the connection back to the pool when the request is done

		// Store the connection in the Gin context
		c.Set("dbConn", conn)

		c.Next() // Proceed to the next middleware/handler
	}
}

// GetDBConnFromContext retrieves the *pgxpool.Conn from the Gin context.
func GetDBConnFromContext(c *gin.Context) (*pgxpool.Conn, error) {
	val, exists := c.Get("dbConn")
	if !exists {
		return nil, fmt.Errorf("database connection not found in context")
	}

	conn, ok := val.(*pgxpool.Conn)
	if !ok {
		return nil, fmt.Errorf("invalid type for database connection in context")
	}

	return conn, nil
}
