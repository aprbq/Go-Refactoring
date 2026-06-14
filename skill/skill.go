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
	storage storage
}

func NewHandler(db *sql.DB) handler {
	return handler{storage: storage{db: db}}
}

type storage struct {
	db *sql.DB
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
	err := json.Unmarshal(r.Levels, &lvl)
	return lvl, err
}

func (r record) decode(row *sql.Row) (Skill, error) {
	if err := row.Scan(&r.Key, &r.Name, &r.Description, &r.Logo, &r.Levels, &r.Tags); err != nil {
		return Skill{}, err
	}

	lvl, err := r.unmarshalLevels()
	return r.toSkills(lvl), err
}

func (s storage) findSkillByKey(key string) (Skill, error) {
	row := s.db.QueryRow("SELECT key, name, description, logo, levels, tags FROM skill WHERE key = $1", key)

	r := record{}
	return r.decode(row)
}

func (h handler) GetSkillByKey(c *gin.Context) {
	key := c.Param("key")

	skill, err := h.storage.findSkillByKey(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": skill})
}
