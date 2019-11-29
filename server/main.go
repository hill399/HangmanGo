// Hangman CLI server Package
// Runs at localhost:8080
// Author: hill399

// Usage: Launches http server which the client-side application can interact with.
// "/newgame" Generates new game and stores active game data.
// "/games" Generates list of all currently running games.
// "/guess" Accepts and evaluates user guesses.

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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

/* Array to store created games */
var openGames []gameStore

/* Main function */
func main() {

	/* http handler for "/newgame" */
	newGameHandler := func(w http.ResponseWriter, req *http.Request) {
		gameNo := newGame()
		fmt.Fprintf(w, "Game %s Created", strconv.Itoa(gameNo))
	}

	/* http handler for "/games" */
	/* Iterates over openGames array and prints to client list of games and status */
	openGamesHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "GAME ID | PLAYABLE | TURNS |      WORD STATE      |  WINNER \n")
		for _, game := range openGames {
			io.WriteString(w, strconv.Itoa(game.gameID)+
				" 	  "+strconv.FormatBool(game.gameState)+
				"        "+strconv.Itoa(game.turns)+
				"       "+strings.Join(game.completeWord, ",")+
				"        "+game.winner+"\n")
		}
	}

	/* http handler for "/guess" */
	/* http POST request with game/guess data in body */
	guessHandler := func(w http.ResponseWriter, req *http.Request) {
		/* Read and isolate guess data */
		body, _ := ioutil.ReadAll(req.Body)
		decBody := strings.Split(string(body), ",")
		gameNo, _ := strconv.Atoi(decBody[0])
		guess := strings.ToLower(decBody[1])
		username := decBody[2]

		/* mutex lock game to alter for concurrency purposes */
		openGames[gameNo].mux.Lock()

		/* Logic to check game is active and evaluate guess */
		if openGames[gameNo].gameState != true {
			io.WriteString(w, "Game is finished, cannot make guess\n")
			/* Function to check if letter has already been played */
		} else if checkUsedLetters(gameNo, guess) != true {
			io.WriteString(w, "Letter already played, try again\n")
		} else {
			/* Counter to track number of letters found on guess */
			var lettersFound int

			/* Add guessed letter to burnt letters array */
			openGames[gameNo].lettersGuessed = append(openGames[gameNo].lettersGuessed, guess)

			/* Loop through win word and evaluate against char guess */
			for i := range openGames[gameNo].playWord {
				if openGames[gameNo].playWord[i] == guess {
					openGames[gameNo].completeWord[i] = openGames[gameNo].playWord[i]
					lettersFound++
				}
			}

			fmt.Fprintf(w, "%d Correct letters found!\n", lettersFound)

			/* If char not found, reduce number of turns */
			if lettersFound == 0 {
				openGames[gameNo].turns--
				if openGames[gameNo].turns == 0 {
					openGames[gameNo].gameState = false
				}
			}

			/* Evaluate state of guess word to determine win condition */
			win := true
			for i, letter := range openGames[gameNo].completeWord {
				if openGames[gameNo].playWord[i] != letter {
					win = false
				}
			}

			/* Modify game state and set winner if game over */
			if win == true {
				fmt.Fprintf(w, "You are the winner of Game %d!\n", gameNo)
				openGames[gameNo].winner = username
				openGames[gameNo].gameState = false
			}

			/* Server-side debug statement */
			//fmt.Println(openGames[gameNo])

			/* Write game state to client after turn in complete */
			io.WriteString(w, "GAME ID | PLAYABLE | TURNS | WORD STATE | WINNER \n")
			io.WriteString(w,
				strconv.Itoa(openGames[gameNo].gameID)+
					" 	  "+strconv.FormatBool(openGames[gameNo].gameState)+
					"        "+strconv.Itoa(openGames[gameNo].turns)+
					" 	 "+strings.Join(openGames[gameNo].completeWord, ",")+
					"    "+openGames[gameNo].winner+"\n")
		}

		/* Unlock mutex to allow for next user to attempt */
		openGames[gameNo].mux.Unlock()
	}

	fmt.Println("---------------------------")
	fmt.Println("Hangman CLI Server Side App")
	fmt.Println("---------------------------")
	fmt.Println("Listening at localhost:8080")

	/* Set http function handlers and start listening on port 8080 */
	http.HandleFunc("/newgame", newGameHandler)
	http.HandleFunc("/games", openGamesHandler)
	http.HandleFunc("/guess", guessHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

/* Evaluates user guess against array of characters already used in active game */
func checkUsedLetters(gameNo int, guess string) bool {
	fmt.Println("guess", guess)
	for _, letter := range openGames[gameNo].lettersGuessed {
		if letter == guess {
			return false
		}
	}
	return true
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
	tempStore := gameStore{gameID: len(openGames), gameState: true, playWord: tempPlayWord, completeWord: tempCompleteWord, turns: 8, winner: "N/A"}
	fmt.Println(tempStore)
	openGames = append(openGames, tempStore)

	/* return game ID */
	return tempStore.gameID
}
