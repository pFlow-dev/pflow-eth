pflow-eth
=========

Pflow.eth is a Web3 Metamodel Repository.

Embed links to formal models of your smart contract.
Build models with the Solidity Metamodel libraries.

Motivation
----------

We want to compose smart contracts that are Visualizable and Formally verified.
Why not embed the proofs directly in the code?

Status
------
WIP / Alpha


Testing on  Sepolia Testnet: https://sepolia.etherscan.io/address/0x33908630f6ee75e05d78cfbfaa4a0e3dc76c8c6e#code

Testing on Scroll Testnet: https://sepolia-blockscout.scroll.io/address/0x2A9862692E1d681dA5986B4C96D582dD0Ef29433

```
Metadata and sources of "tictactoe" were published successfully.
contracts/TicTacToe_flattened.sol : 
dweb:/ipfs/QmVymBsnnW9gKVCD9XyJENtP2pmetggtMJKR21gT6i35nR
metadata.json : 
dweb:/ipfs/QmT6LNecQRawvpTH4maofgFjgWXsRDMjMYv6ajBfhJTkg7
```


Models
======

TicTacToe
---------

The first example developed in solidity. Compare with an interactive JS model at https://pflow.dev

![image](https://github.com/pFlow-dev/pflow-eth/assets/243500/10f14403-59ce-4539-b411-75471b678530)

KonamiCode
----------

A simple example of a petri-net model that can be deployed to Ethereum.
The model is a simple state machine that recognizes the Konami Code.
