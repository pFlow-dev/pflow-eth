// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;

library Declaration {

    struct place {
        string label;
        uint8 x;
        uint8 y;
        uint256 initial;
        uint256 capacity;
    }

    struct transition {
        string label;
        uint8 x;
        uint8 y;
        uint8 role;
    }

    struct arc {
        string source;
        string target;
        uint256 weight;
        bool consume;
        bool produce;
        bool inhibit;
        bool read;
    }

    struct PetriNet {
        place[] places;
        transition[] transitions;
        arc[] arcs;
    }

}

library Model {

    event SignalEvent(
        uint8 indexed role,
        uint8 indexed actionId,
        uint256 indexed scalar
    );

    struct PetriNet {
        Place[] places;
        Transition[] transitions;
        Arc[] arcs;
    }

    struct Position {
        uint8 x;
        uint8 y;
    }

    // REVIEW: should we replace int256[] w/ the more compact Vector
    // struct Scalar {
    //     int256 value;
    //     uint8 offset;
    // }
    //
    // struct Vector {
    //     Scalar[] values;
    // }

    struct Transition {
        string label;
        uint8 offset;
        Position position;
        uint8 role;
        int256[] delta;
        int256[] guard;
    }

    struct Place {
        string label;
        uint8 offset;
        Position position;
        uint256 initial;
        uint256 capacity;
    }

    enum NodeKind {
        PLACE,
        TRANSITION
    }

    struct Node {
        string label;
        uint8 offset;
        NodeKind kind;
    }

    struct Arc {
        uint256 weight;
        Node source;
        Node target;
        bool inhibitor;
        bool read;
    }

}

interface ModelInterface {
    function model() external returns (Model.PetriNet memory);

    function declaration() external returns (Declaration.PetriNet memory);

    function signal(uint8 action, uint256 scalar) external;

    function signalMany(uint8[] calldata actions, uint256[] calldata scalars) external;
}

abstract contract Metamodel is ModelInterface {
    uint256 public sequence = 0;

    Model.Place[] internal places;
    Model.Transition[] internal transitions;
    Model.Arc[] internal arcs;

    // transform is a hook for derived contracts to implement state transitions
    function transform(uint8 i, Model.Transition memory t, uint256 scalar) internal virtual;

    // signal executes a transition and broadcasts the event
    function signal(uint8 action, uint256 scalar) external {
        Model.Transition memory t = transitions[action];
        assert(action == t.offset);
        for (uint8 i = 0; i < uint8(places.length); i++) {
            transform(i, t, scalar);
        }
        sequence++;
        emit Model.SignalEvent(t.role, action, scalar );
    }

    // send multiple signals
    function signalMany(uint8[] calldata actions, uint256[] calldata scalars) external {
        require(actions.length == scalars.length, "ModelRegistry: invalid input");
        for (uint256 i = 0; i < actions.length; i++) {
            this.signal(actions[i], scalars[i]);
        }
    }

    // model returns an indexed model of the PetriNet
    function model() external view returns (Model.PetriNet memory) {
        return Model.PetriNet(places, transitions, arcs);
    }

    // declaration returns a minimal model declaration suitable for serialization
    function declaration() external view returns (Declaration.PetriNet memory) {
        Declaration.place[] memory p = new Declaration.place[](places.length);
        for (uint8 i = 0; i < uint8(places.length); i++) {
            p[i] = Declaration.place(places[i].label, places[i].position.x, places[i].position.y, places[i].initial, places[i].capacity);
        }
        Declaration.transition[] memory t = new Declaration.transition[](transitions.length);
        for (uint8 i = 0; i < uint8(transitions.length); i++) {
            t[i] = Declaration.transition(transitions[i].label, transitions[i].position.x, transitions[i].position.y, transitions[i].role);
        }
        Declaration.arc[] memory a = new Declaration.arc[](arcs.length);
        for (uint8 i = 0; i < uint8(arcs.length); i++) {
            assert(arcs[i].source.kind != arcs[i].target.kind);
            a[i] = Declaration.arc(
                arcs[i].source.label,
                arcs[i].target.label,
                arcs[i].weight,
                arcs[i].source.kind == Model.NodeKind.PLACE, // consume
                arcs[i].target.kind == Model.NodeKind.PLACE, // produce
                arcs[i].inhibitor,
                arcs[i].read
            );
        }
        return Declaration.PetriNet(p, t, a);
    }

    function placeNode(string memory label, uint8 offset) internal pure returns (Model.Node memory) {
        return Model.Node(label, offset, Model.NodeKind.PLACE);
    }

    function transitionNode(string memory label, uint8 offset) internal pure returns (Model.Node memory) {
        return Model.Node(label, offset, Model.NodeKind.TRANSITION);
    }

    function cell(string memory label, uint256 initial, uint256 capacity, Model.Position memory position) internal returns (Model.Place memory) {
        Model.Place memory p = Model.Place(label, uint8(places.length), position, initial, capacity);
        places.push(p);
        return p;
    }

    function func(string memory label, uint8 vectorSize, uint8 action, uint8 role, Model.Position memory position) internal returns (Model.Transition memory) {
        require(uint8(transitions.length) == action, "Transition offset must match Actions enum");
        Model.Transition memory t = Model.Transition(label, action, position, role, new int256[](vectorSize), new int256[](vectorSize));
        transitions.push(t);
        return t;
    }

    function arrow(uint8 weight, Model.Place memory p, Model.Transition memory t) internal {
        arcs.push(Model.Arc(weight, placeNode(p.label, p.offset), transitionNode(t.label, t.offset), false, false));
        transitions[t.offset].delta[p.offset] = 0 - int8(weight);
    }

    function arrow(uint8 weight, Model.Transition memory t, Model.Place memory p) internal {
        arcs.push(Model.Arc(weight, transitionNode(t.label, t.offset), placeNode(p.label, p.offset), false, false));
        transitions[t.offset].delta[p.offset] = int8(weight);
    }

    // inhibit transition after threshold weight is reached
    function guard(uint8 weight, Model.Place memory p, Model.Transition memory t) internal {
        arcs.push(Model.Arc(weight, placeNode(p.label, p.offset), transitionNode(t.label, t.offset), true, false));
        transitions[t.offset].guard[p.offset] = 0 - int8(weight);
    }

    // inhibit transition until threshold weight is reached
    function guard(uint8 weight, Model.Transition memory t, Model.Place memory p) internal {
        arcs.push(Model.Arc(weight, transitionNode(t.label, t.offset), placeNode(p.label, p.offset), true, true));
        transitions[t.offset].guard[p.offset] = int8(weight);
    }

}
