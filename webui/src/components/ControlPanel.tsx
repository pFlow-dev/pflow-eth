import React, {useState} from 'react';

function ControlPanel() {
    const [response, setResponse] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState('');

    const sendControlCommand = async (action: string) => {
        setIsLoading(true);
        setError('');
        try {
            const response = await fetch(`/v0/control?cmd=${action}`, {
                method: 'GET',
            });
            const data = await response.json();
            if (response.ok) {
                setResponse(data.result);
            } else {
                setError(data.result || 'An error occurred');
            }
        } catch (err) {
            setError('Network error');
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div>
            <h2>Control Panel</h2>
            <p>This control panel allows you to send commands to the server for various administrative tasks.</p>
            <ul>
                <li><strong>Reset Database:</strong> Clears all data from the database and resets it to its initial
                    state.
                </li>
                <li><strong>Initialize Block Numbers:</strong> Sets up the initial block numbers for the blockchain
                    synchronization.
                </li>
                <li><strong>Sync with Blockchain:</strong> Starts the process of synchronizing the local data with the
                    blockchain.
                </li>
            </ul>
            <button onClick={() => sendControlCommand('reset_db')}>Reset Database</button>
            <button onClick={() => sendControlCommand('init_block_numbers')}>Initialize Block Numbers</button>
            <button onClick={() => sendControlCommand('sync')}>Sync with Blockchain</button>
            {isLoading && <p>Loading...</p>}
            {response && <p>Response: {response}</p>}
            {error && <p>Error: {error}</p>}
        </div>
    );
}

export default ControlPanel;