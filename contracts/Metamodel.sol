// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;

library Uint8Model {

    event Action(uint256 indexed session, uint8 indexed sequence, uint8 actionId, uint8 role, uint256 when);

    struct PetriNet {
        Place[] places;
        Transition[] transitions;
    }

    struct Transition {
        uint8 offset;
        uint8 role;
        int8[] delta;
        int8[] guard;
    }

    struct Place {
        uint8 offset;
        int8 initial;
        int8 capacity;
    }

}

interface Uint8ModelFactory {
    function model() external returns (Uint8Model.PetriNet memory);
}

abstract contract MetamodelUint8 is Uint8ModelFactory {

    Uint8Model.Place[] internal places;
    Uint8Model.Transition[] internal transitions;

    function transform(uint8 i, Uint8Model.Transition memory t)  internal virtual;

    function model() external view returns (Uint8Model.PetriNet memory) {
        return Uint8Model.PetriNet(places, transitions);
    }

    function cell(int8 initial, int8 capacity) internal returns (Uint8Model.Place memory) {
        Uint8Model.Place memory p =  Uint8Model.Place(uint8(places.length), initial, capacity);
        places.push(p);
        return p;
    }

    function fn(uint8 vectorSize, uint8 action, uint8 role) internal returns (Uint8Model.Transition memory) {
        require(uint8(transitions.length) == action, "Transition offset must match Actions enum");
        Uint8Model.Transition memory t = Uint8Model.Transition(action, role, new int8[](vectorSize), new int8[](vectorSize));
        transitions.push(t);
        return t;
    }

    function txn(uint8 weight, Uint8Model.Place memory p, Uint8Model.Transition memory t) internal {
        transitions[t.offset].delta[p.offset] = 0-int8(weight);
    }

    function txn(uint8 weight, Uint8Model.Transition memory t, Uint8Model.Place memory p) internal {
        transitions[t.offset].delta[p.offset] = int8(weight);
    }

    function guard(uint8 weight, Uint8Model.Place memory p, Uint8Model.Transition memory t) internal {
       transitions[t.offset].guard[p.offset] = 0-int8(weight);
    }

}