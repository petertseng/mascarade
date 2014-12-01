package role

import (
	"fmt"
	"strings"
)

type Role int

const (
	NoSuchRole Role = iota

	Judge
	Bishop
	King
	Fool
	Queen
	Thief
	Witch
	Spy
	Peasant
	Cheat
	Inquisitor
	Widow

	Alchemist
	Actress
	Courtesan
	Gambler
	PuppetMaster
	Damned
	Patron
	Beggar
	Necromancer
	Princess
	Sage
	Usurper
)

var ids = map[string]Role{}
var idsInitialized = false

func FromString(s string) (Role, error) {
	if !idsInitialized {
		for id, nameAndPower := range namesAndPowers {
			ids[strings.ToLower(nameAndPower[0])] = id
		}
		idsInitialized = true
	}

	role, ok := ids[strings.ToLower(s)]
	if ok {
		return role, nil
	}

	return NoSuchRole, fmt.Errorf("No such role %s", s)
}

var namesAndPowers = map[Role][2]string{
	Judge:        {"Judge", "Take all of the courthouse's gold"},
	Bishop:       {"Bishop", "Take 2 coins from the richest of the other players"},
	King:         {"King", "Take 3 coins"},
	Fool:         {"Fool", "Take 1 coin, then swap (or not) two cards not your own"},
	Queen:        {"Queen", "Take 2 coins"},
	Thief:        {"Thief", "Take one coin from each adjacent player"},
	Witch:        {"Witch", "May swap fortune with another player"},
	Spy:          {"Spy", "Look at own card and another, then swap (or not)"},
	Peasant:      {"Peasant", "Take 1 coin, or 2 coins if both Peasant reveal"},
	Cheat:        {"Cheat", "Wins with 10 coins"},
	Inquisitor:   {"Inquisitor", "Target must guess own character or pay 4 coins"},
	Widow:        {"Widow", "Take coins from the bank until at 10 coins"},
	Alchemist:    {"Alchemist", "Everyone passes coins left or right"},
	Actress:      {"Actress", "Use power of previous character played"},
	Courtesan:    {"Courtesan", "Previous player reveals and pays 3 coins if male"},
	Gambler:      {"Gambler", "Play guessing game with another play to win 1, 2, or 3 coins"},
	PuppetMaster: {"Puppet Master", "Switch places of two players and take 1 coin from each"},
	Damned:       {"Damned", "You are eliminated!!!"},
	Patron:       {"Patron", "Take 3 coins, and your neighbors take 1 coin each"},
	Beggar:       {"Beggar", "Take 1 coin from every player richer"},
	Necromancer:  {"Necromancer", "Reveal a character from the Cemetery and use its power"},
	Princess:     {"Princess", "Take 2 coins and choose a player to show character to others"},
	Sage:         {"Sage", "Take 1 coin and look at two players' cards"},
	Usurper:      {"Usurper", "Guess another player's character and use that power if successful"},
}

func (r Role) String() string {
	nameAndPower, ok := namesAndPowers[r]
	if ok {
		return nameAndPower[0]
	}

	return fmt.Sprintf("Unknown role ID %d", r)
}

func (r Role) PowerDescription() string {
	nameAndPower, ok := namesAndPowers[r]
	if ok {
		return nameAndPower[1]
	}

	return fmt.Sprintf("Power for Unknown role ID %d", r)
}

func (r Role) Pair() bool {
	return r == Peasant
}

func (r Role) CanAnnounce() bool {
	return r != Damned
}
