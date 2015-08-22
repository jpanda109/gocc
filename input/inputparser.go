package input

import "strings"

// command is a type alias which defines different commands
type commandType int

const (
	save commandType = iota
)

// parseInput return a command type and its args if the input string is a valid
// command (e.g. /save 1) but returns nil if there's no such command. The last
// return value is a boolean which is true if it was a valid command
func parseInput(input string) (commandType, []string, bool) {
	if len(input) == 0 {
		return -1, nil, false
	}
	tokens := strings.Split(input, " ")
	if tokens[0][0] == '/' {
		if tokens[0] == "/save" {
			return save, tokens[1:], true
		}
	}
	return -1, nil, false
}
