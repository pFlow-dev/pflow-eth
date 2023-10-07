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

    string public constant metamodelUri = "ipns://pflow.eth?contract=KonamiCode";

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
        cell(2, 2, Uint8Model.Position(1,1)); // TwoUps
        cell(2, 2, Uint8Model.Position(1,1)); // TwoDowns
        cell(2, 2, Uint8Model.Position(1,1)); // TwoLefts
        cell(0, 2, Uint8Model.Position(1,1)); // TwoRights
        cell(0, 1, Uint8Model.Position(1,1)); // ThenRight
        cell(0, 1, Uint8Model.Position(1,1)); // ThenSelect
        cell(0, 1, Uint8Model.Position(1,1)); // ThenStart
        cell(0, 1, Uint8Model.Position(1,1)); // ThenA
    }

    function _action(KonamiCodeModel.Actions action, uint8 x, uint8 y) internal {
        fn(uint8(KonamiCodeModel.Properties.SIZE), uint8(action), uint8(KonamiCodeModel.Roles.PLAYER), Uint8Model.Position(x,y));
    }

    function _actions() internal {
        for (uint8 i = 0; i < uint8(KonamiCodeModel.Properties.SIZE); i++) {
            _action(KonamiCodeModel.Actions(i), 1, 1); // FIXME: set proper coordinates
        }
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

    address internal owner;

    bool internal paused = false;
    uint8 internal sequence = 0;
    uint256 public session = 0;

    int8[] public state = new int8[](8);

    bytes32 public constant PLAYER = keccak256("PLAYER");
    bytes32 public constant UNLOCK = keccak256("UNLOCK");

    constructor() {
        owner = msg.sender;
        _setupRole(DEFAULT_ADMIN_ROLE, owner);
        resetSession(KonamiCodeModel.Roles.HALT);
    }

    modifier startSession() {
        sequence = 0;
        session++;
        _;
        Uint8Model.Place[] memory pl = places;
        for (uint8 i = 0; i < uint8(KonamiCodeModel.Properties.SIZE); i++) {
            _init(pl[i]);
        }
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

    function resetSession(KonamiCodeModel.Roles role) internal startSession {
        emit Uint8Model.SignalEvent(session, sequence, uint8(KonamiCodeModel.Actions.HALT), uint8(role), block.timestamp);
    }

    modifier inhibitAction(Uint8Model.Transition memory t) {
        require(!paused, "Game is paused.");
        require(!actionInhibited(t), "Action is inhibited.");
        _;
    }

    function transform(uint8 i, Uint8Model.Transition memory t)  internal override {
        if (t.delta[i] != 0) {
            state[i] = state[i] + t.delta[i];
            // require(state[i] >= 0, 'invalid state');
            if (state[i] < 0) {
                resetLast(t.offset);
            }
        }
    }

    function resetLast(uint8 i) internal {
        resetSession(KonamiCodeModel.Roles.HALT);
        if (i == uint8(KonamiCodeModel.Actions.Up)) {
            signal(KonamiCodeModel.Actions.Up); // re-apply first move after reset
        }
    }

    function _transform(Uint8Model.Transition memory t) internal {
        Uint8Model.Place[] memory p = places;
        if (actionInhibited(t)) {
            resetLast(t.offset);
        }
        for (uint8 i = 0; i < uint8(KonamiCodeModel.Properties.SIZE); i++) {
            transform(i, t);
            if (p[i].capacity != 0) {
                if (state[i] > p[i].capacity) {
                    resetLast(t.offset);
                }
            }
        }
    }

    function actionInhibited(Uint8Model.Transition memory t) internal view returns (bool) {
        for (uint8 i = 0; i < uint8(KonamiCodeModel.Properties.SIZE); i++) {
            if (_isInhibited(i, t)) {
                return true;
            }
        }
        return false;
    }

    function _isInhibited(uint8 i, Uint8Model.Transition memory t) internal view returns (bool) {
        if (t.guard[i] != 0) {
            return state[i] + t.guard[i] >= 0;
        }
        return false;
    }

    modifier autoIncrement() {
        _;
        sequence++;
    }

    function signal(KonamiCodeModel.Actions action) internal autoIncrement {
        // REVIEW: do we have calldata access here?
        uint8 txnId = uint8(action);
        Uint8Model.Transition memory t = transitions[txnId];
        _transform(t);
        emit Uint8Model.SignalEvent(session, sequence, uint8(action), uint8(KonamiCodeModel.Roles.PLAYER), block.timestamp);
    }

    function Up() external {
        signal(KonamiCodeModel.Actions.Up);
    }

    function Down() external {
        signal(KonamiCodeModel.Actions.Down);
    }

    function Left() external {
        signal(KonamiCodeModel.Actions.Left);
    }

    function Right() external {
        signal(KonamiCodeModel.Actions.Right);
    }

    function Select() external {
        signal(KonamiCodeModel.Actions.Select);
    }

    function Start() external {
        signal(KonamiCodeModel.Actions.Start);
    }

    function B() external {
        signal(KonamiCodeModel.Actions.B);
    }

    function A() external {
        signal(KonamiCodeModel.Actions.A);
    }

}
