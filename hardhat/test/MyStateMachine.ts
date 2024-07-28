import {expect} from "chai";
import hre from "hardhat";

describe("MyStateMachine", function () {
    async function deployFixture() {
        const [owner] = await hre.ethers.getSigners();
        const MyStateMachine = await hre.ethers.getContractFactory("MyStateMachine");
        const myStateMachine = await MyStateMachine.deploy();
        await myStateMachine.waitForDeployment();
        return {
            owner,
            myStateMachine,
        }

    }

    it("Should deploy the contract", async function () {
        const {myStateMachine} = await deployFixture();
        expect(await myStateMachine.getAddress()).to.be.properAddress;
        //const addr = await model.getAddress();
        //console.log({ address: addr});
        //expect(addr).to.equal('0x5FbDB2315678afecb367f032d93F642f64180aa3');
    });

});
