package game

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/petertseng/mascarade/format"
	"github.com/petertseng/mascarade/player"
	"github.com/petertseng/mascarade/role"
)

type Power func(game GameResolver, user *player.Player, numCorrect int, format format.Formatter)

func waitForSwappables(choiceGetter func(string) []string, format format.Formatter, user string, possibleChoices map[string]player.Swappable, num int) []player.Swappable {
	for {
		names := choiceGetter(user)
		choices := make([]player.Swappable, 0)
		seen := make(map[string]bool)
		for _, name := range names {
			if _, ok := seen[name]; ok {
				format.Error(user, fmt.Errorf("You must select %d different cards but you selected %s twice", num, name))
				continue
			}

			choice, ok := possibleChoices[name]
			if ok {
				choices = append(choices, choice)
				seen[name] = true
			} else {
				format.Error(user, fmt.Errorf("No such player %s", name))
			}
		}
		if len(choices) != num {
			format.Error(user, fmt.Errorf("You must select %d players", num))
			continue
		}
		return choices
	}
}

func waitForPlayers(choiceGetter func(string) []string, format format.Formatter, user string, possibleChoices map[string]*player.Player, num int) []*player.Player {
	for {
		names := choiceGetter(user)
		choices := make([]*player.Player, 0)
		seen := make(map[string]bool)
		for _, name := range names {
			if _, ok := seen[name]; ok {
				format.Error(user, fmt.Errorf("You must select %d different players but you selected %s twice", num, name))
				continue
			}

			choice, ok := possibleChoices[name]
			if ok {
				choices = append(choices, choice)
				seen[name] = true
			} else {
				format.Error(user, fmt.Errorf("No such player %s", name))
			}
		}
		if len(choices) != num {
			format.Error(user, fmt.Errorf("You must select %d players", num))
			continue
		}
		return choices
	}
}

func waitForCoinOwner(choiceGetter func(string) []string, format format.Formatter, user string, possibleChoices map[string]player.CoinOwner) player.CoinOwner {
	for {
		names := choiceGetter(user)
		if len(names) == 0 {
			continue
		}
		choice, ok := possibleChoices[names[0]]
		if ok {
			return choice
		} else {
			format.Error(user, fmt.Errorf("No such player %s", names[0]))
		}
	}
}

func waitForRole(choiceGetter func(string) []string, format format.Formatter, user string) role.Role {
	for {
		choices := choiceGetter(user)
		if len(choices) == 0 {
			continue
		}
		r, err := role.FromString(choices[0])
		if err == nil {
			return r
		} else {
			format.Error(user, err)
		}
	}
}

func waitForBoolean(choiceGetter func(string) []string, format format.Formatter, user string) bool {
	for {
		choices := choiceGetter(user)
		if len(choices) == 0 {
			continue
		}
		actual, err := strconv.ParseBool(choices[0])
		if err == nil {
			return actual
		} else {
			format.Error(user, err)
		}
	}
}

var powers map[role.Role]Power = map[role.Role]Power{
	role.Judge: func(game GameResolver, user *player.Player, _ int, format format.Formatter) {
		coins := game.TakeCourthouse()
		user.AddCoins(coins)
		format.GainCoins(user.Name(), coins, user.Coins())
		format.Courthouse(0)
	},

	role.Bishop: func(game GameResolver, user *player.Player, _ int, format format.Formatter) {
		richest := game.RichestOtherThan(user.Name())
		var coinOwner player.CoinOwner
		if len(richest) == 1 {
			for _, co := range richest {
				coinOwner = co
			}
		} else {
			names := make([]string, 0)
			for name, _ := range richest {
				names = append(names, name)
			}
			format.PromptForPlayer(user.Name(), role.Bishop, 1, fmt.Sprintf("The richest players are: %s.", strings.Join(names, ", ")))
			coinOwner = waitForCoinOwner(game.UserChoice, format, user.Name(), richest)
		}

		coinOwner.Pay(user, 2)
		format.PayCoins(coinOwner.Name(), coinOwner.Coins(), 2, user.Name(), user.Coins())
	},

	role.King: func(game GameResolver, user *player.Player, _ int, format format.Formatter) {
		user.AddCoins(3)
		format.GainCoins(user.Name(), 3, user.Coins())
	},

	role.Fool: func(game GameResolver, user *player.Player, _ int, format format.Formatter) {
		user.AddCoins(1)
		format.GainCoins(user.Name(), 1, user.Coins())

		swappables := game.SwappablesOtherThan(user.Name())
		format.PromptForSwappable(user.Name(), role.Fool, 2)
		choices := waitForSwappables(game.UserChoice, format, user.Name(), swappables, 2)
		format.PromptForSwap(user.Name())
		actualSwap := waitForBoolean(game.UserChoice, format, user.Name())
		if actualSwap {
			choices[0].SwapRoles(choices[1])
		}
	},

	role.Queen: func(game GameResolver, user *player.Player, _ int, format format.Formatter) {
		user.AddCoins(2)
		format.GainCoins(user.Name(), 2, user.Coins())
	},

	role.Thief: func(game GameResolver, user *player.Player, _ int, format format.Formatter) {
		victim1, victim2, err := game.CoinOwnersNextTo(user.Name())
		if err != nil {
			format.Error(user.Name(), err)
		}
		victim1.Pay(user, 1)
		format.PayCoins(victim1.Name(), victim1.Coins(), 1, user.Name(), user.Coins())
		victim2.Pay(user, 1)
		format.PayCoins(victim2.Name(), victim2.Coins(), 1, user.Name(), user.Coins())
	},

	role.Witch: func(game GameResolver, user *player.Player, _ int, format format.Formatter) {
		coinOwners := game.CoinOwners()
		format.PromptForPlayer(user.Name(), role.Witch, 1, "Choose yourself to not swap.")
		coinOwner := waitForCoinOwner(game.UserChoice, format, user.Name(), coinOwners)
		var coinsToGive uint64
		var giver, receiver player.CoinOwner

		// Who has more coins?
		if coinOwner.Coins() > user.Coins() {
			coinsToGive = coinOwner.Coins() - user.Coins()
			giver = coinOwner
			receiver = user
		} else if coinOwner.Coins() < user.Coins() {
			coinsToGive = user.Coins() - coinOwner.Coins()
			giver = user
			receiver = coinOwner
		}

		if coinsToGive != 0 {
			giver.Pay(receiver, coinsToGive)
			format.PayCoins(giver.Name(), giver.Coins(), coinsToGive, receiver.Name(), receiver.Coins())
		}
	},

	role.Peasant: func(game GameResolver, user *player.Player, numCorrect int, format format.Formatter) {
		if numCorrect == 2 {
			user.AddCoins(2)
			format.GainCoins(user.Name(), 2, user.Coins())
		} else {
			user.AddCoins(1)
			format.GainCoins(user.Name(), 1, user.Coins())
		}
	},

	role.Spy: func(game GameResolver, user *player.Player, _ int, format format.Formatter) {
		format.TellOwnCard(user.Name(), user.Role())
		swappables := game.SwappablesOtherThan(user.Name())
		format.PromptForSwappable(user.Name(), role.Spy, 1)
		choice := waitForSwappables(game.UserChoice, format, user.Name(), swappables, 1)[0]
		format.TellCard(user.Name(), choice.Name(), choice.Role())
		format.PromptForSwap(user.Name())
		actualSwap := waitForBoolean(game.UserChoice, format, user.Name())
		if actualSwap {
			user.SwapRoles(choice)
		}
	},

	role.Cheat: func(game GameResolver, user *player.Player, _ int, format format.Formatter) {
		if user.Coins() >= 10 {
			game.CheaterWins(user)
		}
	},

	role.Inquisitor: func(game GameResolver, user *player.Player, _ int, format format.Formatter) {
		players := game.PlayersOtherThan(user.Name())
		format.PromptForPlayer(user.Name(), role.Inquisitor, 1, "")
		choice := waitForPlayers(game.UserChoice, format, user.Name(), players, 1)[0]
		format.PromptForRole(choice.Name())
		guess := waitForRole(game.UserChoice, format, choice.Name())

		game.RevealCard(choice)

		if guess == choice.Role() {
			format.GoodClaim(choice.Name(), guess)
		} else {
			format.BadClaim(choice.Name(), choice.Role(), guess)
			choice.Pay(user, 4)
			format.PayCoins(choice.Name(), choice.Coins(), 4, user.Name(), user.Coins())
		}
	},

	role.Widow: func(game GameResolver, user *player.Player, _ int, format format.Formatter) {
		if user.Coins() < 10 {
			coinsToGive := 10 - user.Coins()
			user.AddCoins(10 - user.Coins())
			format.GainCoins(user.Name(), coinsToGive, 10)
		}
	},
}
