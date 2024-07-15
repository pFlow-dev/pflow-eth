// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

library Model {

    event SignaledEvent(
        uint8 indexed role,
        uint8 indexed actionId,
        uint256 indexed scalar,
        uint256 sequence
    );

    struct PetriNet {
        Place[] places;
        Transition[] transitions;
    }

    struct Position {
        uint8 x;
        uint8 y;
    }

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

}

interface ModelInterface {
    function model() external returns (Model.PetriNet memory);
    function signal(uint8 action, uint256 scalar) external;
    function signalMany(uint8[] calldata actions, uint256[] calldata scalars) external;
}

abstract contract PflowDSL {
    Model.Place[] internal places;
    Model.Transition[] internal transitions;

    function cell(string memory label, uint256 initial, uint256 capacity, Model.Position memory position) internal returns (Model.Place memory) {
        Model.Place memory p = Model.Place(label, uint8(places.length), position, initial, capacity);
        places.push(p);
        return p;
    }

    function func(string memory label, uint8 vectorSize, uint8 action, uint8 role, Model.Position memory position) internal returns (Model.Transition memory) {
        require(uint8(transitions.length) == action, "transaction => enum miss match");
        Model.Transition memory t = Model.Transition(label, action, position, role, new int256[](vectorSize), new int256[](vectorSize));
        transitions.push(t);
        return t;
    }

    function arrow(int256 weight, Model.Place memory p, Model.Transition memory t) internal {
        require(weight > 0, "weight must be > 0");
        transitions[t.offset].delta[p.offset] = 0 - weight;
    }

    function arrow(int256 weight, Model.Transition memory t, Model.Place memory p) internal {
        require(weight > 0, "weight must be > 0");
        transitions[t.offset].delta[p.offset] = weight;
    }

    // inhibit transition after threshold weight is reached
    function guard(int256 weight, Model.Place memory p, Model.Transition memory t) internal {
        require(weight > 0, "weight must be > 0");
        transitions[t.offset].guard[p.offset] = 0 - weight;
    }

    // inhibit transition until threshold weight is reached
    function guard(int256 weight, Model.Transition memory t, Model.Place memory p) internal {
        require(weight > 0, "weight must be > 0");
        transitions[t.offset].guard[p.offset] = weight;
    }
}

abstract contract Metamodel is PflowDSL, ModelInterface {

    // sequence is a monotonically increasing counter for each signal
    uint256 public sequence = 0;

    // transform is a hook for derived contracts to implement state transitions
    function transform(uint8 i, Model.Transition memory t, uint256 scalar) internal virtual;

    // isInhibited is a hook for derived contracts to implement transition guards
    function isInhibited(Model.Transition memory t) internal view virtual returns (bool);

    // hasPermission implements an ACL for transitions based on user roles
    function hasPermission(Model.Transition memory t) internal view virtual returns (bool);

    function _signal(uint8 action, uint256 scalar) internal {
        Model.Transition memory t = transitions[action];
        require(!isInhibited(t), 'inhibited');
        assert(action == t.offset);
        for (uint8 i = 0; i < uint8(places.length); i++) {
            transform(i, t, scalar);
        }
        sequence++;
        emit Model.SignaledEvent(t.role, action, scalar, sequence);
    }

    function signal(uint8 action, uint256 scalar) external {
        _signal(action, scalar);
    }

    function signalMany(uint8[] calldata actions, uint256[] calldata scalars) external {
        require(actions.length == scalars.length, "ModelRegistry: invalid input");
        for (uint256 i = 0; i < actions.length; i++) {
            _signal(actions[i], scalars[i]);
        }
    }

    // model returns the model in a format suited for execution
    function model() external view returns (Model.PetriNet memory) {
        return Model.PetriNet(places, transitions);
    }

}


abstract contract MyModel is Metamodel {

    enum Roles {DEFAULT, X, O, HALT}
    enum Properties {p00, p01, p02, p10, p11, p12, p20, p21, p22, next, SIZE}
    enum Actions {X00, X01, X02, X10, X11, X12, X20, X21, X22, O00, O01, O02, O10, O11, O12, O20, O21, O22, HALT}

    int256[] public state = new int256[](uint8(Properties.SIZE));

    constructor() {
        cell("p00", 1, 1, Model.Position(1, 1));
        cell("p01", 1, 1, Model.Position(2, 1));
        cell("p02", 1, 1, Model.Position(3, 1));
        cell("p10", 1, 1, Model.Position(1, 3));
        cell("p11", 1, 1, Model.Position(2, 3));
        cell("p12", 1, 1, Model.Position(3, 3));
        cell("p20", 1, 1, Model.Position(1, 4));
        cell("p21", 1, 1, Model.Position(2, 4));
        cell("p22", 1, 1, Model.Position(3, 4));
        cell("next", 0, 1, Model.Position(6, 8));

        func("X00", uint8(Properties.SIZE), uint8(0), uint8(Roles.X), Model.Position(5, 1));
        func("X01", uint8(Properties.SIZE), uint8(1), uint8(Roles.X), Model.Position(6, 1));
        func("X02", uint8(Properties.SIZE), uint8(2), uint8(Roles.X), Model.Position(7, 1));
        func("X10", uint8(Properties.SIZE), uint8(3), uint8(Roles.X), Model.Position(5, 3));
        func("X11", uint8(Properties.SIZE), uint8(4), uint8(Roles.X), Model.Position(6, 3));
        func("X12", uint8(Properties.SIZE), uint8(5), uint8(Roles.X), Model.Position(7, 3));
        func("X20", uint8(Properties.SIZE), uint8(6), uint8(Roles.X), Model.Position(5, 4));
        func("X21", uint8(Properties.SIZE), uint8(7), uint8(Roles.X), Model.Position(6, 4));
        func("X22", uint8(Properties.SIZE), uint8(8), uint8(Roles.X), Model.Position(7, 4));
        func("O00", uint8(Properties.SIZE), uint8(9), uint8(Roles.O), Model.Position(1, 7));
        func("O01", uint8(Properties.SIZE), uint8(10), uint8(Roles.O), Model.Position(2, 7));
        func("O02", uint8(Properties.SIZE), uint8(11), uint8(Roles.O), Model.Position(3, 7));
        func("O10", uint8(Properties.SIZE), uint8(12), uint8(Roles.O), Model.Position(1, 8));
        func("O11", uint8(Properties.SIZE), uint8(13), uint8(Roles.O), Model.Position(2, 8));
        func("O12", uint8(Properties.SIZE), uint8(14), uint8(Roles.O), Model.Position(3, 8));
        func("O20", uint8(Properties.SIZE), uint8(15), uint8(Roles.O), Model.Position(1, 10));
        func("O21", uint8(Properties.SIZE), uint8(16), uint8(Roles.O), Model.Position(2, 10));
        func("O22", uint8(Properties.SIZE), uint8(17), uint8(Roles.O), Model.Position(3, 10));

        arrow(1, places[0], transitions[0]);
        arrow(1, transitions[0], places[9]);
        arrow(1, places[1], transitions[1]);
        arrow(1, transitions[1], places[9]);
        arrow(1, places[2], transitions[2]);
        arrow(1, transitions[2], places[9]);
        arrow(1, places[3], transitions[3]);
        arrow(1, transitions[3], places[9]);
        arrow(1, places[4], transitions[4]);
        arrow(1, transitions[4], places[9]);
        arrow(1, places[5], transitions[5]);
        arrow(1, transitions[5], places[9]);
        arrow(1, places[6], transitions[6]);
        arrow(1, transitions[6], places[9]);
        arrow(1, places[7], transitions[7]);
        arrow(1, transitions[7], places[9]);
        arrow(1, places[8], transitions[8]);
        arrow(1, transitions[8], places[9]);
        arrow(1, places[0], transitions[9]);
        arrow(1, places[9], transitions[9]);
        arrow(1, places[1], transitions[10]);
        arrow(1, places[9], transitions[10]);
        arrow(1, places[2], transitions[11]);
        arrow(1, places[9], transitions[11]);
        arrow(1, places[3], transitions[12]);
        arrow(1, places[9], transitions[12]);
        arrow(1, places[4], transitions[13]);
        arrow(1, places[9], transitions[13]);
        arrow(1, places[5], transitions[14]);
        arrow(1, places[9], transitions[14]);
        arrow(1, places[6], transitions[15]);
        arrow(1, places[9], transitions[15]);
        arrow(1, places[7], transitions[16]);
        arrow(1, places[9], transitions[16]);
        arrow(1, places[8], transitions[17]);
        arrow(1, places[9], transitions[17]);

        for (uint8 i = 0; i < uint8(Properties.SIZE); i++) {
            state[i] = int256(places[i].initial);
        }
    }
}

contract Tictactoe is MyModel {

    function isInhibited(Model.Transition memory t) internal view override returns (bool) {
        for (uint8 i = 0; i < uint8(Properties.SIZE); i++) {
            if (t.guard[i] != 0) {
                if (t.guard[i] < 0) {
                    // inhibit unless condition is met
                    if ((state[i] + t.guard[i]) > 0) {
                        return true;
                    }
                } else {
                    // inhibit until condition is met
                    if ((state[i] - t.guard[i]) < 0) {
                        return true;
                    }
                }
            }
        }
        return false;
    }

    function hasPermission(Model.Transition memory t) internal view override returns (bool) {
        return t.role < uint8(Roles.HALT);
    }

    function transform(uint8 i, Model.Transition memory t, uint256 scalar) internal override {
        require(scalar > 0, "invalid scalar");
        if (t.delta[i] != 0) {
            state[i] = state[i] + t.delta[i] * int256(scalar);
            require(state[i] >= 0, "underflow");
            if (places[i].capacity > 0) {
                require(state[i] <= int256(places[i].capacity), "overflow");
            }
        }
    }
}