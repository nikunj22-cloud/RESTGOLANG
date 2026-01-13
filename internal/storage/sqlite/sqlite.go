package sqlite

import (
	"database/sql"
	"fmt"

	// _ for the behind the seen usecases
	"github.com/nikunj/rest-api/internal/config"
	"github.com/nikunj/rest-api/internal/types"
	_ "modernc.org/sqlite"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite", cfg.Storage)

	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students(
id INTEGER PRIMARY KEY AUTOINCREMENT,
name TEXT,
email TEXT,
age INTEGER

)`)
	if err != nil {
		return nil, err
	}
	return &Sqlite{
		Db: db,
	}, nil
}

// () is known as receiver to attach in struct
func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO students(name , email , age) VALUES (? , ? , ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close() //statement close
	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, nil
	}
	lastid, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	return lastid, nil

}
func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("Select * From Students WHERE id = ? LIMIT 1")
	if err != nil {
		return types.Student{}, err
	}

	defer stmt.Close()

	var student types.Student

	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Age, &student.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student find with id %s", fmt.Sprint(id))
		}
		return types.Student{}, fmt.Errorf("query Error: %w", err)
	}
	return student, nil
}
