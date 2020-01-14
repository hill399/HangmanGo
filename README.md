# HangmanGo
Simple hangman game written in Go.

gRPC-connected game with the following features:
- Able to run multiple game instances
- Players permitted to play any running game instance
- Tracks turns taken, game state and winner (if any)
- Client-side CLI interface

Consists of the following components:

## Server
Utilises `grpc` to start a server running by default at `localhost:50051`. 

Links game functions into server so that requests to the below RPC endpoints can be used to change/view the game state. 

`NewGame`: Generates new game template and pushes it into active games array.

`List`: Retrieves list of currently open games.

`Guess`: Evaluates validity of user guess and processes guess. Determines win/lose state.


## Client

Interacts with the server via RPC requests. Control is handled by CLI interface `urfave/cli`.

`newgame`: Generates new game at server and responds with game no. created.

`listgames`: Retrieves list of active games.

`guess [game_no] [letter_guess] [username (opt)]`: Attempts guess of game specified.


## Usage 

Build and run server with:
```
go build server.go game.go
```

Build client with:
```
go build client.go
```

Execute `/client` on client executable to see usage options.




