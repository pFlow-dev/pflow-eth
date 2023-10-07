// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Test, console2} from "forge-std/Test.sol";
import {KonamiCode} from "../src/KonamiCode.sol";
import "forge-std/console.sol";

contract KonamiCodeTest is Test {
    KonamiCode public konamiCode;
    address public owner;
    address public unlocker;

    function setUp() public {
        owner = address(this);
        konamiCode = new KonamiCode();
    }

    function test_Increment() public {
        vm.startPrank(unlocker);
        // assertEq(konamiCode.session(), 1);
        // konamiCode.Up();
        assertEq(konamiCode.session(), 1);
        konamiCode.Up();
        konamiCode.Up();
        assertEq(konamiCode.session(), 1);
        konamiCode.Down();
        konamiCode.Down();
        konamiCode.Left();
        konamiCode.Right();
        konamiCode.Left();
        konamiCode.Right();
        konamiCode.B();
        konamiCode.A();
        konamiCode.Select();
        konamiCode.Start();
        assertEq(konamiCode.session(), 1);
    }
}
