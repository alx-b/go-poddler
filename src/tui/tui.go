package tui

import (
	"fmt"
	"sync"

	"github.com/alx-b/go-poddler/src/database"
	"github.com/alx-b/go-poddler/src/episode"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TUI struct {
	App         *tview.Application
	UrlField    *tview.InputField
	PodcastList *tview.List
	EpisodeList *tview.List
	StatusList  *tview.List
	DB          *database.Database
	WaitGroup   sync.WaitGroup
	Podcasts    []database.PodcastOut
	Grid        *tview.Grid
}

func (t *TUI) InitAll() {
	t.Podcasts, _ = t.DB.GetPodcasts()
	t.initGrid()
	t.setInputCapture()
	t.setDoneFunc()
	t.loadList()

	if err := t.App.SetRoot(t.Grid, true).EnableMouse(false).SetFocus(t.PodcastList).Run(); err != nil {
		panic(err)
	}
}

func (t *TUI) setDoneFunc() {
	t.UrlField.SetDoneFunc(func(key tcell.Key) {
		if key != tcell.KeyEnter {
			// Do nothing
		} else {
			url := t.UrlField.GetText()
			feed, err := episode.GetPodcastFeed(url)

			if err != nil {
				fmt.Errorf("%v", err)
			}

			t.DB.SavePodcast(database.PodcastIn{Title: feed.Title, URL: url})

			t.UrlField.SetText("")
			t.PodcastList.Clear()
			t.loadList()
		}
	})
}

func (t *TUI) initGrid() {
	t.Grid = tview.NewGrid().
		SetRows(3, 0, 0, 0).
		SetColumns(0).
		SetBorders(false).
		AddItem(t.UrlField, 0, 0, 1, 1, 0, 0, false).
		AddItem(t.PodcastList, 1, 0, 1, 1, 0, 0, false).
		AddItem(t.EpisodeList, 2, 0, 1, 1, 0, 0, false).
		AddItem(t.StatusList, 3, 0, 1, 1, 0, 0, false)
}

func (t *TUI) setInputCapture() {
	t.App.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		if e.Key() == tcell.KeyTAB {
			switch {
			case t.UrlField.HasFocus():
				t.App.SetFocus(t.PodcastList)
			case t.PodcastList.HasFocus():
				t.App.SetFocus(t.EpisodeList)
			case t.EpisodeList.HasFocus():
				t.App.SetFocus(t.StatusList)
			case t.StatusList.HasFocus():
				t.App.SetFocus(t.UrlField)
			}
		}
		if e.Key() == tcell.KeyEsc {
			t.App.Stop()
		}
		return e
	})

	t.PodcastList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab || event.Key() == tcell.KeyBacktab {
			return nil
		}
		if event.Rune() == 'x' {
			item_idx := t.PodcastList.GetCurrentItem()
			t.DB.DeletePodcast(t.Podcasts[item_idx].ID)
			t.PodcastList.RemoveItem(item_idx)
			t.Podcasts, _ = t.DB.GetPodcasts()
			return nil
		}
		return event
	})

	t.EpisodeList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab || event.Key() == tcell.KeyBacktab {
			return nil
		}
		return event
	})

	t.StatusList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			return nil
		}
		return event
	})
}

func (t *TUI) loadList() {
	t.Podcasts, _ = t.DB.GetPodcasts()

	for _, podcast := range t.Podcasts {
		// When Enter is pressed, get episodes list of a podcast
		t.PodcastList.AddItem(podcast.Title, "", 0, func() {
			t.EpisodeList.Clear()
			pod := t.Podcasts[t.PodcastList.GetCurrentItem()]

			feed, err := episode.GetPodcastFeed(pod.URL)

			if err != nil {
				fmt.Errorf("%v", err)
			}

			episodes := episode.GetEpisodes(feed)
			for _, ep := range episodes {
				// Add said episodes to the episodeList
				// When Enter is pressed, download selected episode
				t.EpisodeList.AddItem(fmt.Sprintf("%d. %s (%s)", ep.Number, ep.Title, ep.Date), "", 0, func() {
					anEpisode := episodes[t.EpisodeList.GetCurrentItem()]
					t.WaitGroup.Add(1)
					go func() {
						t.StatusList.AddItem(fmt.Sprintf("Downloading %d._%s", anEpisode.Number, anEpisode.Title), "", 0, nil)
						index := t.StatusList.GetItemCount() - 1
						episode.DownloadFile(anEpisode)
						t.StatusList.SetItemText(index, fmt.Sprintf("Done %d._%s", anEpisode.Number, anEpisode.Title), "")
						//t.StatusList.AddItem(fmt.Sprintf("Done %d._%s", episode.Number, episode.Title), "", 0, nil)
						t.WaitGroup.Done()
					}()
				})
			}
		})
	}
}

func initUrlField() *tview.InputField {
	inputField := tview.NewInputField().
		SetLabel("Enter url: ").
		SetFieldWidth(80)

	inputField.SetFieldBackgroundColor(tcell.ColorGray).
		SetBorder(true).
		SetTitle("Add new feed URL").
		SetBorderPadding(0, 0, 1, 1)

	return inputField
}

func initPodcastList() *tview.List {
	podcastList := tview.NewList().ShowSecondaryText(false)
	podcastList.SetTitle("Podcasts").SetBorder(true).SetBorderPadding(0, 0, 1, 1)
	return podcastList
}

func initEpisodeList() *tview.List {
	episodeList := tview.NewList().ShowSecondaryText(false)
	episodeList.SetTitle("Episodes").SetBorder(true).SetBorderPadding(0, 0, 1, 1)
	return episodeList
}

func initStatusList() *tview.List {
	statusList := tview.NewList().ShowSecondaryText(false)
	statusList.SetTitle("Status").SetBorder(true).SetBorderPadding(0, 0, 1, 1)
	return statusList
}

func CreateTUI() TUI {
	return TUI{
		App:         tview.NewApplication(),
		UrlField:    initUrlField(),
		PodcastList: initPodcastList(),
		EpisodeList: initEpisodeList(),
		StatusList:  initStatusList(),
		DB:          database.CreateDB("poddler_db"),
		Podcasts:    []database.PodcastOut{},
	}
}
