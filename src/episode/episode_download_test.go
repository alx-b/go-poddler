package episode

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/mmcdole/gofeed"
)

func TestGetEpisodes(t *testing.T) {
	testFeed := gofeed.Feed{
		Title:       "podcast title1",
		FeedType:    "SomeType",
		FeedVersion: "SomeVersion",
		Items: []*gofeed.Item{
			{
				Enclosures: []*gofeed.Enclosure{{URL: "http: link"}},
				Title:      "title1",
				Published:  "date1",
			},
			{
				Enclosures: []*gofeed.Enclosure{{URL: "http: link"}},
				Title:      "title2",
				Published:  "date2",
			},
			{
				Enclosures: []*gofeed.Enclosure{{URL: "http: link"}},
				Title:      "title3",
				Published:  "date3",
			},
		},
	}

	t.Run("Get Episodes from Feed", func(t *testing.T) {
		got := GetEpisodes(&testFeed)
		want := []Episode{
			{PodcastTitle: "podcast title1", Number: 3, Title: "title1", FileURL: "http: link", Date: "date1"},
			{PodcastTitle: "podcast title1", Number: 2, Title: "title2", FileURL: "http: link", Date: "date2"},
			{PodcastTitle: "podcast title1", Number: 1, Title: "title3", FileURL: "http: link", Date: "date3"},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func TestFileExists(t *testing.T) {
	t.Run("check that file exists at the path", func(t *testing.T) {
		file, err := ioutil.TempFile("./", "temp")
		if err != nil {
			t.Errorf("%s", err)
		}
		fmt.Println(file.Name())
		defer os.Remove(file.Name())

		got := fileExists(file.Name())
		want := true

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("check that file does NOT exists at the path", func(t *testing.T) {
		got := fileExists("./temp434234324")
		want := false

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}
