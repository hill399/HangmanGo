// Hangman CLI server Package
// Runs at localhost:50051
// Author: hill399

// Usage: Launches rpc server which the client-side application can interact with.
// "NewGame" Generates new game and stores active game data.
// "List" Generates list of all currently running games.
// "Guess" Accepts and evaluates user guesses.


package main

import (
	"fmt"
	"log"
	"context"
	"net"

	"github.com/hill399/HangmanGo/hangmanpb"
	"google.golang.org/grpc"
)

type server struct{}

/* Array to store created games */
var openGames []gameStore

func main() {
	fmt.Println("---------------------------")
	fmt.Println("Hangman CLI Server Side App")
	fmt.Println("---------------------------")
	fmt.Println("Listening at localhost:50051")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	hangmanpb.RegisterGuessServiceServer(s, &server{})
	hangmanpb.RegisterNewGameServiceServer(s, &server{})
	hangmanpb.RegisterListServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}

func (*server) Guess(ctx context.Context, req *hangmanpb.GuessRequest) (*hangmanpb.GuessResponse, error) {

	fmt.Printf("Guess function was invoked with %v", req)

	gameNo := req.GetGuess().GetGameNumber()
	guess := req.GetGuess().GetGuessLetter()
	username := req.GetGuess().GetUsername()

	det := []string{}

	/* mutex lock game to alter for concurrency purposes */
	openGames[gameNo].mux.Lock()

	pGame := &openGames[gameNo]

	/* Check if game is active */
	pGame.IsGameActive(&det)

	var validLetter bool

	/* If game active, validate letter */
	if pGame.gameState == true {
		validLetter = pGame.IsLetterValid(guess, &det)
	}

	if pGame.gameState == true && validLetter == true {
		/* Loop through win word and evaluate against char guess */
		pGame.EvaluateGuess(guess, &det)

		pGame.EvaluateWinState(username)

		if pGame.turns == 0 {
			pGame.gameState = false
		}

		if pGame.gameState != true {
			if pGame.turns == 0 {
				det = append(det, fmt.Sprintf("You lose, Game %d over!\n", gameNo))
			} else {
				det = append(det, fmt.Sprintf("You are the winner of Game %d!\n", gameNo))
			}
		}

		/* Print to server console */
		fmt.Printf("Guess made on game %d\n", gameNo)
		fmt.Println(pGame)
	}

	/* Unlock mutex to allow for next user to attempt */
	openGames[gameNo].mux.Unlock()

	fmt.Println("Printing detail slice")
	for _, line := range det {
		fmt.Println(line)
	}

	res := &hangmanpb.GuessResponse{
		Response: pGame.PrintGame(),
		Detail: det,
	}

	return res, nil
}

func (*server) NewGame(ctx context.Context, req *hangmanpb.NewGameRequest) (*hangmanpb.NewGameResponse, error) {
	fmt.Printf("NewGame function was invoked")

	gameNo := newGame()

	res := &hangmanpb.NewGameResponse{
		GameNumber: int32(gameNo),
	}

	return res, nil
}

func (*server) List(ctx context.Context, req *hangmanpb.ListRequest) (*hangmanpb.ListResponse, error) {
	fmt.Printf("List function was invoked") 

	sa := make([]string, len(openGames))
	
	sa = append(sa, fmt.Sprintf("\nGAME ID| WINNER | PLAYABLE | TURNS | WORD STATE\n"))

	for _, game := range openGames {
		sa = append(sa, game.PrintGame())
	}

	res := &hangmanpb.ListResponse{
		GameDetails: sa,
	}

	return res, nil
}
