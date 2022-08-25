package core

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type chatItem struct {
	Title string
	Id    int64
}

type ChatList struct {
	*tview.Box
	chats         []*chatItem
	selectedIndex int
	startIndex    int
}

func NewChatList() *ChatList {
	chatList := &ChatList{
		Box: tview.NewBox(),
	}
	chatList.Box.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyDown:
			{
				chatList.onKeyDown(event)
			}
		case tcell.KeyUp:
			{
				chatList.onKeyUp(event)
			}
		}
		return event
	})
	return chatList
}

func (chatList *ChatList) onKeyDown(event *tcell.EventKey) {

	if chatList.isScrollDownRequired() {
		chatList.scrollDown()
	}
	chat := chatList.selectNext()
	PublishEvent(NewChatSelectedEvent(chat.Id))
}

func (chatList *ChatList) onKeyUp(event *tcell.EventKey) {
	if chatList.isScrollUpRequired() {
		chatList.scrollUp()
	}
	chat := chatList.selectPrevious()
	PublishEvent(NewChatSelectedEvent(chat.Id))
}

func (chatList *ChatList) AddChat(title string, id int64) *ChatList {

	chat := &chatItem{
		Title: title,
		Id:    id,
	}

	chatList.chats = append(chatList.chats, chat)
	return chatList
}

func (chatList *ChatList) SelectChat(index int) {
	chatList.selectedIndex = index
}

func (chatList *ChatList) GetSelectedChatId() int64 {
	return chatList.chats[chatList.selectedIndex].Id
}

func (chatList *ChatList) Draw(screen tcell.Screen) {
	chatList.Box.DrawForSubclass(screen, chatList)
	innerLeft, innerTop, width, height := chatList.Box.GetInnerRect()
	x, y := innerLeft, innerTop
	for index, chat := range chatList.getVisibleItems(height) {
		title := chat.Title
		if index == chatList.selectedIndex {
			for pos := 0; pos < width; pos++ {
				screen.SetContent(x+pos, y+index, 1, nil, tcell.StyleDefault.Background(tcell.ColorWhite))
			}

		}

		tview.Print(screen, substr(title, width), x, y+index, 100, 0, tcell.ColorOlive)
	}
}

func (chatList *ChatList) SetBorder(show bool) *ChatList {
	chatList.Box.SetBorder(show)
	return chatList
}

func (chatList *ChatList) SetTitle(title string) *ChatList {
	chatList.Box.SetTitle(title)
	return chatList
}

func (chatList *ChatList) SetFocus() *ChatList {
	chatList.Box.Focus(nil)
	return chatList
}

func (chatList *ChatList) getVisibleItems(height int) []*chatItem {
	return chatList.chats[chatList.startIndex : chatList.startIndex+height]
}

func (chatList *ChatList) scrollDown() {
	chatList.startIndex++
	chatList.selectedIndex--
}

func (chatList *ChatList) scrollUp() {
	chatList.startIndex--
	chatList.selectedIndex++
	if chatList.startIndex < 0 {
		chatList.startIndex = 0
	}
}

func (chatList *ChatList) isScrollDownRequired() bool {
	_, _, _, height := chatList.Box.GetInnerRect()
	return chatList.selectedIndex == height-1
}

func (chatList *ChatList) isScrollUpRequired() bool {
	return chatList.selectedIndex == 0
}

func (chatList *ChatList) selectNext() *chatItem {
	chatList.selectedIndex++
	return chatList.chats[chatList.startIndex+chatList.selectedIndex]
}

func (chatList *ChatList) selectPrevious() *chatItem {
	chatList.selectedIndex--
	if chatList.selectedIndex < 0 {
		chatList.selectedIndex = 0
	}
	return chatList.chats[chatList.startIndex+chatList.selectedIndex]
}

func substr(input string, length int) string {
	asRunes := []rune(input)

	if length > len(asRunes) {
		length = len(asRunes)
	}

	return string(asRunes[0:length])
}
