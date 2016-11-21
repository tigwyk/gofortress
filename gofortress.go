package main

import (
	"fmt"
	"github.com/fatih/color"
	"math"
	"math/rand"
	"os"
	"time"
)

var newPlayer *Player

var yellow = color.New(color.FgYellow).SprintFunc()
var red = color.New(color.FgRed).SprintFunc()
var cyan = color.New(color.FgCyan).SprintFunc()
var green = color.New(color.FgGreen).SprintFunc()

var tSecond = time.NewTicker(time.Second * 2)

const XP_CONSTANT = 0.04

var allNPCs []*NPC
var Rooms []*Room
var Weapons []*Weapon

type Room struct {
	x, y, z int
}

type Living struct {
	name         string
	age          int
	gender       string
	gold_on_hand int
	max_health   int
	health       int
	race         string
	x            int
	y            int
	z            int
	actions      []string
	hunger       int
	thirst       int
	desire       string
	killer       string
	xp           int
	level        int
	wielded      *Weapon
}

type Player struct {
	Living
	gold_in_bank int
	class        string
	possessive   string
	objective    string
}

type NPC struct {
	Living
	bonus_xp int
}

type Item struct {
	name string
}

type Weapon struct {
	Item
	verb string
}

func (p *Player) Coords() string {
	return fmt.Sprintf("%d,%d,%d", p.x, p.y, p.z)
}

func (p *Player) Damage(d int, e Entity) {
	if p.health > 0 {
		p.health = p.health - d
	}
	if p.health < 0 {
		p.health = 0
	}
	fmt.Println(p.name, "takes", d, "damage from", e.Name()+"!", p.health, "/", p.max_health)
	if p.health < 1 {
		p.Die(e)
	}
}

func (n *NPC) Coords() string {
	return fmt.Sprintf("%d,%d,%d", n.x, n.y, n.z)
}

func (n *NPC) Damage(d int, e Entity) {
	if n.health > 0 {
		n.health = n.health - d
	}
	if n.health < 0 {
		n.health = 0
	}
	fmt.Println(n.name, "takes", d, "damage from", e.Name()+"!", n.health, "/", n.max_health)
	if n.health < 1 {
		n.Die(e)
	}
}

func (n *NPC) HP() int {
	return n.health
}

func (n *NPC) Name() string {
	return n.name
}

func (p *Player) HP() int {
	return p.health
}

func (p *Player) Name() string {
	return p.name
}

func (i *Item) Name() string {
	return i.name
}

func (r *Room) Coords() string {
	return fmt.Sprintf("%d,%d,%d", r.x, r.y, r.z)
}

func (n *NPC) XP() int {
	return (n.xp + n.bonus_xp)
}

func (p *Player) XP() int {
	return p.xp
}

func (p *Player) XPToNextLevel() int {
	fLevel := float64(p.level)
	result := math.Pow(fLevel/XP_CONSTANT, 2)
	return int(result)
}

func (n *NPC) CalcXP() {
	n.xp = int(math.Pow(float64(n.level)/XP_CONSTANT, 2) / 180)
}

func (p *Player) GainXP(x int) {
	var xpLeft int
	xpLeft = p.XPToNextLevel() - p.xp
	if xpLeft < 0 {
		xpLeft = 0
	}
	fmt.Printf(green("%s has gained %dXP and needs %d more to level up!\n"), p.name, x, xpLeft)
	p.xp = p.xp + x
	fmt.Println(p.XP(), "/", p.XPToNextLevel())
}

func (n *NPC) GainXP(x int) {
	n.xp = n.xp + x
}

func (p *Player) Wield(w *Weapon) bool {
	fmt.Println(p.name, "attempts to wield a", w.name+".")
	switch {
	case p.wielded.name == "fists":
		p.wielded = w
		fmt.Println(p.name, "wields a", w.name)
		return true
	default:
		fmt.Println(p.name, "cannot wield a", w.name, "without putting something down.")
		return false
	}
}

type Entity interface {
	HP() int
	Damage(d int, e Entity)
	Name() string
	XP() int
	GainXP(x int)
}

type generic_thing interface {
	Name() string
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
	newPlayer.max_health = 100
	newPlayer.race = "human"
	newPlayer.class = "mercenary"
	newPlayer.level = 1
	newPlayer.possessive = pos
	newPlayer.objective = obj
	newPlayer.wielded = createWeapon("fists", "punch")
	newPlayer.x = 0
	newPlayer.y = 0
	newPlayer.z = 0
	newPlayer.xp = 0
	newPlayer.hunger = 100
	newPlayer.thirst = 0
	fmt.Printf("%s the %s %s has %s.\n", yellow(newPlayer.name), cyan(newPlayer.race), cyan(newPlayer.class), red("entered the world"))
	var newRoom = new(Room)
	fmt.Println(newRoom)
	createNewWeapons()
}

func ProcessPlayerActions() {
	fmt.Println("Player action.")
	newPlayer.hunger = newPlayer.hunger + 1
	newPlayer.thirst = newPlayer.thirst + 1
	newPlayer.processDesires()
	newPlayer.processEnvironment()
	newPlayer.processStats()
}

func (p *Player) processDesires() {
	if p.hunger > 70 {
		p.desire = "food"
	}
	if p.thirst > 40 {
		p.desire = "drink"
	}
	if p.desire != "" {
		fmt.Println("Current desire:", p.desire)
	}
}

func (p *Player) processEnvironment() {
	for v := range allNPCs {
		if p.Coords() == allNPCs[v].Coords() {
			fmt.Println(p.name, "has noticed", allNPCs[v].name+".")
			p.Attack(allNPCs[v], 5)
		}
	}
}

func (p *Player) processStats() {
	var xtnl int
	xtnl = p.XPToNextLevel()
	if xtnl < 0 {
		xtnl = -xtnl
	}
	if p.xp >= xtnl {
		p.xp = p.xp - xtnl
		p.level = p.level + 1
		fmt.Printf(cyan("%s has achieved level %d!\n"), p.name, p.level)
		xtnl = 0
	}
}

func ProcessNPCActions() {
	fmt.Println("NPC action.")
	if len(allNPCs) == 0 {
		var newNPC = new(NPC)
		newNPC.name = "Bob"
		newNPC.age = 40
		newNPC.gender = "male"
		newNPC.gold_on_hand = 10
		newNPC.health = 10
		newNPC.max_health = 10
		newNPC.race = "human"
		newNPC.level = rand.Intn(10)
		newNPC.CalcXP()
		allNPCs = append(allNPCs, newNPC)
	}
}

func tick() {
	ProcessPlayerActions()
	ProcessNPCActions()
}

func (p *Player) Attack(e Entity, d int) {
	fmt.Println(p.name, p.wielded.verb+"s", e.Name(), "for", d, "damage with", p.possessive, p.wielded.name+"!")
	e.Damage(d, p)
}

func (n *NPC) Attack(e Entity, d int) {
	e.Damage(d, n)
}

func (p *Player) Die(e Entity) {
	switch {
	case e != nil:
		fmt.Printf(red(e.Name() + " has killed " + p.name + "!\n"))
	case e == nil:
		fmt.Printf(red(p.name + " has died.\n"))
	default:
		fmt.Printf(red(p.name + " has died.\n"))
	}
	GameOver()
}

func (n *NPC) Die(e Entity) {
	switch {
	case e != nil:
		fmt.Printf(red(e.Name() + " has killed " + n.name + "!\n"))
		e.GainXP(n.XP())
	case e == nil:
		fmt.Printf(red(n.name + " has died.\n"))
	default:
		fmt.Printf(red(n.name + " has died.\n"))
	}
	allNPCs = append(allNPCs[:0], allNPCs[1:]...)
}

func createNewWeapons() {
	weapon_names := []string{"handgun", "shotgun", "rifle", "assault rifle", "submachine gun"}
	weapon_verbs := []string{"shoot", "blast"}
	for i := 0; i < 5; i++ {
		w := createWeapon(weapon_names[rand.Intn(len(weapon_names))], weapon_verbs[rand.Intn(len(weapon_verbs))])
		//w.name = weapon_names[rand.Intn(len(weapon_names))]
		Weapons = append(Weapons, w)
	}
	fmt.Println("Generated new weapons:", Weapons)
	newPlayer.Wield(Weapons[rand.Intn(len(Weapons))])
}

func createWeapon(name string, verb string) *Weapon {
	w := new(Weapon)
	w.name = name
	w.verb = verb
	return w
}

func GameOver() {
	fmt.Printf(red("Game Over\n"))
	os.Exit(0)
}

func main() {
	rand.Seed(1337)
	new_game()
	fmt.Println(newPlayer.Name(), "begins", newPlayer.possessive, "adventure.")
	// Initiate game loop
	for {
		tick()
		<-tSecond.C
	}
}
