package skill

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
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

type handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) handler {
	return handler{db: db}
}

type record struct {
	Key         string
	Name        string
	Description string
	Logo        string
	Levels      []byte
	Tags        pq.StringArray
}

func (r record) toSkills(lvl []Level) Skill {
	return Skill{
		Key:         r.Key,
		Name:        r.Name,
		Description: r.Description,
		Logo:        r.Logo,
		Tags:        r.Tags,
		Levels:      lvl,
	}
}

func (r record) unmarshalLevels() ([]Level, error) {
	lvl := []Level{}
	if err := json.Unmarshal(r.Levels, &lvl); err != nil {
		return []Level{}, err
	}
	return lvl, nil
}

func findSkillByKey(db *sql.DB, key string) (Skill, error) {
	row := db.QueryRow("SELECT key, name, description, logo, levels, tags FROM skill WHERE key = $1", key)

	r := record{}
	if err := row.Scan(&r.Key, &r.Name, &r.Description, &r.Logo, &r.Levels, &r.Tags); err != nil {
		return Skill{}, err
	}

	lvl, err := r.unmarshalLevels()
	if err != nil {
		return Skill{}, err
	}

	return r.toSkills(lvl), nil
}

func (h handler) GetSkillByKey(c *gin.Context) {
	key := c.Param("key")

	skill, err := findSkillByKey(h.db, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": skill})
}
