package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func LoadUI(favorites []string) {
	app := tview.NewApplication()

	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

	grid := tview.NewGrid().SetBorders(true)

	for index, fave := range favorites {
		item := newPrimitive(fave)
		grid.AddItem(item, index/3, index%3, 1, 1, 0, 0, false)
	}

	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == 'q' {
			app.Stop()
		}
		return nil
	})

	if err := app.SetRoot(grid, true).SetFocus(grid).Run(); err != nil {
		panic(err)
	}
}
