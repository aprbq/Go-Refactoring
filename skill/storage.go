package skill

import (
	"database/sql"
	"encoding/json"

	"github.com/lib/pq"
)

type Storage interface {
	FindSkillByKey(key string) (Skill, error)
	FindSkills() ([]Skill, error)
	CreateSkill(skill Skill) error
}

type storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) storage {
	return storage{db: db}
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

func (s storage) FindSkillByKey(key string) (Skill, error) {
	row := s.db.QueryRow("SELECT key, name, description, logo, levels, tags FROM skill WHERE key = $1", key)

	r := record{}
	return r.decode(row)
}

func (s storage) FindSkills() ([]Skill, error) {
	rows, err := s.db.Query("SELECT key, name, description, logo, levels, tags FROM skill")
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
		lvl, err := r.unmarshalLevels()
		if err != nil {
			return nil, err
		}
		skills = append(skills, r.toSkills(lvl))
	}
	return skills, rows.Err()
}

func (s storage) CreateSkill(skill Skill) error {
	levels, err := json.Marshal(skill.Levels)
	if err != nil {
		return err
	}

	tags := pq.StringArray(skill.Tags)
	stmt, err := s.db.Prepare("INSERT INTO skill (key, name, description, logo, levels, tags) VALUES ($1, $2, $3, $4, $5, $6)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(skill.Key, skill.Name, skill.Description, skill.Logo, levels, tags)
	return err
}
