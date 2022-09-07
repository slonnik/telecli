package core

import (
	"fmt"
	"github.com/Arman92/go-tdlib"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"time"
)

type tItemStyle struct {
	headerColor tcell.Color
	bodyColor   tcell.Color
}

var greyStyle = tItemStyle{
	headerColor: tcell.ColorDarkGray,
	bodyColor:   tcell.ColorDarkGray,
}

var mainStyle = tItemStyle{
	headerColor: tcell.ColorOlive,
	bodyColor:   tcell.ColorGreen,
}

type tListItem struct {
	text         string    // The main text of the list item.
	CreationTime time.Time // item date-time creation
	style        tItemStyle
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

func (teleList *TeleList) AddItemFromMessage(message tdlib.Message) {

	var eventText string
	var style tItemStyle = greyStyle
	switch message.Content.GetMessageContentEnum() {
	case tdlib.MessageTextType:
		eventText = message.Content.(*tdlib.MessageText).Text.Text
		style = mainStyle
	default:
		eventText = string(message.Content.GetMessageContentEnum())
	}

	item := &tListItem{
		text:         eventText,
		CreationTime: time.Unix(int64(message.Date), 0),
		style:        style,
	}

	teleList.addRow(item)
}

/*func (teleList *TeleList) AddItem(mainText, secondaryText string, timeStamp int64) {

	item := &tListItem{
		text:      mainText,
		SecondaryText: secondaryText,
		CreationTime:  time.Unix(timeStamp, 0),
	}

	teleList.addRow(item)
}*/

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
		year, month, day := item.CreationTime.Date()
		hour, minute, _ := item.CreationTime.Clock()
		timeText := fmt.Sprintf("%v-%v-%v %v:%v", year, month, day, hour, minute)
		tview.Print(screen, timeText, x, currentYPos, 100, 0, item.style.headerColor)
		tview.Print(screen, item.text, x+4, currentYPos+1, 100, 0, item.style.bodyColor)
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
