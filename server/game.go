package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/tjarratt/babble"
)

/* type struct to unique game data */
type gameStore struct {
	mux            sync.Mutex
	gameID         int
	gameState      bool
	playWord       []string
	completeWord   []string
	lettersGuessed []string
	turns          int
	winner         string
}

/* Creates new game and returns game ID */
func newGame() int {
	/* Initiate babbler library to use an RWG */
	babbler := babble.NewBabbler()
	babbler.Count = 1
	tempPlayWord := strings.Split(strings.ToLower(babbler.Babble()), "")
	/* Create blank play word for user to view */
	var tempCompleteWord []string
	for range tempPlayWord {
		tempCompleteWord = append(tempCompleteWord, "_")
	}

	/* Generate and push new game in active games array */
	openGames = append(openGames, gameStore{gameID: len(openGames), gameState: true, playWord: tempPlayWord, completeWord: tempCompleteWord, turns: 8, winner: "N/A"})

	/* Print status to server console */
	fmt.Printf("Game %d Created:\n", len(openGames)-1)
	fmt.Println(openGames[len(openGames)-1])

	/* return game ID */
	return openGames[len(openGames)-1].gameID
}

/* Print function to display game data to client */
func (g gameStore) PrintGame(w http.ResponseWriter) {
	fmt.Fprintf(w, "   %d	   %s       %t       %d      %s\n",
		g.gameID,
		g.winner,
		g.gameState,
		g.turns,
		strings.Join(g.completeWord, ","),
	)
}

/* Function to check if the queried game is active */
func (g gameStore) IsGameActive(w http.ResponseWriter) bool {
	if g.gameState == false || g.turns == 0 {
		io.WriteString(w, "Game is finished, cannot make guess\n")
		return false
	}

	return true
}

/* Evaluates user guess against array of characters already used in active game */
func (g gameStore) IsLetterValid(w http.ResponseWriter, guess string) (bool, []string) {
	for _, letter := range g.lettersGuessed {
		if letter == guess {
			/* Add guessed letter to burnt letters array */
			io.WriteString(w, "Letter already played, try again\n")
			return false, g.lettersGuessed
		}
	}

	g.lettersGuessed = append(g.lettersGuessed, guess)
	return true, g.lettersGuessed
}

/* Function to evaluate the clients guess against the current game state */
func (g gameStore) EvaluateGuess(w http.ResponseWriter, guess string) int {
	var ls int

	for i := range g.playWord {
		if g.playWord[i] == guess {
			g.completeWord[i] = g.playWord[i]
			ls++
		}
	}

	fmt.Fprintf(w, "%d Correct letters found!\n", ls)

	/* If char not found, reduce number of turns */
	if ls == 0 {
		g.turns = g.turns - 1
		if g.turns == 0 {
			g.gameState = false
		}
	}

	return g.turns
}

/* Function to determine if a game has been won */
func (g gameStore) EvaluateWinState(name string) (string, bool) {

	/* Evaluate state of guess word to determine win condition */
	win := true
	for i, letter := range g.completeWord {
		if g.playWord[i] != letter {
			win = false
		}
	}

	/* Modify game state and set winner if game over */
	if win == true {
		g.winner = name
		g.gameState = false
		return g.winner, g.gameState
	}

	return "N/A", true
}
