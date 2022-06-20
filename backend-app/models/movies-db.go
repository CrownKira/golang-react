package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type DBModel struct {
	DB *sql.DB
}

// exported name: capital letter
// https://stackoverflow.com/questions/38616687/which-way-to-name-a-function-in-go-camelcase-or-semi-camelcase/38617771
// Exported names (that is, identifiers that can be used from a package other than the one where they are defined) begin with a capital letter. Thus if your method is part of your public API, it should be written:

// WriteToDB

// but if it is an internal helper method it should be written:

// writeToDB
// get one movie from db
// Get returns one movie and error, if any
// qn: what is context?
// m: DB model is the receiver
func (m *DBModel) Get(id int) (*Movie, error) {
	// when work with database, must make use of context
	// this db will timeout after 3 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// query the db
	// $1: placeholder
	query := `select id, title, description, year, release_date, runtime, rating, mpaa_rating, 
						created_at, updated_at from movies where id = $1
	`
	// get one row and exit if more than 3 sec, so need to pass the timeout
	// substitute id for the $1
	// get row struct
	row := m.DB.QueryRowContext(ctx, query, id)

	var movie Movie
	// scan the sql row
	// scan to movie struct var
	// copy from one struct to another struct
	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.Description,
		&movie.Year,
		&movie.ReleaseDate,
		&movie.Runtime,
		&movie.Rating,
		&movie.MPAARating,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// get the genres
	// alias: mg: movie genre
	// g: genre
	// select from ... (movies genres <alias> ) where ... (the id)
	// qn: what is left join; ans: combine to the movies genres table
	// select syntax:
	// select <table>.<field> from <left_table> <alias> left join <right_table> <alias> on (<left_table>.<field> = <right_table>.<field>) where [filter the table]
	query = `select
						mg.id, mg.movie_id, mg.genre_id, g.genre_name
					from 
						movies_genres mg
						left join genres g on (g.id = mg.genre_id)
					where
					mg.movie_id = $1
	`

	// get all the rows, not just one
	// use next to advance from row to row
	rows, _ := m.DB.QueryContext(ctx, query, id)
	// qn: why need close?
	defer rows.Close() // so no resource leak

	// fmt.Println("rows is:", rows)

	var genres = make(map[int]string)
	// rows.Next() ?
	for rows.Next() {
		var mg MovieGenre
		err := rows.Scan(
			&mg.ID,
			&mg.MovieID,
			&mg.GenreID,
			&mg.Genre.GenreName,
		)
		if err != nil {
			return nil, err
		}
		// genres = append(genres, mg)
		// movie genre id : genre name
		genres[mg.ID] = mg.Genre.GenreName
	}

	movie.MovieGenre = genres

	return &movie, nil
}

// get all db from db
// All returns all movies and error, if any
func (m *DBModel) All(genre ...int) ([]*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	where := ""
	if len(genre) > 0 {
		// qn: take only one genre ?
		where = fmt.Sprintf("where id in (select movie_id from movies_genres where genre_id = %d)", genre[0])
	}

	// select <table>.<field> from <left_table> <alias> left join <right_table> <alias> on (<left_table>.<field> = <right_table>.<field>) where (...) order by [field]
	query := fmt.Sprintf(`select id, title, description, year, release_date, runtime, rating, mpaa_rating, 
						created_at, updated_at from movies %s order by title`, where)

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []*Movie

	for rows.Next() {
		var movie Movie
		err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.Description,
			&movie.Year,
			&movie.ReleaseDate,
			&movie.Runtime,
			&movie.Rating,
			&movie.MPAARating,
			&movie.CreatedAt,
			&movie.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		genreQuery := `select
		mg.id, mg.movie_id, mg.genre_id, g.genre_name
	from 
		movies_genres mg
		left join genres g on (g.id = mg.genre_id)
	where
	mg.movie_id = $1
`

		// get all the rows, not just one
		// use next to advance from row to row
		genreRows, _ := m.DB.QueryContext(ctx, genreQuery, movie.ID)
		// qn: why need close?
		// defer rows.Close() // so no resource leak

		// fmt.Println("rows is:", rows)

		var genres = make(map[int]string)
		// rows.Next() ?
		for genreRows.Next() {
			var mg MovieGenre
			err := genreRows.Scan(
				&mg.ID,
				&mg.MovieID,
				&mg.GenreID,
				&mg.Genre.GenreName,
			)
			if err != nil {
				return nil, err
			}
			// genres = append(genres, mg)
			// movie genre id : genre name
			genres[mg.ID] = mg.Genre.GenreName
		}

		genreRows.Close()
		movie.MovieGenre = genres
		movies = append(movies, &movie)
	}

	return movies, nil
}

func (m *DBModel) GenresAll() ([]*Genre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, genre_name, created_at, updated_at from genres order by genre_name`

	// query context returns the rows (iterator with next())
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	// close the rows
	defer rows.Close()

	var genres []*Genre

	for rows.Next() {
		var g Genre
		// scan and append one by one
		err := rows.Scan(
			&g.ID,
			&g.GenreName,
			&g.CreatedAt,
			&g.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		genres = append(genres, &g)
	}

	return genres, nil
}
