package core

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type listItem struct {
	MainText      string // The main text of the list item.
	SecondaryText string // A secondary text to be shown underneath the main text.
}

type TeleList struct {
	*tview.Box
	items []*listItem
}

func NewTeleList() *TeleList {
	return &TeleList{
		Box: tview.NewBox(),
	}
}

func (teleList *TeleList) AddItem(mainText, secondaryText string) *TeleList {

	item := &listItem{
		MainText:      mainText,
		SecondaryText: secondaryText,
	}

	teleList.items = append(teleList.items, item)
	return teleList
}

func (teleList *TeleList) Draw(screen tcell.Screen) {
	teleList.Box.DrawForSubclass(screen, teleList)
	if len(teleList.items) == 0 {
		return
	}
	for index, item := range teleList.items {
		tview.Print(screen, item.MainText, 0, index*2, 100, 0, tcell.ColorWhite)
		tview.Print(screen, item.SecondaryText, 4, index*2+1, 100, 0, tcell.ColorGreen)
	}

}
