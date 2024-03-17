// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;

import "@openzeppelin/contracts/access/AccessControl.sol";
import "../../contracts/Metamodel.sol";
library TicTacToeModel {

    enum Roles{X, O, HALT}

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
        // o moves
        O00, O01, O02,
        O10, O11, O12,
        O20, O21, O22,
        HALT
    }
}

abstract contract TicTacToeMetaModel is Metamodel {

    // add an action to the model
    function _action(string memory label, TicTacToeModel.Properties prop, TicTacToeModel.Actions action, TicTacToeModel.Roles role, uint8 x, uint8 y) internal {
        arrow(1, places[uint8(prop)], func(
            label,
            uint8(TicTacToeModel.Properties.SIZE),
            uint8(action),
            uint8(role),
            Model.Position(x, y)
        ));
    }

    function _place(string memory label, uint8 x, uint8 y) internal {
        cell(label, 1, 1, Model.Position(x, y));
    }

    // declare model properties
    function _props() internal {
        _place("00", 1, 1); // _00
        _place("01", 2, 1); // _01
        _place("02", 3, 1); // _02

        _place("10", 1, 2); // _10
        _place("11", 2, 2); // _11
        _place("12", 3, 2); // _12

        _place("20", 1, 3); // _20
        _place("21", 2, 3); // _21
        _place("22", 3, 3); // _22
    }

    // declare model actions
    function _actions() internal {
        _action("X00", TicTacToeModel.Properties._00, TicTacToeModel.Actions.X00, TicTacToeModel.Roles.X, 5, 1);
        _action("X01", TicTacToeModel.Properties._01, TicTacToeModel.Actions.X01, TicTacToeModel.Roles.X, 6, 1);
        _action("X02", TicTacToeModel.Properties._02, TicTacToeModel.Actions.X02, TicTacToeModel.Roles.X, 7, 1);

        _action("X10", TicTacToeModel.Properties._10, TicTacToeModel.Actions.X10, TicTacToeModel.Roles.X, 5, 2);
        _action("X11", TicTacToeModel.Properties._11, TicTacToeModel.Actions.X11, TicTacToeModel.Roles.X, 6, 2);
        _action("X12", TicTacToeModel.Properties._12, TicTacToeModel.Actions.X12, TicTacToeModel.Roles.X, 7, 2);

        _action("X20", TicTacToeModel.Properties._20, TicTacToeModel.Actions.X20, TicTacToeModel.Roles.X, 5, 3);
        _action("X21", TicTacToeModel.Properties._21, TicTacToeModel.Actions.X21, TicTacToeModel.Roles.X, 6, 3);
        _action("X22", TicTacToeModel.Properties._22, TicTacToeModel.Actions.X22, TicTacToeModel.Roles.X, 7, 3);

        _action("000", TicTacToeModel.Properties._00, TicTacToeModel.Actions.O00, TicTacToeModel.Roles.O, 1, 5);
        _action("O01", TicTacToeModel.Properties._01, TicTacToeModel.Actions.O01, TicTacToeModel.Roles.O, 1, 6);
        _action("O02", TicTacToeModel.Properties._02, TicTacToeModel.Actions.O02, TicTacToeModel.Roles.O, 1, 7);

        _action("O10", TicTacToeModel.Properties._10, TicTacToeModel.Actions.O10, TicTacToeModel.Roles.O, 2, 5);
        _action("O11", TicTacToeModel.Properties._11, TicTacToeModel.Actions.O11, TicTacToeModel.Roles.O, 2, 6);
        _action("O12", TicTacToeModel.Properties._12, TicTacToeModel.Actions.O12, TicTacToeModel.Roles.O, 2, 7);

        _action("O20", TicTacToeModel.Properties._20, TicTacToeModel.Actions.O20, TicTacToeModel.Roles.O, 3, 5);
        _action("O21", TicTacToeModel.Properties._21, TicTacToeModel.Actions.O21, TicTacToeModel.Roles.O, 3, 6);
        _action("O22", TicTacToeModel.Properties._22, TicTacToeModel.Actions.O22, TicTacToeModel.Roles.O, 3, 7);
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

    int256[] public state = new int256[](uint256(TicTacToeModel.Properties.SIZE));

    bytes32 public constant PLAYER_X = keccak256("PLAYER_X");
    bytes32 public constant PLAYER_O = keccak256("PLAYER_O");

    address PlayerX; // REVIEW: should we replace role usage with address? or is there a way to lookup from Access control
    address PlayerO;

    constructor(address p0, address p1) {
        owner = tx.origin;
        require(p0 != p1, "Players must not have the same address.");
        PlayerX = p0;
        PlayerO = p1;
        _grantRole(DEFAULT_ADMIN_ROLE, owner);
        _grantRole(PLAYER_X, p0);
        _grantRole(PLAYER_O, p1);
        resetGame(TicTacToeModel.Roles.HALT);
    }

    // REVIEW: should have additional checks for gateway to limit access from outside the router?

    // function openGateway(address gateway) external override {
    //     require(msg.sender == owner, "Only owner can open gateway");
    //     gateways[gateway] = true;
    // }

    // function closeGateway(address gateway) external override {
    //     require(msg.sender == owner, "Only owner can close gateway");
    //     gateways[gateway] = false;
    // }

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
        _;
        for (uint8 i = 0; i < uint8(places.length); i++) {
            if (state[places[i].offset] != int256(places[i].initial)) {
                state[places[i].offset] = int256(places[i].initial);
            }
        }
    }

    function resetGame(TicTacToeModel.Roles role) internal startGame {
        emit Model.SignalEvent(uint8(role), uint8(TicTacToeModel.Actions.HALT), 1);
    }

    function transform(uint8 i, Model.Transition memory t, uint256 scalar) internal override {
        testIsMyTurn();
        require(scalar == 1, "Invalid multiple");
        if (t.delta[i] != 0) {
            state[i] = state[i] + t.delta[i];
            require(state[i] >= 0, "Invalid state");
        }
    }

    function reset() public {
        resetGame(getRole());
    }

    function testIsMyTurn() public view {
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
            revert("msg.sender is not a player");
        }
    }

}
