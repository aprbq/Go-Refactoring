package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/practice-refactoring/skill"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: No .env file found, using system environment variables")
	}
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	uri := os.Getenv("POSTGRES_URI")
	if uri == "" {
		panic("POSTGRES_URI environment variable not set")
	}
	db, err := sql.Open("postgres", uri)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	s := skill.NewStorage(db)
	h := skill.NewHandler(s)

	r.GET("/skills/:key", h.GetSkillByKey)
	r.GET("/skills", h.GetSkill)
	r.POST("/skills", h.CreateSkill)
	r.Run("127.0.0.1:8080")
}
