package format

import (
	"github.com/petertseng/mascarade/role"
)

type Formatter interface {
	YourTurn(player string) error
	SwapOrNot(swapper, swapee string) error
	Peek(peeker string) error
	TellOwnCard(peeker string, r role.Role) error
	ClaimRole(claimant string, r role.Role) error

	YourTurnToChallenge(player, claimant string, r role.Role) error
	Counterclaim(claimant, original string, r role.Role) error
	NoCounterclaim(claimant, original string, r role.Role) error
	NobodyChallenged(claimant string, r role.Role) error
	GoodClaim(claimant string, r role.Role) error
	BadClaim(claimant string, had, want role.Role) error
	UsePower(user string, r role.Role) error

	GainCoins(gainer string, coins, now uint64) error
	PayFine(gainer string, now uint64) error
	PayCoins(giver string, giverCoins, paid uint64, receiver string, receiverCoins uint64) error
	Courthouse(coins uint64) error

	CheaterWins(cheater string) error
	WinTargetReached(winners []string) error
	WinBroke(winners, broke []string) error

	TellCard(peeker, whoseCard string, r role.Role) error
	PromptForRole(player string) error
	PromptForPlayer(player string, r role.Role, num int, extra string) error
	PromptForSwap(player string) error
	PromptForSwappable(player string, r role.Role, num int) error
	Error(player string, e error) error
}
