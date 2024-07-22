import React, {useEffect, useState} from "react";
import styles from "./Transactions.module.css";

interface Transaction {
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

export default function TransitionList() {
    const [transactions, setTransactions] = useState<Transaction[]>([]);

    const fetchTransactions = () => {
        fetch("/v0/transactions")
            .then(response => response.json())
            .then(data => setTransactions(data))
            .catch(error => console.error("Error fetching transactions:", error));
    };

    useEffect(() => {
        fetchTransactions(); // Initial fetch
        const interval = setInterval(fetchTransactions, 30000); // Refresh every 30 seconds
        return () => clearInterval(interval); // Cleanup on unmount
    }, []);

    const copyToClipboard = (text: string) => {
        navigator.clipboard.writeText(text).then(() => {
            console.log("Text copied to clipboard");
        }).catch(err => {
            console.error("Failed to copy text: ", err);
        });
    };

    return (
        <div>
            <button onClick={fetchTransactions}>Refresh Transactions</button>
            <div className={styles.transactionsContainer}>
                {transactions.length > 0 ? (
                    <table className={styles.transactionTable}>
                        <thead>
                        <tr>
                            <th>Seq</th>
                            <th>Block</th>
                            <th>Role</th>
                            <th>Action</th>
                            <th>Scalar</th>
                            <th>From Address</th>
                            <th>Transaction Hash</th>
                            <th>Removed</th>
                        </tr>
                        </thead>
                        <tbody>
                        {transactions.map((transaction, index) => (
                            <tr key={index} className={transaction.removed ? styles.removedTransaction : ""}>
                                <td>{transaction.sequence}</td>
                                <td>{transaction.block_number}</td>
                                <td>{transaction.role}</td>
                                <td>{transaction.action}</td>
                                <td>{transaction.scalar}</td>
                                <td className={styles.truncated}
                                    onClick={() => copyToClipboard(transaction.from_address)}>{transaction.from_address}</td>
                                <td className={styles.truncated}
                                    onClick={() => copyToClipboard(transaction.transaction_hash)}>{transaction.transaction_hash}</td>
                                <td>{transaction.removed ? "Yes" : "No"}</td>
                            </tr>
                        ))}
                        </tbody>
                    </table>
                ) : (
                    <p>No transactions found.</p>
                )}
            </div>
            &nbsp;EventLog Topic: {transactions.length > 0 ? transactions[0].topic_hash : ""}
        </div>
    );
}