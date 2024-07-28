import '@typechain/hardhat'
import {HardhatUserConfig} from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";
import '@nomicfoundation/hardhat-ethers'
import '@nomicfoundation/hardhat-chai-matchers'
import "hardhat-contract-sizer"


const config: HardhatUserConfig = {
    solidity: {
        version: "0.8.24",
        settings: {
            optimizer: {
                enabled: true,
                runs: 200,
            },
        },
    },
    networks: {
        localhost: {
            url: "http://localhost:8545",
        },
        hardhat: {
            allowUnlimitedContractSize: false,
            throwOnTransactionFailures: false,
            throwOnCallFailures: false,
        },
    },
};

export default config;
