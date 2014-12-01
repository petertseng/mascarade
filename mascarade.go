package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/petertseng/mascarade/game"
)

func main() {
	gameBuilder := game.NewBuilder()

	if len(os.Args) == 1 {
		fmt.Printf("usage: %s num_players player1 player2... playerN role1 role2... roleN\n", os.Args[0])
		return
	}

	numPlayers, err := strconv.ParseInt(os.Args[1], 0, 64)
	if err != nil {
		fmt.Println(err)
		return
	}

	if int64(len(os.Args)) < numPlayers+2 {
		fmt.Printf("Expected %d player names, but only have %d\n", numPlayers, len(os.Args)-2)
		return
	}

	for i := 0; int64(i) < numPlayers; i++ {
		gameBuilder.AddPlayer(os.Args[i+2])
	}

	for _, role := range os.Args[numPlayers+2:] {
		gameBuilder.AddRole(role)
	}

	game, err := gameBuilder.MakeGame(os.Stdout)
	if err != nil {
		fmt.Println(err)
		return
	}

	input := bufio.NewReader(os.Stdin)
	for len(game.Winners()) == 0 {
		str, err := input.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}

		fields := strings.Fields(str)
		if len(fields) == 0 {
			continue
		}

		switch strings.ToLower(fields[0]) {
		case "swap":
			if len(fields) >= 3 {
				var actual bool
				actual, err = strconv.ParseBool(fields[2])
				if err == nil {
					err = game.SwapOrNot(fields[1], actual)
				}
			} else {
				err = fmt.Errorf("usage: swap <player_or_table> <actually_swap>")
			}
		case "peek":
			err = game.Peek()
		case "claim":
			if len(fields) >= 2 {
				err = game.ClaimRole(fields[1])
			} else {
				err = fmt.Errorf("usage: claim <role>")
			}
		case "cc":
			err = game.Challenge()
		case "pass":
			err = game.NoChallenge()
		default:
			err = fmt.Errorf("usage: <swap|peek|claim|cc|pass> [args]")
		}

		if err != nil {
			fmt.Println(err)
		}
	}
}
