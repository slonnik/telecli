package core

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"time"
)

type listItem struct {
	MainText      string  // The main text of the list item.
	SecondaryText string  // A secondary text to be shown underneath the main text.
	TimeStamp     float64 // item date-time creation
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

func (teleList *TeleList) AddItem(mainText, secondaryText string, timeStamp float64) *TeleList {

	item := &listItem{
		MainText:      mainText,
		SecondaryText: secondaryText,
		TimeStamp:     timeStamp,
	}

	teleList.items = append(teleList.items, item)
	return teleList
}

func (teleList *TeleList) ClearItems() *TeleList {

	teleList.items = nil
	return teleList
}

func (teleList *TeleList) Draw(screen tcell.Screen) {
	teleList.Box.DrawForSubclass(screen, teleList)
	innerLeft, innerTop, _, _ := teleList.Box.GetInnerRect()
	x, y := innerLeft, innerTop
	for index, item := range teleList.items {
		tview.Print(screen, item.MainText, x, y+index*2, 100, 0, tcell.ColorOlive)
		itemTime := time.Unix(int64(item.TimeStamp), 0).Local()
		year, month, day := itemTime.Date()
		timeText := fmt.Sprintf("%v-%v-%v %v:%v", year, month, day, itemTime.Hour(), itemTime.Minute())
		tview.Print(screen, fmt.Sprintf("[%v] %v", timeText, item.SecondaryText), x+4, y+index*2+1, 100, 0, tcell.ColorGreen)
	}
}

func (teleList *TeleList) SetBorder(show bool) *TeleList {
	teleList.Box.SetBorder(show)
	return teleList
}

func (teleList *TeleList) SetTitle(title string) *TeleList {
	teleList.Box.SetTitle(fmt.Sprintf("= %v =", title))
	return teleList
}
