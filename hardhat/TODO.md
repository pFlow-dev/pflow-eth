Embrace the Medium
------------------

- Contract Size Limit
    - roughly 100 objects total
    - of Places, Transitions, Arcs, and Guards
- Grid
    - encourage 16*16 - for 'cozy design space'

WIP
---

- [ ] optimize solidity code exports
- [ ] experiment with composition
- [ ] is it worth +9 extra objects to get rid of Model.Position element?
    - [ ] readability is the tradeoff

DONE
----

- [x] ? Why is cost of a guard 3x vs Func ? A: Arc index is unnecessary
- [x] REVIEW init code sizes

```
Jetsam Cost per object (in Kib)

cell() >>> 2.104 / 22  = 0.09563636363636364

func() >>> 4.187 / 48 = 2916666666668

guard() | arrow() >>> 24.912 / 89  = 0.27991011235955054

```
