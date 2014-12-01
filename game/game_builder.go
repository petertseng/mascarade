package game

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/petertseng/mascarade/format"
	"github.com/petertseng/mascarade/output"
	"github.com/petertseng/mascarade/player"
	"github.com/petertseng/mascarade/role"
)

type GameBuilder struct {
	roles       map[role.Role]bool
	playerNames []string
}

func NewBuilder() GameBuilder {
	roles := make(map[role.Role]bool)
	playerNames := make([]string, 0)

	return GameBuilder{roles: roles, playerNames: playerNames}
}

func (gb *GameBuilder) AddPlayer(name string) error {
	// TODO: What if player is already in the game
	gb.playerNames = append(gb.playerNames, name)
	return nil
}

func (gb *GameBuilder) AddRole(name string) error {
	role, err := role.FromString(name)
	if err != nil {
		return err
	}
	gb.roles[role] = true
	return nil
}

func (gb *GameBuilder) MakeGame(out io.Writer) (Game, error) {
	// Make the roles array
	roles := make([]role.Role, 0)
	rolesPresent := make(map[role.Role]bool)
	for role, included := range gb.roles {
		if included {
			roles = append(roles, role)
			rolesPresent[role] = true
			if role.Pair() {
				roles = append(roles, role)
			}
		}
	}

	if len(roles) < len(gb.playerNames) {
		err := fmt.Errorf("Not enough roles (%d) for the players (%d).", len(roles), len(gb.playerNames))
		return Game{}, err
	}

	rand.Seed(time.Now().UTC().UnixNano())
	playerPerm := rand.Perm(len(gb.playerNames))
	rolePerm := rand.Perm(len(roles))

	playerOrder := make([]*player.Player, len(gb.playerNames))
	playerMap := make(map[string]*player.Player)

	for i, name := range gb.playerNames {
		seatingOrder := playerPerm[i]
		role := roles[rolePerm[i]]
		p := player.New(name, role)
		playerOrder[seatingOrder] = &p
		playerMap[name] = &p
	}

	// table cards
	numTableCards := len(roles) - len(gb.playerNames)
	tableCards := make([]*player.TableCard, numTableCards)
	for i, _ := range tableCards {
		tc := player.NewTableCard(i, roles[rolePerm[i+len(gb.playerNames)]])
		tableCards[i] = &tc
	}

	g := Game{
		roles:       rolesPresent,
		players:     playerMap,
		playerOrder: playerOrder,
		tableCards:  tableCards,
		input:       bufio.NewReader(os.Stdin),
		format:      format.NewText(output.NewPrefixed(out)),
	}
	g.startGame()
	return g, nil
}
