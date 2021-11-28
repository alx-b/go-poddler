package database

import (
	"fmt"
	"sort"
)

type PodcastIn struct {
	Title string
	URL   string
}

type PodcastOut struct {
	ID    int
	Title string
	URL   string
}

func (db *Database) SortPodcastByTitle() []PodcastOut {
	podcasts, err := db.GetPodcasts()

	if err != nil {
		panic(err)
	}

	sort.Slice(podcasts, func(a, b int) bool {
		return podcasts[a].Title < podcasts[b].Title
	})

	return podcasts
}

func (db *Database) GetPodcasts() ([]PodcastOut, error) {
	rows, err := db.conn.Query("SELECT * FROM podcast")

	if err != nil {
		return nil, fmt.Errorf("%+v: %v", ErrCouldNotQueryDatabase, err)
	}

	list := []PodcastOut{}

	for rows.Next() {
		podcast := PodcastOut{}
		rows.Scan(
			&podcast.ID,
			&podcast.Title,
			&podcast.URL,
		)
		list = append(list, podcast)
	}

	return list, nil
}

func (db *Database) SavePodcast(newPodcast PodcastIn) error {
	_, err := db.conn.Exec(
		"INSERT INTO podcast (title, url) VALUES (?,?)",
		newPodcast.Title,
		newPodcast.URL,
	)

	if err != nil {
		return ErrCouldNotExecuteQuery
	}

	return nil
}

func (db *Database) DeletePodcast(idx int) error {
	result, err := db.conn.Exec("DELETE FROM podcast WHERE id=?", idx)

	if err != nil {
		return ErrCouldNotDeleteFromDatabase
	}

	numRow, _ := result.RowsAffected()

	if numRow < 1 {
		return ErrNoValueModified

	}

	return nil
}
