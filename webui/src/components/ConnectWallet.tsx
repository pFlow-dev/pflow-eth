import React, {useState} from 'react';
import {ethers} from 'ethers';

declare global {
    interface Window {
        ethereum: any;
    }
}
type ConnectWalletProps = {
    metaModel: any;
}

const ConnectWallet = ({metaModel}: ConnectWalletProps) => {
    const [connected, setConnected] = useState(false);
    const [walletAddress, setWalletAddress] = useState("");

    async function connectWallet() {
        const provider = new ethers.BrowserProvider(window.ethereum);
        if (!connected) {
            const signer = await provider.getSigner();
            const _walletAddress = await signer.getAddress();
            setConnected(true);
            setWalletAddress(_walletAddress);
        } else {
            provider.destroy();
            setConnected(false);
            setWalletAddress("");
        }
    }

    async function callFaucet() {
        if (walletAddress) {
            const response = await fetch('/v0/faucet?addr=' + walletAddress, {
                method: 'GET',
            });
            const data = await response.json();
            if (!response.ok) {
                alert(JSON.stringify({...data}))
            } else {
                alert("+1 ETH tx:" + data.txHash)
            }
        }
    }

    return (
        <>
            <div className='app'>
                <div className='main'>
                    <button className='btn' onClick={connectWallet}>
                        {connected ? "Disconnect Wallet" : "Connect Wallet"}
                    </button>
                    {walletAddress}
                    {walletAddress !== "" && <button onClick={callFaucet}>Faucet</button>}
                </div>
            </div>
        </>
    );
}

export default ConnectWallet;