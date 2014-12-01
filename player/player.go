package player

import (
	"fmt"

	"github.com/petertseng/mascarade/role"
)

type roleOwner struct {
	role role.Role
}

type Player struct {
	name  string
	coins uint64

	roleOwner
	lastRevealedTurn uint
}

type TableCard struct {
	id int
	roleOwner
}

func NewTableCard(id int, role role.Role) TableCard {
	return TableCard{id: id, roleOwner: roleOwner{role: role}}
}
func (tc TableCard) Name() string {
	return fmt.Sprintf("Table Card %d", tc.id)
}

func New(name string, role role.Role) Player {
	return Player{name: name, coins: 6, roleOwner: roleOwner{role: role}}
}

type CoinOwner interface {
	Name() string
	Coins() uint64
	AddCoins(uint64)
	PayFine()
	Pay(CoinOwner, uint64) uint64
}

type Swappable interface {
	Name() string
	Role() role.Role
	setRole(r role.Role)
	SwapRoles(swapWith Swappable)
}

func (p Player) Name() string {
	return p.name
}

func (p Player) Coins() uint64 {
	return p.coins
}

func (p *Player) Pay(payTo CoinOwner, coins uint64) (coinsTaken uint64) {
	if p.coins < coins {
		coinsTaken = p.coins
		payTo.AddCoins(coinsTaken)
		p.coins = 0
		return
	}
	coinsTaken = coins
	p.coins -= coinsTaken
	payTo.AddCoins(coinsTaken)
	return
}

func (p *Player) AddCoins(coins uint64) {
	p.coins += coins
}

func (p *Player) PayFine() {
	p.coins--
}

func (p *Player) Reveal(turn uint) {
	p.lastRevealedTurn = turn
}
func (p Player) LastRevealed() uint {
	return p.lastRevealedTurn
}

func (r roleOwner) Role() role.Role {
	return r.role
}
func (r *roleOwner) setRole(role role.Role) {
	r.role = role
}

func (r *roleOwner) SwapRoles(swapWith Swappable) {
	otherRole := swapWith.Role()
	swapWith.setRole(r.Role())
	r.setRole(otherRole)
}
