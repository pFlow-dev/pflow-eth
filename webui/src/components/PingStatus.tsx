import React, {useEffect, useState} from 'react';
import styles from './PingStatus.module.css';
import {MetaModel, NodeStatus} from "../lib/pflow";

type PingStatusProps = {
    metaModel: MetaModel;
}

function PingStatus({metaModel}: PingStatusProps) {
    const [nodeData, setNodeData] = useState<NodeStatus | null>(null);
    const [isLoading, setIsLoading] = useState<boolean>(true);

    const fetchData = async () => {
        setIsLoading(true);
        try {
            metaModel.pollServer().then((data) => {
                setNodeData(data);
            });
        } catch (error) {
            console.error('Error fetching node data:', error);
        } finally {
            setIsLoading(false);
        }
    };

    useEffect(() => {
        fetchData(); // Initial fetch
        const interval = setInterval(fetchData, 30000); // Set up the interval for refreshing data every 30 seconds
        return () => clearInterval(interval); // Clear the interval when the component unmounts
    }, []);

    if (isLoading) {
        return <div>Loading...</div>;
    }

    if (!nodeData) {
        return <div>No data available</div>;
    }

    return (
        <div className={styles.tableContainer}>
            <svg
                className={styles.pflowLogoSvg}
                width="30" height="30" viewBox="0 0 50 50" fill="none" xmlns="http://www.w3.org/2000/svg">
                <g transform="scale(0.1, 0.1) ">
                    <path
                        d="M470.811 366.244V340.103H496.952V287.822H470.811V261.681H340.108V235.54H470.811V209.4H496.952V157.118H470.811V130.978H366.249V157.118H340.108V183.259H313.967V157.118H340.108V130.978H366.249V26.4385H340.108V0.297836H287.827V26.4385H261.686V157.142H235.545V26.4385H209.405V0.297836H157.123V26.4385H130.983V131.001H157.123V157.142H183.264V183.282H157.123V157.142H130.983V131.001H26.42V157.142H0.279297V209.423H26.42V235.564H157.123V261.704H26.42V287.845H0.279297V340.126H26.42V366.267H130.983V340.126H157.123V313.986H183.264V340.126H157.123V366.267H130.983V470.83H157.123V496.97H209.405V470.83H235.545V340.126H261.686V470.83H287.827V496.97H340.108V470.83H366.249V366.267H340.108V340.126H313.967V313.986H340.108V340.126H366.249V366.267H470.811V366.244ZM287.85 287.822V313.962H209.428V287.822H183.287V209.4H209.428V183.259H287.85V209.4H313.99V287.822H287.85Z"
                        fill="black"/>
                </g>
            </svg>

            <table className={styles.table}>
                <thead>
                <tr>
                    <th></th>
                    <th>Sync Timestamp</th>
                    <th>Block</th>
                    <th>Seq</th>
                    <th>Blocks Added</th>
                    <th>Unconfirmed TX</th>
                </tr>
                </thead>
                <tbody>
                <tr>
                    <td>
                    </td>
                    <td>{nodeData.time}</td>
                    <td>{nodeData.block}</td>
                    <td>{nodeData.sequence}</td>
                    <td>{nodeData.blocks_added}</td>
                    <td>{nodeData.unconfirmed_tx}</td>
                </tr>
                </tbody>
            </table>
        </div>
    );
}

export default PingStatus;