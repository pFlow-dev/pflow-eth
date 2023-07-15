// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;
import "@openzeppelin/contracts/access/AccessControl.sol";
import "./Metamodel.sol";

library KonamiCodeModel {

    enum Roles{ PLAYER, HALT }

    enum Properties {
        TwoUps,
        TwoDowns,
        TwoLefts,
        TwoRights,
        ThenRight,
        ThenSelect,
        ThenStart,
        ThenA,
        SIZE
    }

    enum Actions {
        Up,
        Down,
        Left,
        Right,
        Select,
        Start,
        B,
        A,
        HALT
    }
}

abstract contract KonamiCodeMetaModel is MetamodelUint8 {

    // add a new action
    function _txn(KonamiCodeModel.Properties prop, KonamiCodeModel.Actions action, uint8 weight ) internal {
        txn(weight, places[uint8(prop)], transitions[uint8(action)]);
    }

    function _txn(KonamiCodeModel.Actions action, KonamiCodeModel.Properties prop, uint8 weight) internal {
        txn(weight, transitions[uint8(action)], places[uint8(prop)]);
    }

    function _guard(KonamiCodeModel.Properties prop, KonamiCodeModel.Actions action, uint8 weight) internal {
        guard(weight, places[uint8(prop)], transitions[uint8(action)]);
    }

    function _props() internal {
        cell(2, 2); // TwoUps
        cell(2, 2); // TwoDowns
        cell(2, 2); // TwoLefts
        cell(0, 1); // TwoRights
        cell(0, 1); // ThenRight
        cell(0, 1); // ThenSelect
        cell(0, 1); // ThenStart
        cell(0, 1); // ThenA
    }

    function _action(KonamiCodeModel.Actions action) internal {
        fn(uint8(KonamiCodeModel.Properties.SIZE), uint8(action), uint8(KonamiCodeModel.Roles.PLAYER));
    }

    function _actions() internal {
        _action(KonamiCodeModel.Actions.Up);
        _action(KonamiCodeModel.Actions.Down);
        _action(KonamiCodeModel.Actions.Left);
        _action(KonamiCodeModel.Actions.Right);
        _action(KonamiCodeModel.Actions.Select);
        _action(KonamiCodeModel.Actions.Start);
        _action(KonamiCodeModel.Actions.B);
        _action(KonamiCodeModel.Actions.A);
    }

    constructor() {
        _props();
        _actions();

        _txn(KonamiCodeModel.Properties.TwoUps, KonamiCodeModel.Actions.Up, 1);
        _guard(KonamiCodeModel.Properties.TwoUps, KonamiCodeModel.Actions.Down, 1);

        _txn(KonamiCodeModel.Properties.TwoDowns, KonamiCodeModel.Actions.Down, 1);
        _guard(KonamiCodeModel.Properties.TwoDowns, KonamiCodeModel.Actions.Left, 1);

        _txn(KonamiCodeModel.Properties.TwoLefts, KonamiCodeModel.Actions.Left, 1);

        _txn(KonamiCodeModel.Properties.ThenRight, KonamiCodeModel.Actions.Right, 1);
        _txn(KonamiCodeModel.Actions.Left, KonamiCodeModel.Properties.ThenRight, 1);

        _txn(KonamiCodeModel.Properties.TwoRights, KonamiCodeModel.Actions.B, 2);
        _txn(KonamiCodeModel.Actions.Right, KonamiCodeModel.Properties.TwoRights, 1);

        _txn(KonamiCodeModel.Properties.ThenA, KonamiCodeModel.Actions.A, 1);
        _txn(KonamiCodeModel.Actions.B, KonamiCodeModel.Properties.ThenA, 1);

        _txn(KonamiCodeModel.Properties.ThenSelect, KonamiCodeModel.Actions.Select, 1);
        _txn(KonamiCodeModel.Actions.A, KonamiCodeModel.Properties.ThenSelect, 1);

        _txn(KonamiCodeModel.Properties.ThenStart, KonamiCodeModel.Actions.Start, 1);
        _txn(KonamiCodeModel.Actions.Select, KonamiCodeModel.Properties.ThenStart, 1);

    }

}

/// @custom:security-contact security@stackdump.com
contract KonamiCode is KonamiCodeMetaModel, AccessControl {

    string public constant metamodelUri = "ipns://pflow.eth?contract=KonamiCode";

    address internal owner;

    bool internal paused = false;
    uint8 internal sequence = 0;
    uint256 internal session = 0;

    int8[] public state = new int8[](8);

    bytes32 public constant PLAYER = keccak256("PLAYER");
    bytes32 public constant UNLOCK = keccak256("UNLOCK");

    constructor() {
        owner = msg.sender;
        _setupRole(DEFAULT_ADMIN_ROLE, owner);
        resetGame(KonamiCodeModel.Roles.HALT);
    }

    modifier startGame() {
        sequence = 0;
        session++;
        _;
        Uint8Model.Place[] memory p = places;
        _init(p[0]); // TwoUps
        _init(p[1]); // TwoDowns
        _init(p[2]); // TwoLefts
        _init(p[3]); // TwoRights
        _init(p[4]); // ThenRight
        _init(p[5]); // ThenSelect
        _init(p[6]); // ThenStart
        _init(p[7]); // ThenA
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

    function resetGame(KonamiCodeModel.Roles role) internal startGame {
        emit Uint8Model.Action(session, sequence, uint8(KonamiCodeModel.Actions.HALT), uint8(role), block.timestamp);
    }

    function transform(uint8 i, Uint8Model.Transition memory t)  internal override {
        // TODO: compute state per player
        if (t.delta[i] != 0) {
            state[i] = state[i] + t.delta[i];
            require(state[i] >= 0, "Invalid state");
        }
    }

    function send(KonamiCodeModel.Actions action) internal {
        emit Uint8Model.Action(session, sequence, uint8(action), uint8(KonamiCodeModel.Roles.PLAYER), block.timestamp);
    }

    function getRole() public pure returns (KonamiCodeModel.Roles) {
        return KonamiCodeModel.Roles.PLAYER;
    }

    function Up() public {
        send(KonamiCodeModel.Actions.Up);
    }

    function Down() public {
        send(KonamiCodeModel.Actions.Down);
    }

    function Left() public {
        send(KonamiCodeModel.Actions.Left);
    }

    function Right() public {
        send(KonamiCodeModel.Actions.Right);
    }

    function Select() public {
        send(KonamiCodeModel.Actions.Select);
    }

    function Start() public {
        send(KonamiCodeModel.Actions.Start);
    }

    function B() public {
        send(KonamiCodeModel.Actions.B);
    }

    function A() public {
        send(KonamiCodeModel.Actions.A);
    }

}
