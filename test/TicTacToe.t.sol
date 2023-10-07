// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Test, console2} from "forge-std/Test.sol";
import {TicTacToe} from "../src/TicTacToe.sol";

contract TicTacToeTest is Test {
    TicTacToe public ticTacToe;
    address owner;
    address addrX;
    address addrO;

    function setUp() public {
        owner = address(this);
        addrX = address(1);
        addrO = address(2);
        ticTacToe = new TicTacToe(addrX, addrO);
    }

    function test_Increment() public {
        vm.prank(addrX);
        ticTacToe.X11();

        vm.prank(addrO);
        ticTacToe.O00();
    }
}
