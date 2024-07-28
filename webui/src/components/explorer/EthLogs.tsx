import React, {useEffect, useState} from "react";
import {ethers} from "ethers";
import styles from "./EthLogs.module.css";
import {MetaModel} from "../../lib/pflow";


function uint256HexToNumeric(val: string): bigint {
    // Remove leading '0x' if present
    if (val.startsWith('0x')) {
        val = val.slice(2);
    }

    // Convert hex substrings to numeric values and sum them up
    const valA = BigInt('0x' + val.slice(0, 8)) * (2n ** 224n);
    const valB = BigInt('0x' + val.slice(8, 22)) * (2n ** 168n);
    const valC = BigInt('0x' + val.slice(22, 36)) * (2n ** 112n);
    const valD = BigInt('0x' + val.slice(36, 50)) * (2n ** 56n);
    const valE = BigInt('0x' + val.slice(50, 64));

    return valA + valB + valC + valD + valE;
}

interface EthLog {
    sequence: number;
    block_number: number;
    role: number;
    action: number;
    scalar: number;
    from_address: string;
    transaction_hash: string;
    topic_hash: string;
    removed: boolean;
}

interface ethLogOpts {
    metaModel: MetaModel;
}

export default function EthLogs({metaModel}: ethLogOpts) {
    const [logs, setLogs] = useState<EthLog[]>([]);

    const provider = new ethers.BrowserProvider(window.ethereum);

    const fetchEthLogs = async () => {
        try {
            const filter = {
                fromBlock: "earliest",
                toBlock: "latest",
                address: metaModel.address,
                topics: [] // Add any topics if needed
            };

            const rawLogs = await provider.send("eth_getLogs", [filter]);

            const formattedLogs = rawLogs.map((log: any, index: number) => ({
                sequence: uint256HexToNumeric(log.data).toString(10),
                block_number: parseInt(log.blockNumber, 16),
                role: uint256HexToNumeric(log.topics[1]).toString(10),
                action: uint256HexToNumeric(log.topics[2]).toString(10),
                scalar: uint256HexToNumeric(log.topics[3]).toString(10),
                transaction_hash: log.transactionHash,
                topic_hash: log.topics[0],
                removed: log.removed
            }));

            setLogs(formattedLogs.reverse());
        } catch (error) {
            console.error("Error fetching logs:", error);
        }
    };

    useEffect(() => {
        fetchEthLogs(); // Initial fetch
        const interval = setInterval(fetchEthLogs, 30000); // Refresh every 30 seconds
        return () => clearInterval(interval); // Cleanup on unmount
    }, []);

    const copyToClipboard = (text: string) => {
        navigator.clipboard.writeText(text).then(() => {
            console.log("Text copied to clipboard");
        }).catch(err => {
            console.error("Failed to copy text: ", err);
        });
    };
    const actionLabels: { [key: number]: string } = {};

    metaModel.petriNet.def.transitions.forEach((transition) => {
        actionLabels[transition.offset] = transition.label;
    })

    return (
        <div>
            <div className={styles.eventLogTopic}>
                <button onClick={fetchEthLogs}>Refresh Logs</button>
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
                                <td>{actionLabels[log.action] || log.action}</td>
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
    );
}