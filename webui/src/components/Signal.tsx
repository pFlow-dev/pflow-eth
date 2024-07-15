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

function Signal() {
    const [signalData, setSignalData] = useState<SignalData | null>(null);
    const [isLoading, setIsLoading] = useState<boolean>(true);
    const [action, setAction] = useState<string>('');
    const [scalar, setScalar] = useState<string>('');

    useEffect(() => {
        fetchData('/v0/signal?action=2,2&scalar=1,3');
    }, []);

    const handleActionChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setAction(e.target.value);
    };

    const handleScalarChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setScalar(e.target.value);
    };

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        const url = `/v0/signal?action=${action}&scalar=${scalar}`;
        fetchData(url);
    };

    const fetchData = async (url: string) => {
        setIsLoading(true);
        try {
            const response = await fetch(url);
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            const data: SignalData = await response.json();
            setSignalData(data);
        } catch (error) {
            console.error('Error fetching signal data:', error);
        } finally {
            setIsLoading(false);
        }
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