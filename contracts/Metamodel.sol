// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;
import "@openzeppelin/contracts/access/AccessControl.sol";
// import "hardhat/console.sol";

library Uint8Model {

    event Action(uint256 indexed gameId, uint8 indexed seq, uint8 txnId, uint8 role, uint256 when);

    struct PetriNet {
        Place[] places;
        Transition[] transitions;
    }

    struct Transition {
        uint8 offset;
        uint8 role;
        int8[] delta;
        // int8[] guard; // REVIEW: we don't make use of guards in tic-tac-toe example
    }

    struct Place {
        uint8 offset;
        int8 initial;
        // int8 capacity; // REVIEW: capacity check not used
    }

    // REVIEW: should we use this function to unpack declarations?
    // function model(Uint8ModelFactory factory) external returns (PetriNet memory) {
    //     return factory.declaration();
    // }

}

interface Uint8ModelFactory {
    function declaration() external returns (Uint8Model.PetriNet memory);
}

abstract contract MetamodelUint8  {

    Uint8Model.Place[] internal places;
    Uint8Model.Transition[] internal transitions;

    function cell(int8 initial) internal returns (Uint8Model.Place memory) {
        Uint8Model.Place memory p =  Uint8Model.Place(uint8(places.length), initial);
        places.push(p);
        return p;
    }

    function fn(uint8 vectorSize, uint8 action, uint8 role) internal returns (Uint8Model.Transition memory) {
        require(uint8(transitions.length) == action, "Transition offset must match Actions enum");
        Uint8Model.Transition memory t = Uint8Model.Transition(action, role, new int8[](vectorSize));
        transitions.push(t);
        return t;
    }

    function txn(uint8 weight, Uint8Model.Place memory p, Uint8Model.Transition memory t) internal {
        transitions[t.offset].delta[p.offset] = 0-int8(weight);
    }

    function txn(uint8 weight, Uint8Model.Transition memory t, Uint8Model.Place memory p) internal {
        transitions[t.offset].delta[p.offset] = int8(weight);
    }

    // function guard(uint8 weight, Uint8Model.Place memory p, Uint8Model.Transition memory t) internal {
    //    transitions[t.offset].guard[p.offset] = 0-int8(weight);
    // }

}

abstract contract TicTacToeModel is MetamodelUint8, Uint8ModelFactory {

    enum Roles{ X, O, HALT }

    enum Properties {
        _00, _01, _02,
        _10, _11, _12,
        _20, _21, _22,
        SIZE
    }

    enum Actions {
        // x moves
        X00, X01, X02,
        X10, X11, X12,
        X20, X21, X22,
        // O moves
        O00, O01, O02,
        O10, O11, O12,
        O20, O21, O22,
        HALT
    }

    // add a new action
    function _action(Properties prop, Actions action, Roles role) internal {
        require(role < Roles.HALT, "Invalid role");
        require(action < Actions.HALT, "Invalid action");
        txn(1, places[uint8(prop)], fn(uint8(Properties.SIZE), uint8(action), uint8(role)));
    }

    // add a property to the model
    function _prop() internal {
        cell(1);
    }

    constructor() {
        _prop(); // _00
        _prop(); // _01
        _prop(); // _02

        _prop(); // _10
        _prop(); // _11
        _prop(); // _12

        _prop(); // _20
        _prop(); // _21
        _prop(); // _22

        _action(Properties._00, Actions.X00, Roles.X);
        _action(Properties._01, Actions.X01, Roles.X);
        _action(Properties._02, Actions.X02, Roles.X);

        _action(Properties._10, Actions.X10, Roles.X);
        _action(Properties._11, Actions.X11, Roles.X);
        _action(Properties._12, Actions.X12, Roles.X);

        _action(Properties._20, Actions.X20, Roles.X);
        _action(Properties._21, Actions.X21, Roles.X);
        _action(Properties._22, Actions.X22, Roles.X);

        _action(Properties._00, Actions.O00, Roles.O);
        _action(Properties._01, Actions.O01, Roles.O);
        _action(Properties._02, Actions.O02, Roles.O);

        _action(Properties._10, Actions.O10, Roles.O);
        _action(Properties._11, Actions.O11, Roles.O);
        _action(Properties._12, Actions.O12, Roles.O);

        _action(Properties._20, Actions.O20, Roles.O);
        _action(Properties._21, Actions.O21, Roles.O);
        _action(Properties._22, Actions.O22, Roles.O);
    }

    function declaration() public view returns (Uint8Model.PetriNet memory) {
        return Uint8Model.PetriNet(places, transitions);
    }

}

contract TicTacToe is AccessControl, TicTacToeModel {
    address internal owner;
    bool internal paused = false;

    uint256 internal gameId = 0;
    uint8 internal sequence = 0;

    int8[] public state = new int8[](9);

    bytes32 public constant PLAYER_X = keccak256("PLAYER_X");
    bytes32 public constant PLAYER_O = keccak256("PLAYER_O");

    constructor(address p0, address p1) {
        owner = msg.sender;
        require(p0 != p1, "Players must not have the same address.");
        _setupRole(DEFAULT_ADMIN_ROLE, owner);
        _grantRole(PLAYER_X, p0);
        _grantRole(PLAYER_O, p1);
        resetGame(TicTacToeModel.Roles.HALT);
    }

    function _init(Uint8Model.Place memory p) internal {
        if (state[p.offset] != p.initial) {
            state[p.offset] = p.initial;
        }
    }

    // no news is good news - revert if game is paused
    function testIsGameOpen() public view {
        require(!paused, "Game is paused.");
    }

    function pause() public onlyRole(DEFAULT_ADMIN_ROLE) {
        paused = true;
    }

    function unpause() public onlyRole(DEFAULT_ADMIN_ROLE) {
        paused = false;
    }

    modifier startGame() {
        sequence = 0;
        gameId++;
        _;
        Uint8Model.Place[] memory p = places;
        _init(p[0]);
        _init(p[1]);
        _init(p[2]);
        _init(p[3]);
        _init(p[4]);
        _init(p[5]);
        _init(p[6]);
        _init(p[7]);
        _init(p[8]);
    }

    function resetGame(TicTacToeModel.Roles role) internal startGame {
        emit Uint8Model.Action(gameId, sequence, uint8(TicTacToeModel.Actions.HALT), uint8(role), block.timestamp);
    }

    function transform(uint8 i, Uint8Model.Transition memory t)  internal {
        if (t.delta[i] != 0) {
            state[i] = state[i] + t.delta[i];
            require(state[i] >= 0, "Invalid state");
        }
    }

    // no news is good news - revert if it is not the caller's turn
    function testIsMyTurn() public view  {
        testIsGameOpen();
        require(msg.sender != owner, "Gameplay from contract owner is forbidden.");
        if (sequence % 2 == 0) {
            require(getRole() == Roles.X, "X's turn");
        } else {
            require(getRole() == Roles.O, "O's turn");
        }
    }

    modifier takeTurns() {
        testIsMyTurn();
        _;
        sequence++;
    }

    function move(TicTacToeModel.Actions action) internal takeTurns {
        uint8 txnId = uint8(action);
        Uint8Model.Transition memory t = transitions[txnId];
        assert(txnId == t.offset);
        transform(0, t);
        transform(1, t);
        transform(2, t);
        transform(3, t);
        transform(4, t);
        transform(5, t);
        transform(6, t);
        transform(7, t);
        transform(8, t);
        emit Uint8Model.Action(gameId, sequence, txnId, t.role, block.timestamp);
    }

    function reset() public {
        resetGame(getRole());
    }

    function getRole() public view returns (Roles) {
        if (hasRole(PLAYER_X, msg.sender)) {
            return TicTacToeModel.Roles.X;
        } else if (hasRole(PLAYER_O, msg.sender)) {
            return TicTacToeModel.Roles.O;
        } else {
            revert("Unexpected caller");
        }
    }

    function X00() public {
        move(TicTacToeModel.Actions.X00);
    }

    function X01() public {
        move(TicTacToeModel.Actions.X01);
    }

    function X02() public {
        move(TicTacToeModel.Actions.X02);
    }

    function X10() public {
        move(TicTacToeModel.Actions.X10);
    }

    function X11() public {
        move(TicTacToeModel.Actions.X11);
    }

    function X12() public {
        move(TicTacToeModel.Actions.X12);
    }

    function X20() public {
        move(TicTacToeModel.Actions.X20);
    }

    function X21() public {
        move(TicTacToeModel.Actions.X21);
    }

    function X22() public {
        move(TicTacToeModel.Actions.X22);
    }

    function O00() public {
        move(TicTacToeModel.Actions.O00);
    }

    function O01() public {
        move(TicTacToeModel.Actions.O01);
    }

    function O02() public {
        move(TicTacToeModel.Actions.O02);
    }

    function O10() public {
        move(TicTacToeModel.Actions.O10);
    }

    function O11() public {
        move(TicTacToeModel.Actions.O11);
    }

    function O12() public {
        move(TicTacToeModel.Actions.O12);
    }

    function O20() public {
        move(TicTacToeModel.Actions.O20);
    }

    function O21() public {
        move(TicTacToeModel.Actions.O21);
    }

    function O22() public {
        move(TicTacToeModel.Actions.O22);
    }
}