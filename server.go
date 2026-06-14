package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
	"github.com/practice-refactoring/skill"
)

type Level struct {
	Key          string   `json:"key"`
	Name         string   `json:"name"`
	Brief        string   `json:"brief"`
	Descriptions []string `json:"descriptions"`
	Level        int      `json:"level"`
}

type Skill struct {
	Key         string   `json:"key"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Logo        string   `json:"logo"`
	Levels      []Level  `json:"levels"`
	Tags        []string `json:"tags"`
}

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

	r.GET("/skills", func(c *gin.Context) {
		rows, err := db.Query("SELECT key, name, description, logo, levels, tags FROM skill")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		skills := []Skill{}
		for rows.Next() {
			var skill Skill
			var levels []byte
			var tags pq.StringArray
			if err := rows.Scan(&skill.Key, &skill.Name, &skill.Description, &skill.Logo, &levels, &tags); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if err := json.Unmarshal(levels, &skill.Levels); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			skill.Tags = tags
			skills = append(skills, skill)
		}

		c.JSON(http.StatusOK, gin.H{"data": skills})
	})

	r.POST("/skills", func(c *gin.Context) {
		var skill Skill
		if err := c.ShouldBindJSON(&skill); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		stmt, err := db.Prepare("INSERT INTO skill (key, name, description, logo, levels, tags) VALUES ($1, $2, $3, $4, $5, $6)")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		levels, err := json.Marshal(skill.Levels)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		tags := pq.StringArray(skill.Tags)
		_, err = stmt.Exec(skill.Key, skill.Name, skill.Description, skill.Logo, levels, tags)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": skill})
	})
	r.Run("127.0.0.1:8080")
}
