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
	left, top, _, _ := teleList.Box.GetRect()
	innerLeft, innerTop, _, _ := teleList.Box.GetInnerRect()
	x, y := left+innerLeft, top+innerTop
	for index, item := range teleList.items {
		tview.Print(screen, item.MainText, x, y+index*2, 100, 0, tcell.ColorOlive)
		tview.Print(screen, item.SecondaryText, x+4, y+index*2+1, 100, 0, tcell.ColorGreen)
	}
}

func (teleList *TeleList) SetBorder(show bool) *TeleList {
	teleList.Box.SetBorder(show)
	return teleList
}
