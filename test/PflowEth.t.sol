
// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Test, console2} from "forge-std/Test.sol";
import {PflowEth} from "../src/PflowEth.sol";

contract PflowEthTest is Test {
    PflowEth public erc;
    address owner;
    // address addrX;
    // address addrO;

    function setUp() public {
        owner = address(this);
        // addrX = address(1);
        // addrO = address(2);
        erc = new PflowEth();
    }

    function test_Increment() public {
        erc.pause();
    }
}
