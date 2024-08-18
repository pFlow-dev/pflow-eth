import React, { useState, useEffect } from "react";
import PetriNet from "../model/PetriNet";
import {MetaModel, newModel, EthLog} from "../../lib/pflow";
import styles from "./PflowExplorer.module.css";
import { useWeb3ModalProvider } from '@web3modal/ethers/react'
import { BrowserProvider } from 'ethers'

export default function PflowExplorer() {
    const [logs, setLogs] = useState<EthLog[]>([]);
    function copyToClipboard(text: string) {
        navigator.clipboard.writeText(text);
    }
    const { walletProvider } = useWeb3ModalProvider()

    const [provider, setProvider] = useState<BrowserProvider | null>(null);

    const [actionLabels, setActionLabels] = useState<Map<string, string>>(new Map<string, string>());

    const initializeMetaModel = (): MetaModel => {
       return new MetaModel({
           updateHook: (model) => setMetaModel(model),
           initialModel: newModel({ declaration: () => {} }),
        });
    }
    const [metaModel, setMetaModel] = useState<MetaModel>(initializeMetaModel());


    useEffect(() => {
        if (walletProvider) {
            setProvider(new BrowserProvider(walletProvider));
        }
    }, [walletProvider]);

    useEffect(() => {

        async function loadChainData(provider: BrowserProvider): Promise<void> {
            return metaModel.loadFromContract(provider).then(async () => {
                metaModel.petriNet.def.transitions.forEach((transition) => {
                    actionLabels.set(`${transition.offset}`, transition.label);
                });
                setActionLabels(actionLabels);
                const logs = await metaModel.fetchEthLogs(provider);
                logs.reverse();
                setLogs(logs);
            }).finally(() => {
                setMetaModel(metaModel);
            })
        }

        if (provider) {
            loadChainData(provider).finally(() => {});
        }
    }, [provider, actionLabels, metaModel]);

    return (
        <div className={styles.explorerContainer}>
            <svg className={styles.stateMachineContainer} id="pflow-svg-outer">
                <PetriNet metaModel={metaModel} />
            </svg>
            <div>
                <div className={styles.eventLogTopic}>
                    &nbsp;Topic: {logs.length > 0 ? logs[0].topic_hash : ""}
                </div>
                <div className={styles.logsContainer}>
                    {logs.length > 0 ? (
                        <table className={styles.logsTable}>
                            <thead>
                            <tr>
                                <th>Seq</th>
                                <th>Block</th>
                                <th>Role</th>
                                <th>Action</th>
                                <th>Scalar</th>
                                <th>Transaction Hash</th>
                                <th>Removed</th>
                            </tr>
                            </thead>
                            <tbody>
                            {logs.map((log, index) => (
                                <tr key={index} className={log.removed ? styles.removedLog : ""}>
                                    <td>{log.sequence}</td>
                                    <td>{log.block_number}</td>
                                    <td>{log.role}</td>
                                    <td>{actionLabels.get(log.action) || log.action}</td>
                                    <td>{log.scalar}</td>
                                    <td className={styles.truncated}
                                        onClick={() => copyToClipboard(log.transaction_hash)}>{log.transaction_hash}</td>
                                    <td>{log.removed ? "Yes" : "No"}</td>
                                </tr>
                            ))}
                            </tbody>
                        </table>
                    ) : (
                        <p>No logs found.</p>
                    )}
                </div>
            </div>
        </div>
    );
}