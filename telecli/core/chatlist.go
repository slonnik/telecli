package core

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type chatItem struct {
	Title string
	Id    int64
}

func (item chatItem) getHeight() int {
	return 1
}

type ChatList struct {
	*ScrollableBox
}

func NewChatList() *ChatList {

	chatList := &ChatList{
		ScrollableBox: newScrollableBox(),
	}
	chatList.setParent(chatList)
	chatList.subscribe(chatList)
	return chatList
}

func (chatList *ChatList) AddChat(title string, id int64) *ChatList {

	chatList.addRow(chatItem{
		Title: title,
		Id:    id,
	})
	return chatList
}

func (chatList *ChatList) SelectChat(index int) {
	chatList.setSelectedRow(index)
}

func (chatList ChatList) GetSelectedChatId() int64 {
	return chatList.getSelectedRow().(chatItem).Id
}

func (chatList ChatList) itemChanged() {
	CoreEvents <- NewChatSelectedEvent(chatList.GetSelectedChatId())
}

func (chatList *ChatList) Draw(screen tcell.Screen) {
	chatList.Box.DrawForSubclass(screen, chatList)
	innerLeft, innerTop, width, _ := chatList.Box.GetInnerRect()
	x, y := innerLeft, innerTop
	for index, row := range chatList.getVisibleRows() {
		title := row.(chatItem).Title
		if index == chatList.selectedIndex {
			for pos := 0; pos < width; pos++ {
				screen.SetContent(x+pos, y+index, 1, nil, tcell.StyleDefault.Background(tcell.ColorGray))
			}

		}
		tview.Print(screen, substr(title, width), x, y+index, 100, 0, tcell.ColorOlive)
	}
}

func substr(input string, length int) string {
	asRunes := []rune(input)

	if length > len(asRunes) {
		length = len(asRunes)
	}

	return string(asRunes[0:length])
}
