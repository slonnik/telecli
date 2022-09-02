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
	*ScrollableBox
	chats []*chatItem
}

func NewChatList() *ChatList {

	return &ChatList{
		ScrollableBox: newScrollableBox(),
	}
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

func (chatList *ChatList) getVisibleItems(height int) []*chatItem {
	return chatList.chats[chatList.startIndex : chatList.startIndex+height]
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

func substr(input string, length int) string {
	asRunes := []rune(input)

	if length > len(asRunes) {
		length = len(asRunes)
	}

	return string(asRunes[0:length])
}
