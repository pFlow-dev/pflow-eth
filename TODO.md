Overview

_Goals:_

* Help developers build large on-chain structures securely
* Reach out to programers who want to build state-machines on Solidity
  * provide code and documentation
* Build a really large network of composable models on and off chain

_Why?_

Meta-modeling can live between 'programing' and 'AI' - ai models can be derived from meta-models by generating data sets.

Why on the Blockchain?

It has a cost for use, but is immutable and will live forever.

_Why now?_

"We're gonna need some bigger Hyperstructures".  

We want to build data-composability on top of a token infrastructure that does _not_ depend on payable tokens.


WIP
---
* [ ] build out JS front end - re-deploy to pflow.eth
* [ ] use quicknode callbacaks in Dapp - deploy on sepolia
* [ ] deploy end-to-end POC for TicTacToe
* [ ] deploy end-to-end POC for KonamiCode


BACKLOG
-------
* [ ] test quicknode callbacks
* [ ] test w/ Quicknode on Sepolia Testnet w/ quick-alerts
* [ ] investigate using https://eips.ethereum.org/EIPS/eip-712 for offline message signing ?  what if we just send events w/o state?

DONE
----
* [x] test contract locally
* [x] provision callbacks & test w/ Quicknode on Sepolia Testnet
* [x] investigate Q: can we emit the whole model during RESET? [A: why tho?]
* [x] prototype access control for models that have been sucessfully added to the registry
* [x] explore eth logs - https://ethereum.stackexchange.com/questions/1686/ethereum-event-log-maximum-size#:~:text=There%20is%20no%20size%20limit,how%20much%20Ether%20you%20have


ICEBOX
------
* [ ] should we adopt IPFS pubsub? https://github.com/ipfs/js-ipfs/blob/master/docs/core-api/PUBSUB.md
* [ ] REVIEW: how to add signed data https://github.com/rarible/protocol-contracts/blob/master/meta-tx/contracts/EIP712MetaTransaction.sol
* [ ] plan ERC721 adoption for metamodel publishing
* [ ] REVIEW: adopt ethscrition/calldata type model for off-chain rollups
* [ ] could/should we have our own consensus layer? https://chainsafe.github.io/lodestar/
      - what is a beacon-node?
* [ ] https://soliditydeveloper.com/merkle-tree - good overview and links to actual solidity examples
* [ ] patricia trie or mekle tree https://github.com/lidofinance/curve-merkle-oracle/blob/main/contracts/MerklePatriciaProofVerifier.sol
* [ ] review approach to provenance https://boredapeyachtclub.com/#/provenance
* [ ] review alternate code approached to storage to make smaller models
      https://dev.to/mudgen/solidity-libraries-can-t-have-state-variables-oh-yes-they-can-3ke9


