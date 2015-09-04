package friends

import "github.com/jpanda109/gocc/config"

// DetailView defines a view where the user can see detailed information
// about a specific friend, including the description and notes the user
// may  have saved
type DetailView struct {
	w, h   int
	friend *config.Friend
}

// NewDetailView is a future proof DetailView initializer
func NewDetailView(friend *config.Friend) *DetailView {
	view := &DetailView{
		w:      0,
		h:      0,
		friend: friend,
	}
	return view
}

// Start initializes and runs the screen
func (view *DetailView) Start(w, h int) {
	view.w, view.h = w, h
}

// Resize initializes and starts this screen
func (view *DetailView) Resize(w, h int) {
	view.w, view.h = w, h
}
