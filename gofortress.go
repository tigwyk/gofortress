package main

import (
	"fmt"
	"github.com/fatih/color"
)

var newPlayer player
var yellow = color.New(color.FgYellow).SprintFunc()
var red = color.New(color.FgRed).SprintFunc()
var cyan = color.New(color.FgCyan).SprintFunc()

type player struct {
	name         string
	age          int
	gender       string
	gold_on_hand int
	gold_in_bank int
	health       int
	race         string
	class        string
	right_hand   string
	left_hand    string
	possessive   string
	objective    string
	x            int
	y            int
	z            int
}

func new_game() {
	var obj string
	var pos string
	// Create new SprintXxx functions for later
	fmt.Print("What is your name? ")
	playerName := ""
	fmt.Scanln(&playerName)
	fmt.Print("What is your gender? ")
	playerGender := ""
	fmt.Scanln(&playerGender)
	switch playerGender {
	case "male":
		obj = "him"
		pos = "his"
	case "female":
		obj = "her"
		pos = "her"
	default:
		obj = "them"
		pos = "their"
	}
	newPlayer = player{playerName, 25, playerGender, 0, 0, 100, "human", "adventurer", "nothing", "nothing", pos, obj, 0, 0, 0}
	fmt.Printf("%s the %s %s has %s.\n", yellow(newPlayer.name), cyan(newPlayer.race), cyan(newPlayer.class), red("entered the world"))
}

func main() {
	new_game()
	fmt.Println(newPlayer.name, "begins", newPlayer.possessive, "adventure.")
	// Initiate game loop

}
