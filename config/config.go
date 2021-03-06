package config

import (
	"encoding/json"
	"os"
)

// Friend contains info about a friend
// Name is a nickname for the friend
// Addrs is the ip + port at which the friend is located
type Friend struct {
	Name        string
	Description string
	Addr        string
	Ports       []string
}

// Config holds the configuration data from config file
type Config struct {
	Friends []*Friend
	Port    string
	Name    string
}

var config *Config
var save func() error

// Init instantiates the configuration object
// Should be called at the beginning of program and only once.
func Init() error {
	config = &Config{
		Friends: make([]*Friend, 0),
		Port:    "8080",
		Name:    "anon",
	}
	save = saveConfig
	return readConfig()
}

// SetDebug sets the config package to debug mode
func SetDebug(debug bool) {
	if debug {
		save = func() error {
			return nil
		}
	} else {
		save = saveConfig
	}
}

// Name returns the default name
func Name() string {
	return config.Name
}

// SetName persists the name and overwrites the current config file
func SetName(name string) {
	config.Name = name
	save()
}

// Port returns the default port
func Port() string {
	return config.Port
}

// SetPort persists the port to the config file
func SetPort(port string) {
	config.Port = port
	save()
}

// Friends returns your list of friends
func Friends() []*Friend {
	return config.Friends
}

// AddFriend persists the friend to config file
func AddFriend(friend *Friend) {
	config.Friends = append(config.Friends, friend)
	save()
}

func saveConfig() error {
	file, err := os.Create("config.json")
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	err = encoder.Encode(config)
	return nil
}

func readConfig() error {
	file, err := os.Open("config.json")
	if err != nil {
		saveConfig()
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return err
	}
	return nil
}
