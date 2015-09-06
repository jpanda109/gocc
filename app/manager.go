package app

import (
	"errors"
	"sync"
)

// Screen defines an interface which has start and stop functions
type Screen interface {
	Start()
	Stop()
	SetManager(manager *Manager)
}

// Manager manages screens and ensures only one is being used at a time
type Manager struct {
	screen Screen
	wg     *sync.WaitGroup
}

// Start activates the given screen and returns a waitgroup which blocks
// until Quit() called
func (m *Manager) Start(screen Screen) (*sync.WaitGroup, error) {
	if m.screen != nil { // already started
		return nil, errors.New("")
	}
	m.screen = screen
	m.screen.Start()
	var wg sync.WaitGroup
	m.wg = &wg
	m.wg.Add(1)
	return m.wg, nil
}

// SwitchScreeen closes the previous screen and opens a new one
func (m *Manager) SwitchScreeen(screen Screen) {
	m.screen.Stop()
	m.screen = screen
	go m.screen.Start()
}

// Quit closes the current screen and quits the application
func (m *Manager) Quit() {
	m.screen.Stop()
	m.wg.Done()
}
