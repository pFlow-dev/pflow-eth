// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;
import "@openzeppelin/contracts/access/AccessControl.sol";
import "./Metamodel.sol";

library TicTacToeModel {

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
}

abstract contract TicTacToeMetaModel is MetamodelUint8 {

    string public constant metamodelUri = "ipns://pflow.eth?contract=TicTacToe";

    // add an action to the model
    function _action(TicTacToeModel.Properties prop, TicTacToeModel.Actions action, TicTacToeModel.Roles role) internal {
        txn(1, places[uint8(prop)], fn(uint8(TicTacToeModel.Properties.SIZE), uint8(action), uint8(role)));
    }

    // declare model properties
    function _props() internal {
        cell(1, 1); // _00
        cell(1, 1); // _01
        cell(1, 1); // _02

        cell(1, 1); // _10
        cell(1, 1); // _11
        cell(1, 1); // _12

        cell(1, 1); // _20
        cell(1, 1); // _21
        cell(1, 1); // _22
    }

    // declare model actions
    function _actions() internal {
        _action(TicTacToeModel.Properties._00, TicTacToeModel.Actions.X00, TicTacToeModel.Roles.X);
        _action(TicTacToeModel.Properties._01, TicTacToeModel.Actions.X01, TicTacToeModel.Roles.X);
        _action(TicTacToeModel.Properties._02, TicTacToeModel.Actions.X02, TicTacToeModel.Roles.X);

        _action(TicTacToeModel.Properties._10, TicTacToeModel.Actions.X10, TicTacToeModel.Roles.X);
        _action(TicTacToeModel.Properties._11, TicTacToeModel.Actions.X11, TicTacToeModel.Roles.X);
        _action(TicTacToeModel.Properties._12, TicTacToeModel.Actions.X12, TicTacToeModel.Roles.X);

        _action(TicTacToeModel.Properties._20, TicTacToeModel.Actions.X20, TicTacToeModel.Roles.X);
        _action(TicTacToeModel.Properties._21, TicTacToeModel.Actions.X21, TicTacToeModel.Roles.X);
        _action(TicTacToeModel.Properties._22, TicTacToeModel.Actions.X22, TicTacToeModel.Roles.X);

        _action(TicTacToeModel.Properties._00, TicTacToeModel.Actions.O00, TicTacToeModel.Roles.O);
        _action(TicTacToeModel.Properties._01, TicTacToeModel.Actions.O01, TicTacToeModel.Roles.O);
        _action(TicTacToeModel.Properties._02, TicTacToeModel.Actions.O02, TicTacToeModel.Roles.O);

        _action(TicTacToeModel.Properties._10, TicTacToeModel.Actions.O10, TicTacToeModel.Roles.O);
        _action(TicTacToeModel.Properties._11, TicTacToeModel.Actions.O11, TicTacToeModel.Roles.O);
        _action(TicTacToeModel.Properties._12, TicTacToeModel.Actions.O12, TicTacToeModel.Roles.O);

        _action(TicTacToeModel.Properties._20, TicTacToeModel.Actions.O20, TicTacToeModel.Roles.O);
        _action(TicTacToeModel.Properties._21, TicTacToeModel.Actions.O21, TicTacToeModel.Roles.O);
        _action(TicTacToeModel.Properties._22, TicTacToeModel.Actions.O22, TicTacToeModel.Roles.O);
    }

    constructor() {
        _props();
        _actions();
    }

}

/// @custom:security-contact security@stackdump.com
contract TicTacToe is TicTacToeMetaModel, AccessControl {

    address internal owner;

    bool internal paused = false;
    uint8 internal sequence = 0;
    uint256 internal session = 0;

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
        session++;
        _;
        Uint8Model.Place[] memory p = places;
        _init(p[0]); // _00
        _init(p[1]); // _01
        _init(p[2]); // _02
        _init(p[3]); // _10
        _init(p[4]); // _11
        _init(p[5]); // _12
        _init(p[6]); // _20
        _init(p[7]); // _21
        _init(p[8]); // _22
    }

    function resetGame(TicTacToeModel.Roles role) internal startGame {
        emit Uint8Model.SignalEvent(session, sequence, uint8(TicTacToeModel.Actions.HALT), uint8(role), block.timestamp);
    }

    function transform(uint8 i, Uint8Model.Transition memory t)  internal override {
        if (t.delta[i] != 0) {
            state[i] = state[i] + t.delta[i];
            require(state[i] >= 0, "Invalid state");
        }
    }

    modifier takeTurns() {
        testIsMyTurn();
        _;
        sequence++;
    }

    function signal(TicTacToeModel.Actions action) internal takeTurns {
        uint8 txnId = uint8(action);
        Uint8Model.Transition memory t = transitions[txnId];
        assert(txnId == t.offset);
        transform(0, t); // _00
        transform(1, t); // _01
        transform(2, t); // _02
        transform(3, t); // _10
        transform(4, t); // _11
        transform(5, t); // _12
        transform(6, t); // _20
        transform(7, t); // _21
        transform(8, t); // _22
        emit Uint8Model.SignalEvent(session, sequence, txnId, t.role, block.timestamp);
    }

    function reset() public {
        resetGame(getRole());
    }

    function testIsMyTurn() public view  {
        testIsGameOpen();
            require(!hasRole(DEFAULT_ADMIN_ROLE, msg.sender), "Gameplay from contract admin is forbidden");
        if (sequence % 2 == 0) {
            require(getRole() == TicTacToeModel.Roles.X, "X's turn");
        } else {
            require(getRole() == TicTacToeModel.Roles.O, "O's turn");
        }
    }

    function getRole() public view returns (TicTacToeModel.Roles) {
        if (hasRole(PLAYER_X, msg.sender)) {
            return TicTacToeModel.Roles.X;
        } else if (hasRole(PLAYER_O, msg.sender)) {
            return TicTacToeModel.Roles.O;
        } else {
            revert("Unexpected caller");
        }
    }

    function X00() external {
        signal(TicTacToeModel.Actions.X00);
    }

    function X01() external {
        signal(TicTacToeModel.Actions.X01);
    }

    function X02() external {
        signal(TicTacToeModel.Actions.X02);
    }

    function X10() external {
        signal(TicTacToeModel.Actions.X10);
    }

    function X11() external {
        signal(TicTacToeModel.Actions.X11);
    }

    function X12() external {
        signal(TicTacToeModel.Actions.X12);
    }

    function X20() external {
        signal(TicTacToeModel.Actions.X20);
    }

    function X21() external {
        signal(TicTacToeModel.Actions.X21);
    }

    function X22() external {
        signal(TicTacToeModel.Actions.X22);
    }

    function O00() external {
        signal(TicTacToeModel.Actions.O00);
    }

    function O01() external {
        signal(TicTacToeModel.Actions.O01);
    }

    function O02() external {
        signal(TicTacToeModel.Actions.O02);
    }

    function O10() external {
        signal(TicTacToeModel.Actions.O10);
    }

    function O11() external {
        signal(TicTacToeModel.Actions.O11);
    }

    function O12() external {
        signal(TicTacToeModel.Actions.O12);
    }

    function O20() external {
        signal(TicTacToeModel.Actions.O20);
    }

    function O21() external {
        signal(TicTacToeModel.Actions.O21);
    }

    function O22() external {
        signal(TicTacToeModel.Actions.O22);
    }
}
