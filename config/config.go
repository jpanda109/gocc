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
	Addrs       []string
}

// Config holds the configuration data from config file
type Config struct {
	Friends []*Friend
	Ports   []string
	Name    string
}

var config *Config

// Init instantiates the configuration object
// Should be called at the beginning of program and only once.
func Init() error {
	config = &Config{
		Friends: make([]*Friend, 0),
		Ports:   []string{"8080"},
		Name:    "anon",
	}
	return readConfig()
}

// GetPort returns the default port
func GetPort() []string {
	return config.Ports
}

// AddPort persists the port to the config file
func AddPort(port string) {
	config.Ports = append(config.Ports, port)
	saveConfig()
}

// GetFriends returns your list of friends
func GetFriends() []*Friend {
	return config.Friends
}

// AddFriend persists the friend to config file
func AddFriend(friend *Friend) {
	config.Friends = append(config.Friends, friend)
	saveConfig()
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
