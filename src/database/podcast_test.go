package database

import (
	"reflect"
	"testing"
)

func TestPodcastDatabaseQueries(t *testing.T) {
	db := CreateDB(":memory:")
	defer db.CloseConnection()
	db.SavePodcast(PodcastIn{Title: "title1", URL: "http: something1"})
	db.SavePodcast(PodcastIn{Title: "title2", URL: "http: something2"})

	t.Run("Get all podcasts from database", func(t *testing.T) {
		got, _ := db.GetPodcasts()
		want := []PodcastOut{
			{
				ID:    1,
				Title: "title1",
				URL:   "http: something1",
			},
			{
				ID:    2,
				Title: "title2",
				URL:   "http: something2",
			},
		}
		assertPodcastListDeepEqual(t, got, want)
	})

	t.Run("Add podcast to database", func(t *testing.T) {
		db.SavePodcast(PodcastIn{
			Title: "title3",
			URL:   "http: something3",
		})
		got, _ := db.GetPodcasts()
		want := []PodcastOut{
			{
				ID:    1,
				Title: "title1",
				URL:   "http: something1",
			},
			{
				ID:    2,
				Title: "title2",
				URL:   "http: something2",
			},
			{
				ID:    3,
				Title: "title3",
				URL:   "http: something3",
			},
		}
		assertPodcastListDeepEqual(t, got, want)
	})

	t.Run("Delete podcast from database", func(t *testing.T) {
		db.DeletePodcast(2)
		got, _ := db.GetPodcasts()
		want := []PodcastOut{
			{
				ID:    1,
				Title: "title1",
				URL:   "http: something1",
			},
			{
				ID:    3,
				Title: "title3",
				URL:   "http: something3",
			},
		}
		assertPodcastListDeepEqual(t, got, want)
	})

}

func assertPodcastListDeepEqual(t testing.TB, got, want []PodcastOut) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
