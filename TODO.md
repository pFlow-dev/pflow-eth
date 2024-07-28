WIP
---
- [ ] launch MVP: sepolia-optimism.pflow.xyz:
  (no transaction download/storage, just contract browsing and state machine actions)

- [ ] support using multiple contracts(allow inputs)
- [ ] consider optimizing sync by monitoring the sequence value in the contract ping

BACKLOG
-------
- [ ] shift+click to select multiple transitions and adjust repeat count (scalar values)

ICEBOX
-------
- [ ] make a flow for users to register their contracts for transaction monitoring on pflow.xyz
- [ ] experiment with multiple streams of contract events & subscribe/store
- [ ] consider viewing model history: left/right arrows show what happened in each frame


DONE
----
- [x] upgrade web2/web3 when wallet is connected - send signal command from wallet when connected instead of API
- [x] make 'transaction builder' using the petri-net GUI
- [x] finish anonymous session support - 'log in' w/ websocket
- [x] add session time-out
- [x] use model definition to populate UI labels for known values
- [x] finish unpacking the petri-net model from (on-chain call or server)
- [x] enhance the pflowDSL so metamodel declaration looks more like solidity
