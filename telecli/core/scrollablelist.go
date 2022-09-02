package core

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ScrollableBox struct {
	*tview.Box
	selectedIndex int
	startIndex    int
}

func newScrollableBox() *ScrollableBox {
	scrollableBox := &ScrollableBox{
		Box: tview.NewBox(),
	}

	scrollableBox.Box.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyDown:
			{
				scrollableBox.onKeyDown(event)
			}
		case tcell.KeyUp:
			{
				scrollableBox.onKeyUp(event)
			}
		}
		return event
	})
	return scrollableBox
}

func (scrollableBox *ScrollableBox) onKeyDown(event *tcell.EventKey) {

	if scrollableBox.isScrollDownRequired() {
		scrollableBox.scrollDown()
	}

}

func (scrollableBox *ScrollableBox) onKeyUp(event *tcell.EventKey) {
	if scrollableBox.isScrollUpRequired() {
		scrollableBox.scrollUp()
	}
}

func (scrollableBox *ScrollableBox) scrollDown() {
	scrollableBox.startIndex++
	scrollableBox.selectedIndex--
}

func (scrollableBox *ScrollableBox) scrollUp() {
	scrollableBox.startIndex--
	scrollableBox.selectedIndex++
	if scrollableBox.startIndex < 0 {
		scrollableBox.startIndex = 0
	}
}

func (scrollableBox *ScrollableBox) isScrollDownRequired() bool {
	_, _, _, height := scrollableBox.Box.GetInnerRect()
	return scrollableBox.selectedIndex == height-1
}

func (scrollableBox *ScrollableBox) isScrollUpRequired() bool {
	return scrollableBox.selectedIndex == 0
}

/*func (chatList *ScrollableBox) selectNext() *chatItem {
	chatList.selectedIndex++
	return chatList.chats[chatList.startIndex+chatList.selectedIndex]
}

func (chatList *ScrollableBox) selectPrevious() *chatItem {
	chatList.selectedIndex--
	if chatList.selectedIndex < 0 {
		chatList.selectedIndex = 0
	}
	return chatList.chats[chatList.startIndex+chatList.selectedIndex]
}*/
