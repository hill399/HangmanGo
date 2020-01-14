// Hangman CLI client package
// Author: hill399

// Usage: CLI client interface which allows interaction with running server-side application.
// "newgame" Generates new game on server.
// "listgames" Generates list of all currently running games on server.
// "guess" Takes game no., letter guess and optional username for server interaction.
package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"unicode"
	"context"

	"github.com/urfave/cli"
	"github.com/hill399/HangmanGo/hangmanpb"
	"google.golang.org/grpc"
)

/* Main client function */
func main() {

	/* Initiate new CLI application for CLI user interfacing */
	app := cli.NewApp()
	app.Name = "Hangman CLI"
	app.Usage = "Client side CLI for hangman application"
	app.Version = "1.0.0"

	/* Creation of CLI functions and documentation */
	/* "Action" is effective function call */
	app.Commands = []*cli.Command{
		{
			/* Create new game - calls "/newgame" handler on server-side */
			Name:    "newgame",
			Aliases: []string{"n"},
			Usage:   "Query the server to start a new game",
			Action: func(c *cli.Context) error {

				cc, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())

				if err != nil {
					return err
				}
			
				defer cc.Close()

				sc := hangmanpb.NewNewGameServiceClient(cc)

				req := &hangmanpb.NewGameRequest{}

				res, err := sc.NewGame(context.Background(), req)
			
				if err != nil {
					log.Fatalf("Error while calling New Game rpc: %v", err)
				}
			
				log.Printf("Game %v Created", res.GameNumber)

				return nil
			},
		},
		{
			/* List open games - calls "/games" handler on server-side */
			Name:    "listgames",
			Aliases: []string{"l"},
			Usage:   "Print list of currently open games",
			Action: func(c *cli.Context) error {

				cc, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())

				if err != nil {
					return err
				}
			
				defer cc.Close()

				sc := hangmanpb.NewListServiceClient(cc)

				req := &hangmanpb.ListRequest{}

				res, err := sc.List(context.Background(), req)
			
				if err != nil {
					log.Fatalf("Error while calling List Game rpc: %v", err)
				}
			
				log.Printf("Game Details:\n %v", res.GameDetails)

				return nil
			},
		},
		{
			/* List open games - calls "/guess" handler on server-side */
			Name:    "guess",
			Aliases: []string{"g"},
			Usage:   "guess [game number (int)] [guess letter (char)] [optional_username string]",
			Action: func(c *cli.Context) error {
				/* Isolate user arguments for evaluation */
				gameNo := c.Args().Get(0)
				gameGuess := c.Args().Get(1)
				username := c.Args().Get(2)

				/* Set default username if omitted */
				if username == "" {
					username = "guest"
				}

				/* Parse gameNo to check if integer */
				_, err := strconv.ParseInt(gameNo, 10, 8)
				if err != nil {
					return errors.New("Invalid param - game no")
				}

				/* Parse gameGuess to assess if single character */
				for i, l := range gameGuess {
					if !unicode.IsLetter(l) || i > 0 {
						return errors.New("Invalid param - guess letter")
					}
				}

				cc, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())

				if err != nil {
					return err
				}
			
				defer cc.Close()

				sc := hangmanpb.NewGuessServiceClient(cc)

				gn, _ := strconv.Atoi(gameNo)

				req := &hangmanpb.GuessRequest{
					Guess: &hangmanpb.Guess{
						GameNumber: int32(gn),
						GuessLetter: gameGuess,
						Username: username,
					},
				}

				res, err := sc.Guess(context.Background(), req)
			
				if err != nil {
					log.Fatalf("Error while calling Guess rpc: %v", err)
				}
			
				log.Printf("Guess Response:\n %v", res.Response)
				
				for _, line := range res.Detail {
					fmt.Println(line)
				}

				return nil
			},
		},
	}

	/* Start CLI app */
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
