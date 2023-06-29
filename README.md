pflow-eth
=========

An SDK for deploying petri-net models on Ethereum / EVM compatible blockchains.

Status
------
WIP / Alpha


Deployed to Testnet: https://sepolia.etherscan.io/address/0x33908630f6ee75e05d78cfbfaa4a0e3dc76c8c6e#code


Models
======

TicTacToe
---------

The first example developed in solidity. Compare with an interactive JS model at https://pflow.dev

![image](https://github.com/pFlow-dev/pflow-eth/assets/243500/10f14403-59ce-4539-b411-75471b678530)


The `TicTacToe` contract represents a game of Tic-Tac-Toe implemented using a Petri net model. It allows players with the roles of X and O to make moves on the game board, validates the moves, maintains the game state, and emits events for each move made. The contract includes role-based access control and functionality to pause and reset the game.

Let's go through the code to understand its functionality:

1. **Library and Interface Declarations**: The code starts with the declaration of a library called `Uint8Model`. It defines various data structures and functions related to a Petri net model. Additionally, an interface `Uint8ModelFactory` is defined, which specifies a function `declaration()`.

2. **MetamodelUint8 and TicTacToeModel Contracts**: The contract `MetamodelUint8` is an abstract contract that implements functionality related to the Petri net model using the `Uint8Model` library. The contract `TicTacToeModel` is an abstract contract that extends `MetamodelUint8` and adds specific logic for a Tic-Tac-Toe game. It defines enums for roles, properties, and actions in the Tic-Tac-Toe game.

3. **TicTacToe Contract**: The `TicTacToe` contract is the main contract that inherits from `AccessControl` and `TicTacToeModel`. It represents the actual implementation of the Tic-Tac-Toe game.

   - **State Variables**: The contract defines several state variables, including `owner` (address of the contract deployer), `paused` (a boolean flag to pause the game), `gameId` (an integer representing the current game ID), `sequence` (an integer representing the current sequence number), and `state` (an array representing the game board state).

   - **Constructor**: The constructor initializes the contract by setting the `owner` and assigning the roles `PLAYER_X` and `PLAYER_O` to two specified addresses. It also calls the `resetGame` function to initialize the game.

   - **Internal Functions**: The contract defines several internal functions used for game logic and initialization. Notable functions include `_init`, which initializes a place in the Petri net model, `testIsGameOpen`, which checks if the game is paused, `testIsMyTurn`, which verifies if it is the caller's turn to play, and `move`, which performs a move in the game based on the specified action.

   - **Modifier and Event**: The contract defines a modifier `startGame` that resets the game state before each game starts. It also emits an event `Uint8Model.Action` whenever a move is made in the game.

   - **Public Functions**: The contract exposes several public functions, such as `pause` and `unpause` to control the game's pause state, `reset` to reset the game, and functions for each possible move on the Tic-Tac-Toe board (e.g., `X00`, `X01`, `O10`, etc.).

   - **Role Management**: The contract includes role management functionality provided by the `AccessControl` contract. It uses the roles `DEFAULT_ADMIN_ROLE`, `PLAYER_X`, and `PLAYER_O` to control access and permissions.

   - **Helper Functions**: The contract includes a helper function `getRole` to determine the current role (X or O) of the caller.

4. **Uint8ModelFactory Interface**: The `Uint8ModelFactory` interface defines a single function `declaration()`, which returns a `Uint8Model.PetriNet` struct representing the Petri net model.
