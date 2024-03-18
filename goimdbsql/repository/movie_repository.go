package repository

import (
	"database/sql"
	"log"

	"github.com/Pugpaprika21/goimdb/database"
	"github.com/Pugpaprika21/goimdb/dto"
)

type (
	IMovieRepository interface {
		CreateTable()
		GetAll() ([]*dto.Movie, error)
		GetByYear(y int) ([]*dto.Movie, error)
		GetByID(imdbID string) (*dto.Movie, error)
		Create(movie *dto.Movie) (*dto.Movie, error)
		Update(id string, movie *dto.Movie) (*dto.Movie, error)
		Delete(id string) (bool, error)
	}

	MovieRepository struct {
		db *sql.DB
	}
)

func NewMovieRepository() *MovieRepository {
	return &MovieRepository{
		db: database.New().GetDB(),
	}
}

func (m *MovieRepository) CreateTable() {
	createTb := `
	CREATE TABLE IF NOT EXISTS goimdb (
		id INT AUTO_INCREMENT,
		imdbID TEXT NOT NULL UNIQUE,
		title TEXT NOT NULL,
		year INT NOT NULL,
		rating FLOAT NOT NULL,
		isSuperHero BOOLEAN NOT NULL,
		PRIMARY KEY (id)
	);`

	if _, err := m.db.Exec(createTb); err != nil {
		log.Fatal("create table error ", err)
	}
}

func (m *MovieRepository) GetAll() ([]*dto.Movie, error) {
	mvs := []*dto.Movie{}

	rows, err := m.db.Query(`SELECT id, imdbID, title, year, rating, isSuperHero FROM goimdb`)
	if err != nil {
		log.Fatal("query error", err)
	}
	defer rows.Close()

	for rows.Next() {
		var m dto.Movie
		if err := rows.Scan(&m.ID, &m.ImdbID, &m.Title, &m.Year, &m.Rating, &m.IsSuperHero); err != nil {
			log.Fatal("scan error", err)
		}
		mvs = append(mvs, &m)
	}

	if err := rows.Err(); err != nil {
		log.Fatal("rows error", err)
	}

	return mvs, err
}

func (m *MovieRepository) GetByYear(year int) ([]*dto.Movie, error) {
	var mvs []*dto.Movie

	rows, err := m.db.Query(`SELECT id, imdbID, title, year, rating, isSuperHero FROM goimdb WHERE year = ?`, year)
	if err != nil {
		return mvs, err
	}
	defer rows.Close()

	for rows.Next() {
		var m dto.Movie
		if err := rows.Scan(&m.ID, &m.ImdbID, &m.Title, &m.Year, &m.Rating, &m.IsSuperHero); err != nil {
			return mvs, err
		}
		mvs = append(mvs, &m)
	}

	if err := rows.Err(); err != nil {
		return mvs, err
	}

	return mvs, nil
}

func (m *MovieRepository) GetByID(imdbID string) (*dto.Movie, error) {
	movie := &dto.Movie{}
	row := m.db.QueryRow(`SELECT id, imdbID, title, year, rating, isSuperHero FROM goimdb WHERE imdbID=?`, imdbID)

	err := row.Scan(&movie.ID, &movie.ImdbID, &movie.Title, &movie.Year, &movie.Rating, &movie.IsSuperHero)

	return movie, err
}

func (m *MovieRepository) Create(movie *dto.Movie) (*dto.Movie, error) {
	stmt, err := m.db.Prepare(`
		INSERT INTO goimdb(imdbID, title, year, rating, isSuperHero)
		VALUES (?, ?, ?, ?, ?);
	`)
	if err != nil {
		return movie, err
	}
	defer stmt.Close()

	r, err := stmt.Exec(movie.ImdbID, movie.Title, movie.Year, movie.Rating, movie.IsSuperHero)

	if err != nil {
		return movie, err
	}

	id, _ := r.LastInsertId()
	movie.ID = id

	return movie, nil
}

func (m *MovieRepository) Update(id string, movie *dto.Movie) (*dto.Movie, error) {
	stmt, err := m.db.Prepare(`UPDATE goimdb SET rating = ? WHERE imdbID = ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(movie.Rating, id)
	if err != nil {
		return nil, err
	}

	return movie, nil
}

func (m *MovieRepository) Delete(id string) (bool, error) {
	stmt, err := m.db.Prepare(`DELETE FROM goimdb WHERE imdbID = ?`)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return false, err
	}

	return true, nil
}
