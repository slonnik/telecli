package core

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type chatItem struct {
	Title string
}

type ChatList struct {
	*tview.Box
	chats        []*chatItem
	selectedChat int
}

func NewChatList() *ChatList {
	return &ChatList{
		Box: tview.NewBox(),
	}
}

func (chatList *ChatList) AddChat(title string) *ChatList {

	chat := &chatItem{
		Title: title,
	}

	chatList.chats = append(chatList.chats, chat)
	return chatList
}

func (chatList *ChatList) SelectChat(index int) {
	chatList.selectedChat = index
}

func (chatList *ChatList) Draw(screen tcell.Screen) {
	chatList.Box.DrawForSubclass(screen, chatList)
	innerLeft, innerTop, _, _ := chatList.Box.GetInnerRect()
	x, y := innerLeft, innerTop
	for index, chat := range chatList.chats {
		title := chat.Title
		if index == chatList.selectedChat {
			title = fmt.Sprintf(`[red]%v[white]`, title)
		}
		tview.Print(screen, title, x, y+index, 100, 0, tcell.ColorOlive)
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
