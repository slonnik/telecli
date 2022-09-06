package core

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type tListNotification interface {
	itemChanged()
}

type ScrollableBox struct {
	*tview.Box
	selectedIndex    int
	startIndex       int
	listNotification tListNotification
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

func (scrollableBox *ScrollableBox) subscribe(notification tListNotification) *ScrollableBox {
	scrollableBox.listNotification = notification
	return scrollableBox
}

func (scrollableBox *ScrollableBox) onKeyDown(event *tcell.EventKey) {

	if scrollableBox.isScrollDownRequired() {
		scrollableBox.scrollDown()
	} else {
		scrollableBox.selectedIndex++
	}
	scrollableBox.listNotification.itemChanged()
}

func (scrollableBox *ScrollableBox) onKeyUp(event *tcell.EventKey) {
	if scrollableBox.isScrollUpRequired() {
		scrollableBox.scrollUp()
	} else {
		scrollableBox.selectedIndex--
	}
	scrollableBox.listNotification.itemChanged()
}

func (scrollableBox *ScrollableBox) scrollDown() {
	scrollableBox.startIndex++
}

func (scrollableBox *ScrollableBox) scrollUp() {
	scrollableBox.startIndex--
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
