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

func (pGame *gameStore) PrintGame(w http.ResponseWriter) {
	fmt.Fprintf(w, "   %d	   %s       %t       %d      %s\n",
		(*pGame).gameID,
		(*pGame).winner,
		(*pGame).gameState,
		(*pGame).turns,
		strings.Join((*pGame).completeWord, ","),
	)
}

func (pGame *gameStore) IsGameActive(w http.ResponseWriter) {
	if (*pGame).gameState == false || (*pGame).turns == 0 {
		io.WriteString(w, "Game is finished, cannot make guess\n")
	}
}

/* Evaluates user guess against array of characters already used in active game */
func (pGame *gameStore) IsLetterValid(w http.ResponseWriter, guess string) bool {
	for _, letter := range (*pGame).lettersGuessed {
		if letter == guess {
			/* Add guessed letter to burnt letters array */
			io.WriteString(w, "Letter already played, try again\n")
			return false
		}
	}

	(*pGame).lettersGuessed = append((*pGame).lettersGuessed, guess)
	return true
}

func (pGame *gameStore) EvaluateGuess(w http.ResponseWriter, guess string) {
	var ls int

	for i := range (*pGame).playWord {
		if (*pGame).playWord[i] == guess {
			(*pGame).completeWord[i] = (*pGame).playWord[i]
			ls++
		}
	}

	fmt.Fprintf(w, "%d Correct letters found!\n", ls)

	/* If char not found, reduce number of turns */
	if ls == 0 {
		(*pGame).turns = (*pGame).turns - 1
		if (*pGame).turns == 0 {
			(*pGame).gameState = false
		}
	}
}

func (pGame *gameStore) EvaluateWinState(name string) {

	/* Evaluate state of guess word to determine win condition */
	win := true
	for i, letter := range (*pGame).completeWord {
		if (*pGame).playWord[i] != letter {
			win = false
		}
	}

	/* Modify game state and set winner if game over */
	if win == true {
		(*pGame).winner = name
		(*pGame).gameState = false
	}
}
