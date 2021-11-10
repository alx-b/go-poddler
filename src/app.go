package poddler

import (
	"github.com/alx-b/go-poddler/src/tui"
)

func Start() {
	tui := tui.CreateTUI()
	defer tui.DB.CloseConnection()
	tui.InitAll()
	tui.WaitGroup.Wait()
}
