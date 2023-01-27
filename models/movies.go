package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

type Movie struct {
    Id       int     `json:"id"`
    Title    string  `json:"title"`
    Year     int     `json:"year"`
    Director string  `json:"director"`
}


func ConnectDatabase() error {
	db, err := sql.Open("sqlite3", "./sqlite.db")
	if err != nil {
		return err
	}
	DB = db
	return nil
}

func GetMovies() ([]Movie, error) {
	rows, err := DB.Query("SELECT id, title, year, director FROM movies")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	movies := make([]Movie, 0)
	for rows.Next() {
		singleMovie := Movie{}
		err = rows.Scan(&singleMovie.Id, &singleMovie.Title, &singleMovie.Year, &singleMovie.Director)

		if err != nil {
			return nil, err
		}

		movies = append(movies, singleMovie)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return movies, err
}

func GetMovieById(id string) (Movie, error) {

	stmt, err := DB.Prepare("SELECT id, title, year, director FROM movies WHERE id = ?")

	if err != nil {
		return Movie{}, err
	}

	movie := Movie{}

	sqlErr := stmt.QueryRow(id).Scan(&movie.Id, &movie.Title, &movie.Year, &movie.Director)

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return Movie{}, nil
		}
		return Movie{}, sqlErr
	}
	return movie, nil
}

func AddMovie(newMovie Movie) (bool, error) {
	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("INSERT INTO movies (title, year, director) VALUES (?, ?, ?)")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(newMovie.Title, newMovie.Year, newMovie.Director)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func UpdateMovie(ourMovie Movie, id int) (bool, error) {

	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("UPDATE movies SET title = ?, year = ?, director = ? WHERE Id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(ourMovie.Title, ourMovie.Year, ourMovie.Director, id)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func DeleteMovie(movieId int) (bool, error) {

	tx, err := DB.Begin()

	if err != nil {
		return false, err
	}

	stmt, err := DB.Prepare("DELETE FROM movies WHERE id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(movieId)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}
