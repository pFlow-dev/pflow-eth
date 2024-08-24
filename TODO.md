WIP
---
enhance UX for general usage 

- [ ] add interval job to monitor sequence


BACKLOG
-------
- [ ] shift+click to select multiple transitions and adjust repeat count (scalar values)
- [ ] support using multiple contracts(allow input address)
- [ ] experiment with multiple streams of contract events & subscribe/store


ICEBOX
-------
- [ ] add go lang backend to actively monitor blockchains... 

DONE
----


MULTISIG NOTES
--------------

Build a state channel, where each party sends individual events to the aggregator
then when both parties agree - sign and send multi-sig payment, signalMany()

?? what keeps each transaction from being sent out-of-order?

I guess the protocol can just be that each party proposes a new signalMany() call that the other confirms/signs... 

-- likely the address cannot issue any other transactions (via nonce) - maybe give a way to collapse the state channel
