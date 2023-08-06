const {
    time,
    loadFixture,
} = require("@nomicfoundation/hardhat-network-helpers");
require("@nomiclabs/hardhat-web3");
const {ethers} = require("hardhat");

describe("TicTacToe", function () {
    async function deployTestProxy() {
        // Contracts are deployed using the first signer/account by default
        const [owner, p0, p1] = await ethers.getSigners();

        // const Lib = await ethers.getContractFactory("contracts/Metamodel.sol:TicTacToeModel");
        // const lib = await Lib.deploy();
        // await lib.deployed();
        const contractFactory = await ethers.getContractFactory("contracts/TicTacToe.sol:TicTacToe", {
            libraries: {
                // TicTacToeModel: lib.address,
            },
        });
        const k = await contractFactory.deploy(p0.address, p1.address); // REVIEW: swap to test access ctl
        // console.log(api.from);

        return {contract: contractFactory, api: k, p0, p1};
    }

    it("should allow gameplay", async function () {
        const {api, p0, p1} = await loadFixture(deployTestProxy);

        const x = await api.connect(p0);
        const o = await api.connect(p1);

        await x.X11();
        await o.O01();
        await x.X00();
        await o.O02();
        await x.X22(); // X wins

        await o.testIsMyTurn();
        await x.reset();
        await x.testIsMyTurn();

        await x.X11();
        await o.O01();
        await x.X00();
        await o.O02();
        await x.X22(); // X wins
    });

});
