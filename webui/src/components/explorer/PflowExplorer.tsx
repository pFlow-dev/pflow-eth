import React, {Fragment, useState} from "react";
import {ethers} from "ethers";
import PetriNet from "../model/PetriNet";
import {MetaModel, newModel} from "../../lib/pflow";
import styles from "./PflowExplorer.module.css";
import EthLogs from "./EthLogs";

const sepoliaTestModel = "0x9265fd7b41b3f96c3123319b713a5c5a761981f1";
const defaultAddress = sepoliaTestModel;

declare global {
    interface Window {
        ethereum: any;
    }
}

export default function PflowExplorer() {
    const emptyModel: MetaModel = new MetaModel({
        updateHook: (m) => {
            setMetaModel(m);
        },
        address: defaultAddress,
        initialModel: newModel({
            declaration: () => {
            }
        }),
    });
    const [metaModel, setMetaModel] = useState<MetaModel>(emptyModel);
    const [connected, setConnected] = useState(false);
    const [walletAddress, setWalletAddress] = useState("");
    const [userSignedIn, setUserSignedIn] = useState(false);
    const [signingMessage, setSigningMessage] = useState(false);

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
            setUserSignedIn(false);
            setWalletAddress("");
        }
    }

    async function signInWithEthereum() {
        if (window.ethereum) {
            try {
                const provider = new ethers.BrowserProvider(window.ethereum);
                const signer = await provider.getSigner();
                setSigningMessage(true);
                const address = await signer.getAddress();
                const ts = (new Date()).toISOString();
                const msg = "authenticate " + ts;
                const signature = await signer.signMessage(msg);
                const session = metaModel.getSession();
                return fetch("/v0/authenticate?address=" + address + "&signature=" + signature + "&ts=" + ts + "&session=" + session, {
                    method: "GET",
                    headers: {
                        "Content-Type": "application/json",
                    },
                }).then((response) => {
                    return response.json().then((response) => {
                        if (response.success) {
                            console.log({
                                address,
                                signature,
                                session: response.data.session,
                            }, 'login_success');
                            metaModel.setSession(response.data.session);
                        }
                    }).then(() => {
                        return metaModel.pollServer().then(() => {
                            setUserSignedIn(true);
                            setSigningMessage(false);
                            return metaModel.loadFromContract().then(() => {
                                metaModel.update();
                            })
                        });
                    })
                });
            } catch (error) {
                console.error("Error signing message:", error);
                setSigningMessage(false);
            }
        } else {
            alert("Ethereum wallet is not connected");
        }
    }

    const shortAddress = walletAddress ? walletAddress.slice(0, 6) + "..." + walletAddress.slice(-4) : "";

    const connectedButNotSignedIn = connected && !userSignedIn;

    return (<div className={styles.explorerContainer}>
        <div>
            <button className="btn" onClick={connectWallet}>
                {connected ? shortAddress : "Connect Wallet"}
            </button>
            {connectedButNotSignedIn && (
                <button className="btn" onClick={signInWithEthereum} disabled={signingMessage}>
                    {"Sign In"}
                </button>
            )}

            {signingMessage && <p>Please check your wallet to sign the message...</p>}

            &nbsp;<input type="text" size={42} className={styles.contractLink} value={metaModel.address} readOnly/>
        </div>

        {!userSignedIn && <div className={styles.licenseText}>
            Connect Wallet and Sign In to interact with contract.<br/>
        </div>}
        {userSignedIn && <Fragment>
            <svg className={styles.stateMachineContainer} id="pflow-svg-outer">
                <PetriNet metaModel={metaModel}/>
            </svg>
            <EthLogs metaModel={metaModel}/>
        </Fragment>}
    </div>);
}