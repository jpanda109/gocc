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

// SetFriends is simply a publicly accessible mutator method for friends
func (view *ListView) SetFriends(friends []*config.Friend) {
	view.friends = friends
}

// Start initializes and runs this screen
func (view *ListView) Start(w, h int) {
	view.w, view.h = w, h
	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)
	for y, friend := range view.friends {
		for x, c := range friend.Name {
			termbox.SetCell(x, y, rune(c), termbox.ColorWhite, termbox.ColorBlack)
		}
	}
	termbox.Flush()
}

// Resize resets the ListView object's w and h attributes to match the given w
// and w
func (view *ListView) Resize(w, h int) {
	view.w, view.h = w, h
}
