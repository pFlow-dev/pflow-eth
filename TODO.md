WIP
---
This repo is currently a testbed for working with multiple streams of contract events
TODO: experiment with multiple streams of contract events & subscribe/store

BACKLOG
-------
- [ ] finish unpacking the petri-net model from (on-chain call or server)
- [ ] use model definition to populate UI labels for known values
- [ ] upgrade web2/web3 when wallet is connected - send signal command from wallet when connected instead of API
- [ ] make 'transaction builder' using the petri-net GUI

ICEBOX
-------
- [ ] when viewing model history let left/right arrows show what happened in each frame
- [ ] (UX) make status / green / ok compare w/ timeclock
- [ ] consider optimizing sync by monitoring the sequence value in the contract
- [ ] support monitoring multiple contracts

DONE
----
- [x] make faucet route work
- [x] add control panel for developers to reset/restart/update config
- [x] Connect Wallet button (use ethers - since it seems typechain generates with this)
- [x] include the network id in the status bar
- [x] REVIEW: A: (no) should we build a top-level callback system? to trigger gui refresh
  ? should this object still be called Metamodel
