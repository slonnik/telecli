package core

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"time"
)

type tListItem struct {
	MainText      string    // The main text of the list item.
	SecondaryText string    // A secondary text to be shown underneath the main text.
	CreationTime  time.Time // item date-time creation
}

func (item tListItem) getHeight() int {
	return 2
}

type TeleList struct {
	*ScrollableBox
}

func NewTeleList() *TeleList {
	teleList := &TeleList{
		ScrollableBox: newScrollableBox(),
	}
	teleList.subscribe(teleList)
	teleList.setParent(teleList)
	return teleList
}

func (teleList *TeleList) AddItem(mainText, secondaryText string, timeStamp int64) *TeleList {

	item := &tListItem{
		MainText:      mainText,
		SecondaryText: secondaryText,
		CreationTime:  time.Unix(timeStamp, 0),
	}

	teleList.addRow(item)
	return teleList
}

func (teleList *TeleList) ClearItems() *TeleList {
	teleList.clearRows()
	return teleList
}

func (teleList *TeleList) Draw(screen tcell.Screen) {
	teleList.Box.DrawForSubclass(screen, teleList)
	innerLeft, innerTop, width, _ := teleList.Box.GetInnerRect()
	x, currentYPos := innerLeft, innerTop

	for index, row := range teleList.getVisibleRows() {
		item := row.(*tListItem)
		if index == teleList.selectedIndex {
			for pos := 0; pos < width; pos++ {
				screen.SetContent(x+pos, currentYPos, 1, nil, tcell.StyleDefault.Background(tcell.ColorGray))
				screen.SetContent(x+pos, currentYPos+1, 1, nil, tcell.StyleDefault.Background(tcell.ColorGray))
			}

		}
		tview.Print(screen, item.MainText, x, currentYPos, 100, 0, tcell.ColorOlive)
		year, month, day := item.CreationTime.Date()
		hour, minute, _ := item.CreationTime.Clock()
		timeText := fmt.Sprintf("%v-%v-%v %v:%v", year, month, day, hour, minute)
		tview.Print(screen, fmt.Sprintf("[%v] %v", timeText, item.SecondaryText), x+4, currentYPos+1, 100, 0, tcell.ColorGreen)
		currentYPos += 2
	}
}

func (teleList *TeleList) SetTitle(title string) *TeleList {
	teleList.Box.SetTitle(fmt.Sprintf("= %v =", title))
	return teleList
}

func (teleList *TeleList) itemChanged() {
	PublishEvents(NewUpdateScreenEvent())
}
