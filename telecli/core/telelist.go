package core

import (
	"fmt"
	"github.com/Arman92/go-tdlib"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"time"
)

type tItemStyle struct {
	headerColor        tcell.Color
	bodyColor          tcell.Color
	creationTimeFormat string
}

var (
	greyStyle = tItemStyle{
		headerColor:        tcell.ColorDarkGray,
		bodyColor:          tcell.ColorDarkGray,
		creationTimeFormat: "02 Jan 06 15:04 MST",
	}
	mainStyle = tItemStyle{
		headerColor:        tcell.ColorOlive,
		bodyColor:          tcell.ColorGreen,
		creationTimeFormat: "02 Jan 06 15:04 MST",
	}
)

type tListItem struct {
	text         string    // The main text of the list item.
	CreationTime time.Time // item date-time creation
	style        tItemStyle
	height       int
}

func (item tListItem) getHeight() int {
	return item.height
}

func (item *tListItem) setHeight(height int) {
	item.height = height
}

func (item tListItem) getCreationTimeAsText() string {
	return item.CreationTime.Format(item.style.creationTimeFormat)
}

type TMessage tdlib.Message

func (message TMessage) ToListItem() *tListItem {
	var eventText string
	var style tItemStyle
	switch message.Content.GetMessageContentEnum() {
	case tdlib.MessageTextType:
		eventText = message.Content.(*tdlib.MessageText).Text.Text
		style = mainStyle
	default:
		eventText = string(message.Content.GetMessageContentEnum())
		style = greyStyle
	}

	return &tListItem{
		text:         eventText,
		CreationTime: time.Unix(int64(message.Date), 0),
		style:        style,
	}
}

type tPrinter struct {
	screen    tcell.Screen
	curOffset int
	curRow    int
	selected  bool
	width     int
}

func (printer *tPrinter) print(text string, color tcell.Color) *tPrinter {
	return printer.printWithIntent(text, 0, color)
}

func (printer *tPrinter) printWithIntent(text string, intent int, color tcell.Color) *tPrinter {
	if printer.selected {
		for i := printer.curOffset; i < printer.width; i++ {
			printer.screen.SetContent(i, printer.curRow, 1, nil, tcell.StyleDefault.Background(tcell.ColorGray))
		}
	}
	tview.Print(printer.screen, text, printer.curOffset+intent, printer.curRow, len(text), 0, color)
	printer.curRow++
	return printer
}

func (printer *tPrinter) setSelected(selected bool) *tPrinter {
	printer.selected = selected
	return printer
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

func (teleList *TeleList) AddItem(item *tListItem) {
	teleList.addRow(item)
}

func (teleList *TeleList) ClearItems() *TeleList {
	teleList.clearRows()
	return teleList
}

func (teleList *TeleList) Draw(screen tcell.Screen) {
	teleList.Box.DrawForSubclass(screen, teleList)
	const intent = 4
	innerLeft, innerTop, width, _ := teleList.Box.GetInnerRect()

	printer := &tPrinter{
		screen:    screen,
		curRow:    innerTop,
		curOffset: innerLeft,
		width:     width,
	}

	for index, row := range teleList.getVisibleRows() {
		item := row.(*tListItem)
		textRows := chunkText(item.text, width-intent)
		item.setHeight(1 + len(textRows))

		printer.setSelected(index == teleList.selectedIndex)
		printer.print(item.getCreationTimeAsText(), item.style.headerColor)
		for _, text := range textRows {
			printer.printWithIntent(text, intent, item.style.bodyColor)
		}
	}
}

func (teleList *TeleList) SetTitle(title string) *TeleList {
	teleList.Box.SetTitle(fmt.Sprintf("= %v =", title))
	return teleList
}

func (teleList *TeleList) itemChanged() {
	CoreEvents <- NewUpdateScreenEvent()
}

func chunkText(text string, chunkSize int) []string {
	if len(text) == 0 {
		return nil
	}
	if chunkSize >= len(text) {
		return []string{text}
	}
	var chunks []string = make([]string, 0, (len(text)-1)/chunkSize+1)
	currentLen := 0
	currentStart := 0
	for i := range text {
		if currentLen == chunkSize {
			chunks = append(chunks, text[currentStart:i])
			currentLen = 0
			currentStart = i
		}
		currentLen++
	}
	chunks = append(chunks, text[currentStart:])
	return chunks
}
