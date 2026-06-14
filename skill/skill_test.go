package skill

import (
	"database/sql"
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
}
