package main

import (
	"fmt"
	"github.com/fatih/color"
	"time"
)

var newPlayer *Player

var yellow = color.New(color.FgYellow).SprintFunc()
var red = color.New(color.FgRed).SprintFunc()
var cyan = color.New(color.FgCyan).SprintFunc()

var tSecond = time.NewTicker(time.Second)

type Living struct {
	name         string
	age          int
	gender       string
	gold_on_hand int
	health       int
	race         string
	x            int
	y            int
	z            int
}

type Player struct {
	Living
	gold_in_bank int
	class        string
	right_hand   string
	left_hand    string
	possessive   string
	objective    string
}

type NPC struct {
	Living
}

func new_game() {
	var obj string
	var pos string
	// Create new SprintXxx functions for later
	fmt.Print("What is your name? ")
	playerName := ""
	fmt.Scanln(&playerName)
	fmt.Print("Gender? ")
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
	newPlayer = new(Player)
	newPlayer.name = playerName
	newPlayer.age = 25
	newPlayer.gender = playerGender
	newPlayer.gold_on_hand = 0
	newPlayer.gold_in_bank = 0
	newPlayer.health = 100
	newPlayer.race = "human"
	newPlayer.class = "adventurer"
	newPlayer.right_hand = "nothing"
	newPlayer.left_hand = "nothing"
	newPlayer.possessive = pos
	newPlayer.objective = obj
	newPlayer.x = 0
	newPlayer.y = 0
	newPlayer.z = 0
	fmt.Printf("%s the %s %s has %s.\n", yellow(newPlayer.name), cyan(newPlayer.race), cyan(newPlayer.class), red("entered the world"))
}

func tick() {
	fmt.Println("tick")
}

func main() {
	new_game()
	fmt.Println(newPlayer.name, "begins", newPlayer.possessive, "adventure.")
	// Initiate game loop
	for {
		tick()
		<-tSecond.C
	}
}
