// We require the Hardhat Runtime Environment explicitly here. This is optional
// but useful for running the script in a standalone fashion through `node <script>`.
//
// You can also run a script with `npx hardhat run <script>`. If you do that, Hardhat
// will compile your contracts, add the Hardhat Runtime Environment's members to the
// global scope, and execute the script.
const hre = require("hardhat");
/*

/*
hardhat run --network localhost  ./scripts/deploy.js
*/

async function main() {
  // const currentTimestampInSeconds = Math.round(Date.now() / 1000);
  // const unlockTime = currentTimestampInSeconds + 60;
  // const lockedAmount = hre.ethers.utils.parseEther("0.001");


  const ze = "0xCae1d2Aa66E01daCf90a655519099620cbf85B72"
  const z3 = "0xE3BD13fc489a19F59794Eed531f0aF6225D50623"
  // const me = "0xd20F93E2D8f7378946E8642F7579723B9A81544A"

  async function deploy(name) {
    const mm = await hre.ethers.getContractFactory(name);
    const model = await mm.deploy(ze, z3)
    await model.deployed();

    console.log( `${name} deployed to ${model.address}` );
  }

  // await deploy("MetamodelUint8")
  // await deploy("TicTacToeModel")
  const contractName = "contracts/Metamodel_flattened.sol:TicTacToe"
  await deploy(contractName); // KLUDGE note that we are using flattened version

}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
