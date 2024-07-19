import React, {useEffect, useState} from 'react';
import styles from './Signal.module.css';

interface SignalData {
    nonce: number;
    response: {
        contract: string;
        event_log: Array<{
            data: string;
            event: string[];
        }>;
        sender: string;
        transaction_hash: string;
    };
}

type SignalProps = {
    metaModel: any;
}

function Signal({metaModel}: SignalProps) {
    const [signalData, setSignalData] = useState<SignalData | null>(null);
    const [isLoading, setIsLoading] = useState<boolean>(true);
    const [action, setAction] = useState<string>('');
    const [scalar, setScalar] = useState<string>('');

    useEffect(() => {
        setSignalData({
            nonce: 0,
            response: {
                contract: '',
                event_log: [],
                sender: '',
                transaction_hash: ''
            }
        });
        setIsLoading(false);
    }, []);

    const handleActionChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setAction(e.target.value);
    };

    const handleScalarChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setScalar(e.target.value);
    };

    const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        metaModel.loadFromContract();
        e.preventDefault();
    };

    if (isLoading) {
        return <div>Loading...</div>;
    }

    if (!signalData) {
        return <div>No data available</div>;
    }

    return (
        <div className={styles.signalContainer}>
            <div className={styles.signalForm}>
                <p>Contract Call - Metamodel::Signal()</p>
                <form onSubmit={handleSubmit}>
                    <input type="text" value={action} onChange={handleActionChange} placeholder="Action"/>
                    <input type="text" value={scalar} onChange={handleScalarChange} placeholder="Scalar"/>
                    <button type="submit">Submit</button>
                </form>
                <pre>{JSON.stringify(signalData, null, 2)}</pre>
            </div>
        </div>
    );
}

export default Signal;