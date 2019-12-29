// Hangman CLI client package
// Author: hill399

// Usage: CLI client interface which allows interaction with running server-side application.
// "newgame" Generates new game on server.
// "listgames" Generates list of all currently running games on server.
// "guess" Takes game no., letter guess and optional username for server interaction.
package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/urfave/cli"
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
	app.Commands = []cli.Command{
		{
			/* Create new game - calls "/newgame" handler on server-side */
			Name:    "newgame",
			Aliases: []string{"n"},
			Usage:   "Query the server to start a new game",
			Action: func(c *cli.Context) error {
				/* Get body response and check for errors */
				respN, err := http.Get("http://localhost:8080/newgame")
				if err != nil {
					return err
				}
				defer respN.Body.Close()
				/* Readout body of return response and check for errors */
				bodyN, err := ioutil.ReadAll(respN.Body)
				if err != nil {
					return err
				}
				/* Print returned new game status */
				fmt.Println(string(bodyN))
				return nil
			},
		},
		{
			/* List open games - calls "/games" handler on server-side */
			Name:    "listgames",
			Aliases: []string{"l"},
			Usage:   "Print list of currently open games",
			Action: func(c *cli.Context) error {
				/* Get body response and check for errors */
				respL, err := http.Get("http://localhost:8080/games")
				if err != nil {
					fmt.Println("Failed to start new game")
					return err
				}
				defer respL.Body.Close()
				/* Readout body of return response and check for errors */
				bodyL, err := ioutil.ReadAll(respL.Body)
				if err != nil {
					fmt.Println("Failed to start new game")
					return err
				}
				/* Print returned game list */
				fmt.Println(string(bodyL))
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

				/* Format and pack arguments into byte form for POST request */
				guessPack := []string{c.Args().Get(0), c.Args().Get(1), username}
				guessPackBytes := []byte(strings.Join(guessPack, ","))

				/* Make POST request to "/guess" handler, evaluate error */
				respG, err := http.Post("http://localhost:8080/guess", "application/json", bytes.NewBuffer(guessPackBytes))
				if err != nil {
					fmt.Println("Failed to make guess")
					return err
				}
				defer respG.Body.Close()
				/* Readout body of return response and check for errors */
				bodyG, _ := ioutil.ReadAll(respG.Body)

				/* Print returned guess response from server */
				fmt.Println(string(bodyG))
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
