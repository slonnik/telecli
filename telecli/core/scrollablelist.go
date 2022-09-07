package core

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type tListNotification interface {
	itemChanged()
}

type tRow interface {
	getHeight() int
}

type ScrollableBox struct {
	*tview.Box
	selectedIndex    int
	startIndex       int
	listNotification tListNotification
	rows             []tRow
	parent           tview.Primitive
}

func newScrollableBox() *ScrollableBox {
	scrollableBox := &ScrollableBox{
		Box: tview.NewBox(),
	}

	scrollableBox.Box.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		switch action {
		case tview.MouseLeftClick:
			scrollableBox.onMouseLeftClick(event)
		}
		return action, event
	})
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

func (scrollableBox *ScrollableBox) setParent(parent tview.Primitive) {
	scrollableBox.parent = parent
}

func (scrollableBox *ScrollableBox) addRow(row tRow) {
	scrollableBox.rows = append(scrollableBox.rows, row)
}

func (scrollableBox *ScrollableBox) clearRows() {
	scrollableBox.rows = []tRow{}
}

func (scrollableBox *ScrollableBox) getVisibleRows() []tRow {
	_, _, _, height := scrollableBox.Box.GetInnerRect()

	var rowsHeight int
	endIndex := scrollableBox.startIndex
	for _, row := range scrollableBox.rows[scrollableBox.startIndex:] {
		rowsHeight += row.getHeight()
		if rowsHeight <= height {
			endIndex++

		} else {
			break
		}
	}
	return scrollableBox.rows[scrollableBox.startIndex:endIndex]
}

func (scrollableBox *ScrollableBox) getSelectedRow() tRow {
	return scrollableBox.rows[scrollableBox.selectedIndex]
}

func (scrollableBox *ScrollableBox) setSelectedRow(index int) {
	scrollableBox.selectedIndex = index
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

func (scrollableBox *ScrollableBox) onMouseLeftClick(event *tcell.EventMouse) {
	visibleRows := scrollableBox.getVisibleRows()
	_, yClick := event.Position()
	yPos := 1
	for index, row := range visibleRows {
		rowHeight := row.getHeight()
		if yClick >= yPos && yClick < yPos+rowHeight {
			scrollableBox.setSelectedRow(index)
			break
		}
		yPos += row.getHeight()

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

func (scrollableBox *ScrollableBox) SetTitle(title string) *ScrollableBox {
	scrollableBox.Box.SetTitle(title)
	return scrollableBox
}

func (scrollableBox *ScrollableBox) SetFocus() tview.Primitive {
	scrollableBox.Box.Focus(nil)
	return scrollableBox.parent
}

func (scrollableBox *ScrollableBox) SetBorder(show bool) tview.Primitive {
	scrollableBox.Box.SetBorder(show)
	return scrollableBox.parent
}
