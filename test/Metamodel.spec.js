const {
  time,
  loadFixture,
} = require("@nomicfoundation/hardhat-network-helpers");
require("@nomiclabs/hardhat-web3");
const { ethers } = require("hardhat");

describe("Metamodel", function () {
  async function deployTestProxy() {
    // Contracts are deployed using the first signer/account by default
    const [owner, p0, p1] = await ethers.getSigners();

    // const Lib = await ethers.getContractFactory("contracts/Metamodel.sol:TicTacToeModel");
    // const lib = await Lib.deploy();
    // await lib.deployed();
    const contract = await ethers.getContractFactory("contracts/TicTacToe.sol:TicTacToe", {
        libraries: {
            // TicTacToeModel: lib.address,
        },
    });
    const api = await contract.deploy(p0.address, p1.address); // REVIEW: swap to test access ctl
    // console.log(api.from);

    return { contract, api, p0, p1 };
  }

  describe("Deployment", function () {

    it("should allow gameplay", async function () {
      const { api , p0, p1} = await loadFixture(deployTestProxy);
      // console.log({api})
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
});
