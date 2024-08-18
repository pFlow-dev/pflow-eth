WIP
---
enhance UX for general usage 

- [ ] add a header/footer to match pflow.xyz 
- [ ] add pflow logo
- [ ] add interval job to monitor sequence


BACKLOG
-------

- [ ] shift+click to select multiple transitions and adjust repeat count (scalar values)
- [ ] support using multiple contracts(allow input address)
- [ ] experiment with multiple streams of contract events & subscribe/store


ICEBOX
-------
- [ ] make a flow for users to register their contracts for transaction monitoring on pflow.xyz
- [ ] consider viewing model history: left/right arrows show what happened in each frame

- [ ] build an actual state channel implementation where users have a shared multi-sig + individual events

      monitor the pendingNonce for each address in the channel,
      each party sends a new proposed next command, other users sign (w/ a multi-sig?)
      or
      each user builds up state w/ appropriate permissions,
      the current multiSignal transaction is latest state.

DONE
----
- [x] launch MVP: sepolia-optimism.pflow.xyz:
- [x] upgrade web2/web3 when wallet is connected - send signal command from wallet when connected instead of API
- [x] make 'transaction builder' using the petri-net GUI
- [x] finish anonymous session support - 'log in' w/ websocket
- [x] add session time-out
- [x] use model definition to populate UI labels for known values
- [x] finish unpacking the petri-net model from (on-chain call or server)
- [x] enhance the pflowDSL so metamodel declaration looks more like solidity


MULTISIG NOTES
--------------

Build a state channel, where each party sends individual events to the aggregator
then when both parties agree - sign and send multi-sig payment, signalMany()

?? what keeps each transaction from being sent out-of-order?

I guess the protocol can just be that each party proposes a new signalMany() call that the other confirms/signs... 

-- likely the address cannot issue any other transactions (via nonce) - maybe give a way to collapse the state channel
