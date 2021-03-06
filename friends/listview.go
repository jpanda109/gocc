package friends

import "github.com/jpanda109/gocc/config"

import "github.com/nsf/termbox-go"

// ListView defines view where the user can view saved friends along
// with their basic information
// w is the allocated width of the screen
// h is the allocated height of the screen
// friends is a list of friends objects to be displayed
// curline is the current line of the cursor
// topline is the index of the top friend displayed on screen
type ListView struct {
	w, h    int
	friends []*config.Friend
	curline int
	topline int
}

// NewListView is a future-proofed ListView initializer
func NewListView() *ListView {
	view := &ListView{
		w:       0,
		h:       0,
		friends: []*config.Friend{},
		curline: 0,
		topline: 0,
	}
	return view
}

// Start initializes and runs this screen
func (view *ListView) Start(w, h int, friends []*config.Friend) {
	view.w, view.h, view.friends = w, h, friends
	view.Refresh()
}

// Refresh refreshes the screen based on curline and topline
func (view *ListView) Refresh() {
	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)
	for y, friend := range view.friends {
		for x, c := range friend.Name {
			bgColor := termbox.ColorBlack
			fgColor := termbox.ColorWhite
			if view.curline == y {
				bgColor = termbox.ColorMagenta
			}
			termbox.SetCell(x, y, rune(c), fgColor, bgColor)
		}
	}
	termbox.Flush()
}

// Direction is a typedef for how to move the cursor
type Direction int

const (
	// Up moves cursor up
	Up Direction = iota
	// Down moves cursor down
	Down
)

// MoveCursor moves the cursor in the direction specified if possible
func (view *ListView) MoveCursor(dir Direction) {
	switch dir {
	case Up:
		if view.curline > 0 {
			view.curline--
		}
	case Down:
		if view.curline < view.h-1 {
			view.curline++
		}
	}
	view.Refresh()
}

// Resize resets the ListView object's w and h attributes to match the given w
// and w
func (view *ListView) Resize(w, h int) {
	view.w, view.h = w, h
}
