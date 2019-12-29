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
)

/* Array to store created games */
var openGames []gameStore

/* Main function */
func main() {

	/* http handler for "/newgame" */
	newGameHandler := func(w http.ResponseWriter, req *http.Request) {
		gameNo := newGame()
		fmt.Fprintf(w, "Game %d Created", gameNo)
	}

	/* http handler for "/games" */
	/* Iterates over openGames array and prints to client list of games and status */
	openGamesHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "GAME ID | WINNER | PLAYABLE | TURNS | WORD STATE \n")
		for _, game := range openGames {
			game.PrintGame(w)
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

		pGame := &openGames[gameNo]

		/* Check if game is active */
		pGame.IsGameActive(w)

		var validLetter bool

		/* If game active, validate letter */
		if pGame.gameState == true {
			validLetter = pGame.IsLetterValid(w, guess)
		}

		if pGame.gameState == true && validLetter == true {
			/* Loop through win word and evaluate against char guess */
			pGame.EvaluateGuess(w, guess)

			pGame.EvaluateWinState(username)

			if pGame.turns == 0 {
				pGame.gameState = false
			}

			if pGame.gameState != true {
				if pGame.turns == 0 {
					fmt.Fprintf(w, "You lose, Game %d over!\n", gameNo)
				} else {
					fmt.Fprintf(w, "You are the winner of Game %d!\n", gameNo)
				}
			}
		}

		/* Write game state to client after turn in complete */
		io.WriteString(w, "GAME ID | WINNER | PLAYABLE | TURNS | WORD STATE \n")
		pGame.PrintGame(w)

		/* Print to server console */
		fmt.Printf("Guess made on game %d\n", gameNo)
		fmt.Println(pGame)

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
