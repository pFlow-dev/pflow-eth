import React, {useEffect, useState} from 'react';
import PetriNet from "./PetriNet";
import ControlPanel from "./ControlPanel"; // Import ControlPanel component
import {MetaModel} from "../lib/pflow";
import styles from './StateMachine.module.css';
import Signal from "./Signal";
import TransitionList from "./Transactions";
import PingStatus from "./PingStatus";
import {ModelContext} from "../lib/pflow/api";
import ConnectWallet from "./ConnectWallet";

// Deployed MyStateMachine
const defaultAddress = '0x5fbdb2315678afecb367f032d93f642f64180aa3';

function StateMachine() {
    const [modelData, setModelData] = useState<ModelContext | null>(null);
    const [isLoading, setIsLoading] = useState<boolean>(true);
    const [address, setAddress] = useState<string>(defaultAddress);
    const [activeTab, setActiveTab] = useState<string>('model');

    const metaModel = new MetaModel({address: defaultAddress});

    useEffect(() => {
        fetchData();
    }, []);

    const handleAddressChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        metaModel.loadFromContract().then(() => {
            // REVIEW: this is a web3 call
            // TODO: app should support usage without a backend
            setAddress(e.target.value);
        });
    };

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        await fetchData();
    };

    const fetchData = async () => {
        setIsLoading(true);
        try {
            const modelCtx: ModelContext = await metaModel.getModelFromServer();
            // console.log({modelCtx});
            setModelData(modelCtx);
        } catch (error) {
            console.error('Error fetching model data:', error);
        } finally {
            setIsLoading(false);
        }
    };

    const handleTabChange = (tabName: string) => {
        setActiveTab(tabName);
    };

    if (isLoading) {
        return <div>Loading...</div>;
    }

    if (!modelData) {
        return <div>No data available</div>;
    }


    return (
        <div className={styles.stateMachineContainer}>
            <PingStatus metaModel={metaModel}/>
            <div className={styles.tabs}>
                <button onClick={() => handleTabChange('model')}>Model</button>
                <button onClick={() => handleTabChange('signal')}>Signal</button>
                <button onClick={() => handleTabChange('transitionList')}>Transitions</button>
                <button onClick={() => handleTabChange('controlPanel')}>Control Panel</button>
                <ConnectWallet metaModel={metaModel}/>
            </div>
            <form onSubmit={handleSubmit}>
                <input type="text" size={46} value={address} onChange={handleAddressChange} placeholder="Address"/>
                <button type="submit">Load Address</button>
            </form>
            {activeTab === 'model' && <PetriNet metaModel={metaModel}/>}
            {activeTab === 'signal' && <Signal metaModel={metaModel}/>}
            {activeTab === 'transitionList' && <TransitionList/>}
            {activeTab === 'controlPanel' &&
                <ControlPanel/>}
        </div>
    );
}

export default StateMachine;