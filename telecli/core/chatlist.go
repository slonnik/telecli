package core

import (
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
	chatList := &ChatList{
		Box: tview.NewBox(),
	}
	chatList.Box.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyDown:
			{
				chatList.selectedChat++
			}
		case tcell.KeyUp:
			{
				chatList.selectedChat--
				if chatList.selectedChat < 0 {
					chatList.selectedChat = 0
				}
			}
		}
		return event
	})
	return chatList
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
	innerLeft, innerTop, width, _ := chatList.Box.GetInnerRect()
	x, y := innerLeft, innerTop
	for index, chat := range chatList.chats {
		title := chat.Title
		if index == chatList.selectedChat {
			for pos := 0; pos < width; pos++ {
				screen.SetContent(x+pos, y+index, 1, nil, tcell.StyleDefault.Background(tcell.ColorWhite))
			}

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
