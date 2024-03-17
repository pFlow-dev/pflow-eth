// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;
import "@openzeppelin/contracts/access/AccessControl.sol";
import "../../contracts/Metamodel.sol";


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

abstract contract KonamiCodeMetaModel is Metamodel {

    function _txn(KonamiCodeModel.Properties prop, KonamiCodeModel.Actions action, uint8 weight ) internal {
        arrow(weight, places[uint8(prop)], transitions[uint8(action)]);
    }

    function _txn(KonamiCodeModel.Actions action, KonamiCodeModel.Properties prop, uint8 weight) internal {
        arrow(weight, transitions[uint8(action)], places[uint8(prop)]);
    }

    function _guard(KonamiCodeModel.Properties prop, KonamiCodeModel.Actions action, uint8 weight) internal {
        guard(weight, places[uint8(prop)], transitions[uint8(action)]);
    }

    function _props() internal {
        cell("TwoUps", 2, 2, Model.Position(1,1)); // TwoUps
        cell("TwoDowns", 2, 2, Model.Position(1,1)); // TwoDowns
        cell("TwoLefts", 2, 2, Model.Position(1,1)); // TwoLefts
        cell("TwoRights", 0, 2, Model.Position(1,1)); // TwoRights
        cell("ThenRight", 0, 1, Model.Position(1,1)); // ThenRight
        cell("ThenSelect", 0, 1, Model.Position(1,1)); // ThenSelect
        cell("ThenStart", 0, 1, Model.Position(1,1)); // ThenStart
        cell("ThenA", 0, 1, Model.Position(1,1)); // ThenA
    }

    function _action(string memory label, KonamiCodeModel.Actions action, uint8 x, uint8 y) internal {
        func(label, uint8(KonamiCodeModel.Properties.SIZE), uint8(action), uint8(KonamiCodeModel.Roles.PLAYER), Model.Position(x,y));
    }

    function _actions() internal {
        _action("Up", KonamiCodeModel.Actions.Up, 5, 1);
        _action("Down", KonamiCodeModel.Actions.Down, 6, 1);
        _action("Left", KonamiCodeModel.Actions.Left, 7, 1);
        _action("Right", KonamiCodeModel.Actions.Right, 8, 1);
        _action("Select", KonamiCodeModel.Actions.Select, 9, 1);
        _action("Start", KonamiCodeModel.Actions.Start, 10, 1);
        _action("B", KonamiCodeModel.Actions.B, 11, 1);
        _action("A", KonamiCodeModel.Actions.A, 12, 1);
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
    uint256 public session = 0;

    int256[] public state = new int256[](8);

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
        Model.Place[] memory pl = places;
        for (uint8 i = 0; i < uint8(KonamiCodeModel.Properties.SIZE); i++) {
            _init(pl[i]);
        }
    }

    function _init(Model.Place memory p) internal {
        if (state[p.offset] != int256(p.initial)) {
            state[p.offset] = int256(p.initial);
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
        emit Model.SignalEvent(uint8(KonamiCodeModel.Actions.HALT), uint8(role), 1);
    }

    modifier inhibitAction(Model.Transition memory t) {
        require(!paused, "Game is paused.");
        require(!actionInhibited(t), "Action is inhibited.");
        _;
    }

    function transform(uint8 i, Model.Transition memory t, uint256 scalar) internal override {
        require(scalar == 1, 'invalid multiple');

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
            this.signal(uint8(KonamiCodeModel.Actions.Up), 1); // re-apply first move after reset
        }
    }

    function _transform(Model.Transition memory t) internal {
        Model.Place[] memory p = places;
        if (actionInhibited(t)) {
            resetLast(t.offset);
        }
        for (uint8 i = 0; i < uint8(KonamiCodeModel.Properties.SIZE); i++) {
            transform(i, t, 1);
            if (p[i].capacity != 0) {
                if (state[i] > int256(p[i].capacity)) {
                    resetLast(t.offset);
                }
            }
        }
    }

    function actionInhibited(Model.Transition memory t) internal view returns (bool) {
        for (uint8 i = 0; i < uint8(KonamiCodeModel.Properties.SIZE); i++) {
            if (_isInhibited(i, t)) {
                return true;
            }
        }
        return false;
    }

    function _isInhibited(uint8 i, Model.Transition memory t) internal view returns (bool) {
        if (t.guard[i] != 0) {
            return state[i] + t.guard[i] >= 0;
        }
        return false;
    }

    modifier autoIncrement() {
        _;
        sequence++;
    }

    function Up() external {
        this.signal(uint8(KonamiCodeModel.Actions.Up), 1);
    }

    function Down() external {
        this.signal(uint8(KonamiCodeModel.Actions.Down), 1);
    }

    function Left() external {
        this.signal(uint8(KonamiCodeModel.Actions.Left), 1);
    }

    function Right() external {
        this.signal(uint8(KonamiCodeModel.Actions.Right), 1);
    }

    function Select() external {
        this.signal(uint8(KonamiCodeModel.Actions.Select), 1);
    }

    function Start() external {
        this.signal(uint8(KonamiCodeModel.Actions.Start), 1);
    }

    function B() external {
        this.signal(uint8(KonamiCodeModel.Actions.B), 1);
    }

    function A() external {
        this.signal(uint8(KonamiCodeModel.Actions.A), 1);
    }

}
