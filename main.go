package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/pallat/todoapi/router"
	"github.com/pallat/todoapi/store"
	"github.com/pallat/todoapi/todo"
)

// var (
// 	buildcommit = "dev"
// 	buildtime   = time.Now().String()
// )

func main() {
	// _, err := os.Create("/tmp/live")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer os.Remove("tmp/live")

	err := godotenv.Load("local.env")
	if err != nil {
		log.Printf("please consider environment variables: %s\n", err)
	}

	db, err := gorm.Open(mysql.Open(os.Getenv("DB_CONN")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if err := db.AutoMigrate(&todo.Todo{}); err != nil {
		log.Println("auto migrate db", err)
	}

	// move router
	// r := gin.Default()
	// config := cors.DefaultConfig()
	// config.AllowOrigins = []string{
	// 	"http://localhost:8080",
	// }
	// config.AllowHeaders = []string{
	// 	"Origin",
	// 	"Authorization",
	// 	"TransactionID",
	// }
	// r.Use(cors.New(config))
	r := router.NewMyRouter()

	// r.GET("/healthz", func(c *gin.Context) {
	// 	c.Status(200)
	// })
	// r.GET("/limitz", limitedHandler)
	// r.GET("/x", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"buildcommit": buildcommit,
	// 		"buildtime":   buildtime,
	// 	})
	// })

	// gormStore := todo.NewGormStore(db)
	gormStore := store.NewGormStore(db)

	handler := todo.NewTodoHandler(gormStore)
	// r.POST("/todos", handler.NewTask)
	r.POST("/todos", handler.NewTask)
	// r.GET("/todos", handler.List)
	// r.DELETE("/todos/:id", handler.Remove)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	s := &http.Server{
		Addr:           ":" + os.Getenv("PORT"),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	stop()
	fmt.Println("shutting down gracefully, press Ctrl+C again to force")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(timeoutCtx); err != nil {
		fmt.Println(err)
	}
}

// var limiter = rate.NewLimiter(5, 5)

// func limitedHandler(c *gin.Context) {
// 	if !limiter.Allow() {
// 		c.AbortWithStatus(http.StatusTooManyRequests)
// 		return
// 	}
// 	c.JSON(200, gin.H{
// 		"message": "pong",
// 	})
// }
