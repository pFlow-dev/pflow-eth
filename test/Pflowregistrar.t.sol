// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Test} from "forge-std/Test.sol";
import "forge-std/console.sol";

import {Declaration, Model} from "../contracts/Metamodel.sol";
import {TicTacToe} from "./examples/TicTacToe.sol";
import {PflowRegistrar, PflowEth} from "../contracts/PflowRegistrar.sol";
import "../contracts/Metamodel.sol";

contract ModelRegistryTest is Test {
    TicTacToe public ticTacToe;
    PflowRegistrar public api;
    address owner;
    address addrX;
    address addrO;


    function setUp() public {
        owner = address(this);
        addrX = address(1);
        addrO = address(2);

        ticTacToe = new TicTacToe(addrX, addrO);
        api = new PflowRegistrar();
    }

    function test_importModel() public {
        vm.prank(owner);
        api.register(address(ticTacToe));
        PflowEth.FlowInfo[] memory flows = api.getFlows();
        assertEq(flows.length, 1);
        assertNotEq(flows[0].flowId, bytes32(0));
    }
}
