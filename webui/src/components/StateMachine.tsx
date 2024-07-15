import React, {useEffect, useState} from 'react';
import Model from "./Model";
import ControlPanel from "./ControlPanel"; // Import ControlPanel component
import {MetaModel} from "../lib/pflow";
import styles from './StateMachine.module.css';
import Signal from "./Signal";
import TransitionList from "./Transactions";
import NodeStatus from "./NodeStatus";

interface ModelData {
    // Assuming a placeholder structure for ModelData, adjust according to actual data structure
    details: any;
}

function StateMachine() {
    const [modelData, setModelData] = useState<ModelData | null>(null);
    const [isLoading, setIsLoading] = useState<boolean>(true);
    const [address, setAddress] = useState<string>('0x5fbdb2315678afecb367f032d93f642f64180aa3');
    const [activeTab, setActiveTab] = useState<string>('model');

    useEffect(() => {
        const url = `/v0/model?addr=${address}`;
        fetchData(url);
    }, []);

    const handleAddressChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setAddress(e.target.value);
    };

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        const url = `/v0/model?addr=${address}`;
        fetchData(url);
    };

    const fetchData = async (url: string) => {
        setIsLoading(true);
        try {
            const response = await fetch(url);
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            const data: ModelData = await response.json();
            setModelData(data);
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

    const metaModel = new MetaModel({editor: true});

    return (
        <div className={styles.stateMachinContainer}>
            <NodeStatus/>
            <div className={styles.tabs}>
                <button onClick={() => handleTabChange('model')}>Model</button>
                <button onClick={() => handleTabChange('signal')}>Signal</button>
                <button onClick={() => handleTabChange('transitionList')}>Transitions</button>
                <button onClick={() => handleTabChange('controlPanel')}>Control Panel</button>
                {/* Add ControlPanel tab button */}
            </div>
            <form onSubmit={handleSubmit}>
                <input type="text" size={46} value={address} onChange={handleAddressChange} placeholder="Address"/>
                <button type="submit">Load Address</button>
            </form>
            {activeTab === 'model' && <Model metaModel={metaModel}/>}
            {activeTab === 'signal' && <Signal/>}
            {activeTab === 'transitionList' && <TransitionList/>}
            {activeTab === 'controlPanel' &&
                <ControlPanel/>} {/* Render ControlPanel component when its tab is active */}
        </div>
    );
}

export default StateMachine;