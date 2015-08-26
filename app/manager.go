package app

// Screen defines an interface which has start and stop functions
type Screen interface {
	Start()
	Stop()
}

// Manager manages screens and ensures only one is being used at a time
type Manager struct {
	screen *Screen
}
