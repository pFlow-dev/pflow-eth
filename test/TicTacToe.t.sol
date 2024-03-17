// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Test} from "forge-std/Test.sol";
import {console} from "forge-std/console.sol";

import {Declaration} from "../contracts/Metamodel.sol";
import {TicTacToe, TicTacToeModel} from "./examples/TicTacToe.sol";

contract TicTacToeTest is Test {
    TicTacToe public ticTacToe;
    address owner;
    address addrX;
    address addrO;
    address[] players;

    function setUp() public {
        owner = address(this);
        addrX = address(1);
        addrO = address(2);
        ticTacToe = new TicTacToe(addrX, addrO);
        // ticTacToe.openGateway(addrX);
        // ticTacToe.openGateway(addrO);
    }

    function test_Increment() public {
        vm.prank(addrX);
        ticTacToe.signal(uint8(TicTacToeModel.Actions.X00), uint256(1));

        vm.prank(addrO);
        ticTacToe.signal(uint8(TicTacToeModel.Actions.O11), uint256(1));
    }

    function test_declaration() public {
        Declaration.PetriNet memory model = ticTacToe.declaration();
        assertEq(model.places.length, 9);
    }
}
