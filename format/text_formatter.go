package format

import (
	"fmt"
	"strings"

	"github.com/petertseng/mascarade/output"
	"github.com/petertseng/mascarade/role"
)

func NewText(out output.Outputter) Formatter {
	return TextFormatter{out: out}
}

type TextFormatter struct {
	out output.Outputter
}

func (tf TextFormatter) YourTurn(peeker string) error {
	_, err := tf.out.WritePublic([]byte(fmt.Sprintf("%s: It's your turn. What will ye do?\n", peeker)))
	return err
}

func (tf TextFormatter) SwapOrNot(swapper, swapee string) error {
	_, err := tf.out.WritePublic([]byte(fmt.Sprintf("%s swaps (or not) with %s.\n", swapper, swapee)))
	return err
}

func (tf TextFormatter) Peek(peeker string) error {
	_, err := tf.out.WritePublic([]byte(fmt.Sprintf("%s peeks.\n", peeker)))
	return err
}

func (tf TextFormatter) TellOwnCard(peeker string, r role.Role) error {
	_, err := tf.out.WritePrivate(peeker, []byte(fmt.Sprintf("You are the %s.\n", r)))
	return err
}

func (tf TextFormatter) ClaimRole(claimant string, r role.Role) error {
	_, err := tf.out.WritePublic([]byte(fmt.Sprintf("%s claims to be the %s (%s)!\n", claimant, r, r.PowerDescription())))
	return err
}

func (tf TextFormatter) YourTurnToChallenge(player, claimant string, r role.Role) error {
	_, err := tf.out.WritePublic([]byte(fmt.Sprintf("%s: Will you challenge %s's claim of being the %s?\n", player, claimant, r)))
	return err
}

func (tf TextFormatter) Counterclaim(claimant, original string, r role.Role) error {
	_, err := tf.out.WritePublic([]byte(fmt.Sprintf("%s also claims to be the %s!\n", claimant, r)))
	return err
}

func (tf TextFormatter) NoCounterclaim(claimant, original string, r role.Role) error {
	_, err := tf.out.WritePublic([]byte(fmt.Sprintf("%s does not challenge.\n", claimant)))
	return err
}

func (tf TextFormatter) NobodyChallenged(claimant string, r role.Role) error {
	_, err := tf.out.WritePublic([]byte(fmt.Sprintf("Nobody challenged %s's claim of %s!\n", claimant, r)))
	return err
}

func (tf TextFormatter) GoodClaim(claimant string, r role.Role) error {
	_, err := tf.out.WritePublic([]byte(fmt.Sprintf("YES, %s is indeed the %s!\n", claimant, r)))
	return err
}

func (tf TextFormatter) BadClaim(claimant string, had, want role.Role) error {
	_, err := tf.out.WritePublic([]byte(fmt.Sprintf("NO, %s was the %s, not the %s!\n", claimant, had, want)))
	return err
}

func (tf TextFormatter) UsePower(user string, r role.Role) error {
	_, err := tf.out.WritePublic([]byte(fmt.Sprintf("%s now uses the power of the %s: %s!\n", user, r, r.PowerDescription())))
	return err
}

func (tf TextFormatter) GainCoins(gainer string, coins, now uint64) error {
	_, err := tf.out.WritePublic([]byte(fmt.Sprintf("%s gains %d coins and now has %d.\n", gainer, coins, now)))
	return err
}

func (tf TextFormatter) PayFine(gainer string, now uint64) error {
	_, err := tf.out.WritePublic([]byte(fmt.Sprintf("%s pays the fine and now has %d coins.\n", gainer, now)))
	return err
}

func (tf TextFormatter) PayCoins(giver string, giverCoins, paid uint64, receiver string, receiverCoins uint64) error {
	_, err := tf.out.WritePublic([]byte(fmt.Sprintf("%s pays %d coins to %s. %s now has %d and %s now has %d.\n", giver, paid, receiver, giver, giverCoins, receiver, receiverCoins)))
	return err
}

func (tf TextFormatter) Courthouse(coins uint64) error {
	_, err := tf.out.WritePublic([]byte(fmt.Sprintf("The Courthouse now has %d coins.\n", coins)))
	return err
}

func (tf TextFormatter) CheaterWins(cheater string) error {
	_, err := tf.out.WritePublic([]byte(fmt.Sprintf("%s wins the game by cheating!\n", cheater)))
	return err
}

func (tf TextFormatter) WinTargetReached(winners []string) error {
	_, err := tf.out.WritePublic([]byte(fmt.Sprintf("%s wins the game!\n", strings.Join(winners, ", "))))
	return err
}

func (tf TextFormatter) WinBroke(winners, broke []string) error {
	_, err := tf.out.WritePublic([]byte(fmt.Sprintf("%s are broke! The richest player %s wins the game!\n", strings.Join(broke, ", "), strings.Join(winners, ", "))))
	return err
}

func (tf TextFormatter) TellCard(player, whoseCard string, r role.Role) error {
	_, err := tf.out.WritePrivate(player, []byte(fmt.Sprintf("%s is the %s.\n", whoseCard, r)))
	return err
}

func (tf TextFormatter) PromptForRole(player string) error {
	_, err := tf.out.WritePrivate(player, []byte("Choose a role.\n"))
	return err
}

func (tf TextFormatter) PromptForPlayer(player string, r role.Role, num int, extra string) error {
	_, err := tf.out.WritePrivate(player, []byte(fmt.Sprintf("You are using %s's power: %s. Choose %d players. %s\n", r, r.PowerDescription(), num, extra)))
	return err
}

func (tf TextFormatter) PromptForSwap(player string) error {
	_, err := tf.out.WritePrivate(player, []byte(fmt.Sprintf("Would you like to swap these two cards?\n")))
	return err
}

func (tf TextFormatter) PromptForSwappable(player string, r role.Role, num int) error {
	_, err := tf.out.WritePrivate(player, []byte(fmt.Sprintf("You are using %s's power: %s. Choose %d cards.\n", r, r.PowerDescription(), num)))
	return err
}

func (tf TextFormatter) Error(player string, e error) error {
	_, err := tf.out.WritePrivate(player, []byte(fmt.Sprintf("ERROR: %s\n", e)))
	return err
}
