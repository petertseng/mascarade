# mascarade

This is a Go implementation of the game [Mascarade](https://boardgamegeek.com/boardgame/139030/mascarade) by Bruno Faidutti.

# Usage

The package `mascarade/game` contains everything needed to create a game.
create a `NewBuilder`, use `AddPlayer` and `AddRole`, and then `MakeGame` to start the game.

Once the game has started, call `SwapOrNot`, `Peek`, `ClaimRole`, `Challenge`, or `NoChallenge` to perform the respective actions.

When specifying a target to swap with, use #0, #1, #2... etc. to swap with table cards, or a player's name to swap with that player.

`mascarade.go` contains an example that simply runs a game using standard input and standard output.
See the usage message for details on invocation.

## Future work

None of the [Mascarade Expansion](https://boardgamegeek.com/boardgame/163107/mascarade-expansion) characters are implemented, though they are listed in the game data.
