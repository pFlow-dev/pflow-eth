import React from "react";
import ReactDOM from "react-dom/client";
import PflowExplorer from "./components/explorer/PflowExplorer";
import {createWeb3Modal, defaultConfig} from '@web3modal/ethers/react'
import Header from "./components/Header";
import Footer from "./components/Footer";

export const ethereum = {
    chainId: 1,
    name: 'Ethereum',
    currency: 'ETH',
    explorerUrl: 'https://etherscan.io',
    rpcUrl: 'https://eth-mainnet.blastapi.io/4fd309c5-6555-4052-99d5-3a6b646f14b4',
}

export const optimism = {
    chainId: 10,
    name: 'Optimism Mainnet',
    currency: 'ETH',
    explorerUrl: 'https://optimistic.etherscan.io',
    rpcUrl: 'https://optimism-mainnet.blastapi.io/4fd309c5-6555-4052-99d5-3a6b646f14b4',
}

export const hardhat = {
    chainId: 31337,
    name: 'Hardhat Testnet',
    currency: 'ETH',
    explorerUrl: '',
    rpcUrl: 'http://localhost:8545',
    defaultAddress: '0x5FbDB2315678afecb367f032d93F642f64180aa3',
}

export const sepoliaOptimism = {
    chainId: 11155420,
    name: 'Sepolia Optimism Testnet',
    currency: 'ETH',
    explorerUrl: 'https://sepolia-optimism.etherscan.io',
    rpcUrl: "https://optimism-sepolia.blastapi.io/4fd309c5-6555-4052-99d5-3a6b646f14b4"
}

export const sepolia = {
    chainId: 11155111,
    name: 'Sepolia Testnet',
    currency: 'ETH',
    explorerUrl: 'https://sepolia.etherscan.io',
    rpcUrl: 'https://eth-sepolia.blastapi.io/4fd309c5-6555-4052-99d5-3a6b646f14b4',
}

export const projectId = 'd76cef25c0687c2391805b1d2864ba8f';

// 3. Create a metadata object
const metadata = {
    name: 'pflow.xyz',
    description: 'Metamodel Explorer',
    url: 'https://app.pflow.xyz', // origin must match your domain & subdomain
    icons: [ 'https://avatars.githubusercontent.com/u/86532620' ]
}

// 4. Create Ethers config
const ethersConfig = defaultConfig({
    /*Required*/
    metadata,

    auth: {
        email: true,
        socials: ['github'],
        showWallets: true,
        walletFeatures: true
    },

    /*Optional*/
    enableEIP6963: true, // true by default
    enableInjected: true, // true by default
    enableCoinbase: true, // true by default
    defaultChainId: 1, // used for the Coinbase SDK
})

// 5. Create a Web3Modal instance
createWeb3Modal({
    ethersConfig,
    chains: [ethereum, sepolia, optimism, sepoliaOptimism, hardhat],
    projectId,
    enableAnalytics: true // Optional - defaults to your Cloud configuration
})

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
    <React.StrictMode>
        <Header/>
        <PflowExplorer/>
        <Footer/>
    </React.StrictMode>
);
