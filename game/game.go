package game

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/petertseng/mascarade/format"
	"github.com/petertseng/mascarade/player"
	"github.com/petertseng/mascarade/role"
)

type GameResolver interface {
	UserChoice(user string) []string
	PlayersOtherThan(string) map[string]*player.Player
	CoinOwners() map[string]player.CoinOwner
	CoinOwnersNextTo(string) (player.CoinOwner, player.CoinOwner, error)
	SwappablesOtherThan(string) map[string]player.Swappable
	RichestOtherThan(string) map[string]player.CoinOwner
	TakeCourthouse() uint64
	RevealCard(*player.Player)
	CheaterWins(*player.Player)
}

type Game struct {
	roles              map[role.Role]bool
	players            map[string]*player.Player
	playerOrder        []*player.Player
	tableCards         []*player.TableCard
	currentPlayerIndex int

	deadPlayers []*player.Player

	input  *bufio.Reader
	format format.Formatter

	turnCount  uint
	courthouse uint64

	claim            bool
	claimPlayerIndex int
	claimedRole      role.Role
	otherClaimants   []*player.Player

	cheatWinner *player.Player

	winners []string

	// TODO cemetery
}

func (g *Game) startGame() {
	g.format.YourTurn(g.ActivePlayerName())
}

func (g *Game) advanceActivePlayer() {
	g.currentPlayerIndex++
	if g.currentPlayerIndex == len(g.playerOrder) {
		g.currentPlayerIndex = 0
	}
}

func (g *Game) advanceTurn() {
	winner := g.checkVictory()

	if winner {
		return
	}

	g.claim = false
	g.otherClaimants = []*player.Player{}
	g.turnCount++
	g.advanceActivePlayer()
	g.format.YourTurn(g.ActivePlayerName())
}

func (g *Game) advanceClaim() {
	g.advanceActivePlayer()
	if g.currentPlayerIndex == g.claimPlayerIndex {
		g.resolveClaim()
	} else {
		g.format.YourTurnToChallenge(g.ActivePlayerName(), g.AnnouncingPlayerName(), g.claimedRole)
	}
}

func (g *Game) resolveClaim() {
	correct := make([]*player.Player, 0)
	liars := make([]player.CoinOwner, 0)

	if len(g.otherClaimants) == 0 {
		g.format.NobodyChallenged(g.ActivePlayerName(), g.claimedRole)
		correct = append(correct, g.activePlayer())
	} else {
		allClaimants := append(g.otherClaimants, g.activePlayer())
		for _, claimant := range allClaimants {
			role := claimant.Role()
			claimant.Reveal(g.turnCount)
			if role == g.claimedRole {
				correct = append(correct, claimant)
				g.format.GoodClaim(claimant.Name(), g.claimedRole)
			} else {
				liars = append(liars, claimant)
				g.format.BadClaim(claimant.Name(), claimant.Role(), g.claimedRole)
			}
		}
	}

	for _, p := range correct {
		g.usePower(p, g.claimedRole, len(correct))
	}

	// TODO: If using a wildcard power that copied a pair power, a special case here.

	for _, liar := range liars {
		liar.PayFine()
		g.courthouse++
		g.format.PayFine(liar.Name(), liar.Coins())
	}
	if len(liars) > 0 {
		g.format.Courthouse(g.courthouse)
	}

	g.advanceTurn()
}

func (g *Game) usePower(p *player.Player, r role.Role, numCorrect int) {
	g.format.UsePower(p.Name(), r)

	power, ok := powers[r]
	if !ok {
		panic(fmt.Sprintf("No power found for role %s", r))
	}

	power(g, p, numCorrect, g.format)
}

func (g *Game) checkVictory() bool {
	// If anyone is cheater
	if g.cheatWinner != nil {
		g.format.CheaterWins(g.cheatWinner.Name())
		g.winners = []string{g.cheatWinner.Name()}
		return true
	}

	// Figure out whether anyone is broke, and who is the richest, in one pass
	zeroCoins := make([]string, 0)
	highestCoins := uint64(0)
	for _, player := range g.players {
		if player.Coins() == 0 {
			zeroCoins = append(zeroCoins, player.Name())
		}
		if player.Coins() > highestCoins {
			highestCoins = player.Coins()
		}
	}

	// Someone's at 13, all such players win
	if highestCoins >= 13 {
		winners := make([]string, 0)
		for _, player := range g.players {
			if player.Coins() >= 13 {
				winners = append(winners, player.Name())
			}
		}
		g.format.WinTargetReached(winners)
		g.winners = winners
		return true
	}

	// Someone's broke, the richest players win
	if len(zeroCoins) > 0 {
		richest := make([]string, 0)
		for _, player := range g.players {
			if player.Coins() == highestCoins {
				richest = append(richest, player.Name())
			}
		}
		g.format.WinBroke(richest, zeroCoins)
		g.winners = richest
		return true
	}

	return false
}

func (g *Game) activePlayer() *player.Player {
	return g.playerOrder[g.currentPlayerIndex]
}
func (g *Game) announcingPlayer() *player.Player {
	return g.playerOrder[g.claimPlayerIndex]
}

func (g *Game) ActivePlayerName() string {
	return g.activePlayer().Name()
}
func (g *Game) AnnouncingPlayerName() string {
	return g.announcingPlayer().Name()
}

func (g *Game) Peek() error {
	if g.claim {
		return fmt.Errorf("You must respond to %s's claim of %s", g.AnnouncingPlayerName(), g.claimedRole)
	}
	if g.turnCount < 4 {
		return fmt.Errorf("For the first four turns, you must swap (or not)")
	}
	if g.activePlayer().LastRevealed() == g.turnCount-1 {
		return fmt.Errorf("Because you revealed your card on the previous turn, you must swap (or not)")
	}

	g.format.Peek(g.ActivePlayerName())
	g.format.TellOwnCard(g.ActivePlayerName(), g.activePlayer().Role())

	g.advanceTurn()

	return nil
}

func (g *Game) SwapOrNot(target string, actuallySwap bool) error {
	if g.claim {
		return fmt.Errorf("You must respond to %s's claim of %s", g.AnnouncingPlayerName(), g.claimedRole)
	}

	swappable, err := g.ResolveSwappable(target)
	if err != nil {
		return err
	}

	if swappable == g.activePlayer() {
		return fmt.Errorf("You can't swap with yourself, %s", g.ActivePlayerName())
	}

	if actuallySwap {
		g.activePlayer().SwapRoles(swappable)
	}

	g.format.SwapOrNot(g.ActivePlayerName(), swappable.Name())

	g.advanceTurn()

	return nil
}

func (g *Game) ClaimRole(roleName string) error {
	if g.claim {
		return fmt.Errorf("You must respond to %s's claim of %s", g.AnnouncingPlayerName(), g.claimedRole)
	}
	if g.turnCount < 4 {
		return fmt.Errorf("For the first four turns, you must swap (or not)")
	}
	if g.activePlayer().LastRevealed() == g.turnCount-1 {
		return fmt.Errorf("Because you revealed your card on the previous turn, you must swap (or not)")
	}

	roleClaimed, err := role.FromString(roleName)
	if err != nil {
		return err
	}

	if !roleClaimed.CanAnnounce() {
		return fmt.Errorf("The %s can never be announced", roleClaimed)
	}
	if !g.roles[roleClaimed] {
		return fmt.Errorf("The %s is not in this game", roleClaimed)
	}

	g.claim = true
	g.claimPlayerIndex = g.currentPlayerIndex
	g.claimedRole = roleClaimed
	g.format.ClaimRole(g.ActivePlayerName(), g.claimedRole)

	g.advanceClaim()

	return nil
}

func (g *Game) NoChallenge() error {
	if !g.claim {
		return fmt.Errorf("No role has been announced")
	}

	g.format.NoCounterclaim(g.ActivePlayerName(), g.AnnouncingPlayerName(), g.claimedRole)

	g.advanceClaim()

	return nil
}

func (g *Game) Challenge() error {
	if !g.claim {
		return fmt.Errorf("No role has been announced")
	}

	g.format.Counterclaim(g.ActivePlayerName(), g.AnnouncingPlayerName(), g.claimedRole)
	g.otherClaimants = append(g.otherClaimants, g.activePlayer())

	g.advanceClaim()

	return nil
}

func (g *Game) resolvePlayer(name string) (*player.Player, error) {
	player, ok := g.players[name]
	if !ok {
		return nil, fmt.Errorf("No such player %s", name)
	}
	return player, nil
}

func (g *Game) ResolveSwappable(name string) (player.Swappable, error) {
	// If it's a table card...
	if name[0] == '#' {
		index, err := strconv.ParseInt(name[1:], 0, 64)
		if err != nil {
			return nil, fmt.Errorf("No such player %s: %s", name, err)
		}
		if index >= int64(len(g.tableCards)) {
			return nil, fmt.Errorf("There are only %d table cards, so #%d is invalid", len(g.tableCards), index)
		}
		return g.tableCards[index], nil
	}

	// Otherwise, it's a player
	return g.resolvePlayer(name)
}

func (g *Game) UserChoice(user string) []string {
	for {
		str, err := g.input.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}
		return strings.Fields(str)
	}
}

func (g *Game) PlayersOtherThan(exclude string) map[string]*player.Player {
	m := make(map[string]*player.Player)

	for name, player := range g.players {
		if name != exclude {
			m[name] = player
		}
	}

	return m
}

func (g *Game) SwappablesOtherThan(exclude string) map[string]player.Swappable {
	m := make(map[string]player.Swappable)

	for name, player := range g.players {
		if name != exclude {
			m[name] = player
		}
	}

	for i, table := range g.tableCards {
		name := fmt.Sprintf("#%d", i)
		m[name] = table
	}

	return m
}

func (g *Game) CoinOwners() map[string]player.CoinOwner {
	m := make(map[string]player.CoinOwner)

	for name, player := range g.players {
		m[name] = player
	}

	return m
}

func (g *Game) CoinOwnersNextTo(nextto string) (before, after player.CoinOwner, err error) {
	index := 0
	found := false

	for i, player := range g.playerOrder {
		if player.Name() == nextto {
			index = i
			found = true
		}
	}

	if !found {
		return nil, nil, fmt.Errorf("No such player %s", nextto)
	}

	beforeIndex := index - 1
	if beforeIndex < 0 {
		beforeIndex = len(g.playerOrder) - 1
	}
	before = g.playerOrder[beforeIndex]

	afterIndex := index + 1
	if afterIndex >= len(g.playerOrder) {
		afterIndex = 0
	}
	after = g.playerOrder[afterIndex]
	err = nil
	return
}

func (g *Game) RichestOtherThan(exclude string) map[string]player.CoinOwner {
	var highestCoins uint64 = 0
	for name, player := range g.players {
		if player.Coins() > highestCoins && name != exclude {
			highestCoins = player.Coins()
		}
	}

	m := make(map[string]player.CoinOwner)

	for name, player := range g.players {
		if player.Coins() == highestCoins && player.Name() != exclude {
			m[name] = player
		}
	}

	return m
}

func (g *Game) TakeCourthouse() (courthouse uint64) {
	courthouse = g.courthouse
	g.courthouse = 0
	return
}

func (g *Game) CheaterWins(p *player.Player) {
	g.cheatWinner = p
}

func (g *Game) RevealCard(p *player.Player) {
	p.Reveal(g.turnCount)
}

func (g *Game) Winners() []string {
	return g.winners
}
