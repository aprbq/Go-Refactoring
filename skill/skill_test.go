package skill

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

var table = `
	CREATE TABLE IF NOT EXISTS skill (
	key TEXT PRIMARY KEY,
	name TEXT NOT NULL DEFAULT '',
	description TEXT NOT NULL DEFAULT '',
	logo TEXT NOT NULL DEFAULT '',
	levels JSONB NOT NULL DEFAULT '[]',
	tags TEXT [] NOT NULL DEFAULT '{}'
);`

func TestSkillHandler(t *testing.T) {
	t.Run("GetSkillByKey", func(t *testing.T) {
		db, _ := sql.Open("sqlite", "file:TestSkillHandler?mode=memory&cache=shared")
		defer db.Close()

		gosql := `INSERT INTO skill (key, name, description, levels, tags)
        VALUES (
            'go', 'Go', 'Go is an open source programming...',
            '[{"name": "Beginner", "level": 1, "descriptions": ["basic knowledge ..."]},{"name": "Intermediate", "level": 2, "descriptions": ["complex programs..."]}]',
            '{go,golang}'
        );`

		db.Exec(table)
		db.Exec(gosql)

		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Params = append(c.Params, gin.Param{Key: "key", Value: "go"})

		s := NewStorage(db)
		h := NewHandler(s)
		h.GetSkillByKey(c)

		want := `{"data":{"key":"go","name":"Go","description":"Go is an open source programming...","logo":"","levels":[{"key":"","name":"Beginner","brief":"","descriptions":["basic knowledge ..."],"level":1},{"key":"","name":"Intermediate","brief":"","descriptions":["complex programs..."],"level":2}],"tags":["go","golang"]}}`
		got := rec.Body.String()

		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("GetSkillByKe: should respond error when skill not found", func(t *testing.T) {
		db, _ := sql.Open("sqlite", "file:TestSkillHandler?mode=memory&cache=shared")
		defer db.Close()

		db.Exec(table)

		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Params = append(c.Params, gin.Param{Key: "key", Value: "go"})

		s := NewStorage(db)
		h := NewHandler(s)
		h.GetSkillByKey(c)

		want := `{"error":"sql: no rows in result set"}`
		got := rec.Body.String()

		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("GetSkillByKey: should response error when levels json unmarshal failed", func(t *testing.T) {
		db, _ := sql.Open("sqlite", "file:TestSkillHandler?mode=memory&cache=shared")
		defer db.Close()

		gosql := `INSERT INTO skill (key, name, description, levels, tags)
        VALUES (
            'go', 'Go', 'Go is an open source programming...',
            'invalid json levels',
            '{go,golang}'
        );`

		db.Exec(table)
		db.Exec(gosql)

		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Params = append(c.Params, gin.Param{Key: "key", Value: "go"})

		s := NewStorage(db)
		h := NewHandler(s)
		h.GetSkillByKey(c)

		want := `{"error":"invalid character 'i' looking for beginning of value"}`
		got := rec.Body.String()

		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("GetSkill", func(t *testing.T) {
		db, _ := sql.Open("sqlite", "file:TestSkillHandler?mode=memory&cache=shared")
		defer db.Close()

		gosql := `INSERT INTO skill (key, name, description, levels, tags)
        VALUES (
            'go', 'Go', 'Go is an open source programming...',
            '[{"name": "Beginner", "level": 1, "descriptions": ["basic knowledge ..."]},{"name": "Intermediate", "level": 2, "descriptions": ["complex programs..."]}]',
            '{go,golang}'
        );`

		db.Exec(table)
		db.Exec(gosql)

		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)

		s := NewStorage(db)
		h := NewHandler(s)
		h.GetSkill(c)

		want := `{"data":[{"key":"go","name":"Go","description":"Go is an open source programming...","logo":"","levels":[{"key":"","name":"Beginner","brief":"","descriptions":["basic knowledge ..."],"level":1},{"key":"","name":"Intermediate","brief":"","descriptions":["complex programs..."],"level":2}],"tags":["go","golang"]}]}`
		got := rec.Body.String()

		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("GetSkill: should respond empty array when no skills found", func(t *testing.T) {
		db, _ := sql.Open("sqlite", "file:TestSkillHandler?mode=memory&cache=shared")
		defer db.Close()

		db.Exec(table)

		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)

		s := NewStorage(db)
		h := NewHandler(s)
		h.GetSkill(c)

		want := `{"data":[]}`
		got := rec.Body.String()

		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("GetSkill: should respond error when levels json unmarshal failed", func(t *testing.T) {
		db, _ := sql.Open("sqlite", "file:TestSkillHandler?mode=memory&cache=shared")
		defer db.Close()

		gosql := `INSERT INTO skill (key, name, description, levels, tags)
        VALUES (
            'go', 'Go', 'Go is an open source programming...',
            'invalid json levels',
            '{go,golang}'
        );`

		db.Exec(table)
		db.Exec(gosql)

		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)

		s := NewStorage(db)
		h := NewHandler(s)
		h.GetSkill(c)

		want := `{"error":"invalid character 'i' looking for beginning of value"}`
		got := rec.Body.String()

		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("CreateSkill", func(t *testing.T) {
		db, _ := sql.Open("sqlite", "file:TestSkillHandler?mode=memory&cache=shared")
		defer db.Close()

		db.Exec(table)

		skill := Skill{
			Key:         "go",
			Name:        "Go",
			Description: "Go is an open source programming...",
			Logo:        "",
			Levels: []Level{
				{Name: "Beginner", Level: 1, Descriptions: []string{"basic knowledge ..."}},
				{Name: "Intermediate", Level: 2, Descriptions: []string{"complex programs..."}},
			},
			Tags: []string{"go", "golang"},
		}

		body, _ := json.Marshal(skill)
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/skills", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		c, _ := gin.CreateTestContext(rec)
		c.Request = req

		s := NewStorage(db)
		h := NewHandler(s)
		h.CreateSkill(c)

		want := `{"data":{"key":"go","name":"Go","description":"Go is an open source programming...","logo":"","levels":[{"key":"","name":"Beginner","brief":"","descriptions":["basic knowledge ..."],"level":1},{"key":"","name":"Intermediate","brief":"","descriptions":["complex programs..."],"level":2}],"tags":["go","golang"]}}`
		got := rec.Body.String()

		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("CreateSkill: should respond error when binding JSON failed", func(t *testing.T) {
		db, _ := sql.Open("sqlite", "file:TestSkillHandler?mode=memory&cache=shared")
		defer db.Close()

		db.Exec(table)

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/skills", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		c, _ := gin.CreateTestContext(rec)
		c.Request = req

		s := NewStorage(db)
		h := NewHandler(s)
		h.CreateSkill(c)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("got status %d, want %d", rec.Code, http.StatusBadRequest)
		}

		got := rec.Body.String()
		if !bytes.Contains([]byte(got), []byte("error")) {
			t.Errorf("expected error in response, got %s", got)
		}
	})

	t.Run("CreateSkill: should respond error when duplicate key", func(t *testing.T) {
		db, _ := sql.Open("sqlite", "file:TestSkillHandler?mode=memory&cache=shared")
		defer db.Close()

		gosql := `INSERT INTO skill (key, name, description, levels, tags)
        VALUES (
            'go', 'Go', 'Go is an open source programming...',
            '[{"name": "Beginner", "level": 1, "descriptions": ["basic knowledge ..."]}]',
            '{go,golang}'
        );`

		db.Exec(table)
		db.Exec(gosql)

		skill := Skill{
			Key:         "go",
			Name:        "Go",
			Description: "Go is an open source programming...",
			Logo:        "",
			Levels:      []Level{{Name: "Beginner", Level: 1, Descriptions: []string{"basic knowledge ..."}}},
			Tags:        []string{"go", "golang"},
		}

		body, _ := json.Marshal(skill)
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/skills", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		c, _ := gin.CreateTestContext(rec)
		c.Request = req

		s := NewStorage(db)
		h := NewHandler(s)
		h.CreateSkill(c)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("got status %d, want %d", rec.Code, http.StatusInternalServerError)
		}

		got := rec.Body.String()
		if !bytes.Contains([]byte(got), []byte("error")) {
			t.Errorf("expected error in response, got %s", got)
		}
	})

	t.Run("GetSkill: should respond error when database query fails", func(t *testing.T) {
		db, _ := sql.Open("sqlite", "file:TestSkillHandler?mode=memory&cache=shared")
		db.Exec(table)
		db.Close()

		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)

		s := NewStorage(db)
		h := NewHandler(s)
		h.GetSkill(c)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("got status %d, want %d", rec.Code, http.StatusInternalServerError)
		}

		got := rec.Body.String()
		if !bytes.Contains([]byte(got), []byte("error")) {
			t.Errorf("expected error in response, got %s", got)
		}
	})

	t.Run("CreateSkill: should respond error when database prepare fails", func(t *testing.T) {
		db, _ := sql.Open("sqlite", "file:TestSkillHandler?mode=memory&cache=shared")
		db.Exec(table)
		db.Close()

		skill := Skill{
			Key:         "go",
			Name:        "Go",
			Description: "Go is an open source programming...",
			Logo:        "",
			Levels:      []Level{{Name: "Beginner", Level: 1, Descriptions: []string{"basic knowledge ..."}}},
			Tags:        []string{"go", "golang"},
		}

		body, _ := json.Marshal(skill)
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/skills", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		c, _ := gin.CreateTestContext(rec)
		c.Request = req

		s := NewStorage(db)
		h := NewHandler(s)
		h.CreateSkill(c)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("got status %d, want %d", rec.Code, http.StatusInternalServerError)
		}

		got := rec.Body.String()
		if !bytes.Contains([]byte(got), []byte("error")) {
			t.Errorf("expected error in response, got %s", got)
		}
	})

	t.Run("GetSkill: should respond error when rows.Scan fails", func(t *testing.T) {
		db, _ := sql.Open("sqlite", "file:TestSkillHandler?mode=memory&cache=shared")
		defer db.Close()

		db.Exec(table)

		// Insert skill with empty levels that will parse correctly
		gosql := `INSERT INTO skill (key, name, description, levels, tags)
        VALUES (
            'go', 'Go', 'Go is an open source programming...',
            '[]',
            '{go,golang}'
        );`
		db.Exec(gosql)

		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)

		// Use custom storage that queries with wrong column to force scan error
		customStorage := &testStorageForScanError{db: db}
		h := &testHandlerForScanError{storage: customStorage}
		h.GetSkill(c)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("got status %d, want %d", rec.Code, http.StatusInternalServerError)
		}

		got := rec.Body.String()
		if !bytes.Contains([]byte(got), []byte("error")) {
			t.Errorf("expected error in response, got %s", got)
		}
	})

	t.Run("CreateSkill: should respond error when json.Marshal fails", func(t *testing.T) {
		db, _ := sql.Open("sqlite", "file:TestSkillHandler?mode=memory&cache=shared")
		defer db.Close()

		db.Exec(table)

		// Create a custom handler that uses a storage that will fail on CreateSkill
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)

		body := []byte(`{"key":"test","name":"Test","description":"Test","logo":"","levels":null,"tags":["test"]}`)
		req, _ := http.NewRequest("POST", "/skills", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		// Use custom storage that fails on marshal
		customStorage := &testStorageForMarshalError{db: db}
		h := &testHandlerForMarshalError{storage: customStorage}
		h.CreateSkill(c)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("got status %d, want %d", rec.Code, http.StatusInternalServerError)
		}

		got := rec.Body.String()
		if !bytes.Contains([]byte(got), []byte("error")) {
			t.Errorf("expected error in response, got %s", got)
		}
	})

	t.Run("GetSkill: should respond error when rows.Scan fails", func(t *testing.T) {
		db, _ := sql.Open("sqlite", "file:TestSkillHandler?mode=memory&cache=shared")
		defer db.Close()

		db.Exec(table)

		// Insert skill with empty levels that will parse correctly
		gosql := `INSERT INTO skill (key, name, description, levels, tags)
        VALUES (
            'go', 'Go', 'Go is an open source programming...',
            '[]',
            '{go,golang}'
        );`
		db.Exec(gosql)

		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)

		// Use custom storage that queries with wrong column to force scan error
		customStorage := &testStorageForScanError{db: db}
		h := &testHandlerForScanError{storage: customStorage}
		h.GetSkill(c)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("got status %d, want %d", rec.Code, http.StatusInternalServerError)
		}

		got := rec.Body.String()
		if !bytes.Contains([]byte(got), []byte("error")) {
			t.Errorf("expected error in response, got %s", got)
		}
	})

	t.Run("CreateSkill: should respond error when json.Marshal fails", func(t *testing.T) {
		db, _ := sql.Open("sqlite", "file:TestSkillHandler?mode=memory&cache=shared")
		defer db.Close()

		db.Exec(table)

		// Create a custom handler that uses a storage that will fail on CreateSkill
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)

		body := []byte(`{"key":"test","name":"Test","description":"Test","logo":"","levels":null,"tags":["test"]}`)
		req, _ := http.NewRequest("POST", "/skills", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		// Use custom storage that fails on marshal
		customStorage := &testStorageForMarshalError{db: db}
		h := &testHandlerForMarshalError{storage: customStorage}
		h.CreateSkill(c)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("got status %d, want %d", rec.Code, http.StatusInternalServerError)
		}

		got := rec.Body.String()
		if !bytes.Contains([]byte(got), []byte("error")) {
			t.Errorf("expected error in response, got %s", got)
		}
	})
}

type testStorageForScanError struct {
	db *sql.DB
}

func (s *testStorageForScanError) FindSkillByKey(key string) (Skill, error) {
	row := s.db.QueryRow("SELECT key, name, description, logo, levels, tags FROM skill WHERE key = $1", key)
	r := record{}
	return r.decode(row)
}

func (s *testStorageForScanError) FindSkills() ([]Skill, error) {
	// Query with incorrect column name to force scan error
	rows, err := s.db.Query("SELECT key, name, description, logo, nonexistent_column, tags FROM skill")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	skills := []Skill{}
	for rows.Next() {
		r := record{}
		if err := rows.Scan(&r.Key, &r.Name, &r.Description, &r.Logo, &r.Levels, &r.Tags); err != nil {
			return nil, err
		}
		skills = append(skills, r.toSkills([]Level{}))
	}
	return skills, rows.Err()
}

func (s *testStorageForScanError) CreateSkill(skill Skill) error {
	return nil
}

type testHandlerForScanError struct {
	storage Storage
}

func (h *testHandlerForScanError) GetSkill(c *gin.Context) {
	skills, err := h.storage.FindSkills()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": skills})
}

type testStorageForMarshalError struct {
	db *sql.DB
}

func (s *testStorageForMarshalError) FindSkillByKey(key string) (Skill, error) {
	return Skill{}, nil
}

func (s *testStorageForMarshalError) FindSkills() ([]Skill, error) {
	return []Skill{}, nil
}

func (s *testStorageForMarshalError) CreateSkill(skill Skill) error {
	// Simulate json.Marshal error by returning an error
	return json.Unmarshal([]byte("invalid json"), skill.Levels)
}

type testHandlerForMarshalError struct {
	storage Storage
}

func (h *testHandlerForMarshalError) CreateSkill(c *gin.Context) {
	var skill Skill
	if err := c.ShouldBindJSON(&skill); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.storage.CreateSkill(skill); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": skill})
}
