import React from "react";
import {MetaModel} from "../../lib/pflow";
import { useWeb3ModalProvider, useWeb3ModalAccount } from '@web3modal/ethers/react'
import { BrowserProvider, Contract, formatUnits } from 'ethers'



interface TransitionProps {
    id: string;
    metaModel: MetaModel;
}

export default function Transition(props: TransitionProps) {
    const {metaModel} = props;
    const { address, chainId, isConnected } = useWeb3ModalAccount()
    const { walletProvider } = useWeb3ModalProvider()
    const provider =  walletProvider ? new BrowserProvider(walletProvider) : undefined;


    function getHandleWidth() {
        return 36;
    }

    function getStroke() {
        return "#000000";
    }

    function getFill() {
        const res = metaModel.testFire(props.id)
        if (res.ok) {
            return "#62fa75";
        }
        if (res.inhibited) {
            return "#fab5b0";
        }
        return "#ffffff";
    }

    async function onClick(evt: React.MouseEvent) {
        if (provider && isConnected) {
            await metaModel.transitionClick(provider, props.id);
        }
        evt.stopPropagation();
    }

    const t = metaModel.getTransition(props.id);

    function TextLabel() {
        return <text id={props.id + "[label]"} x={t.position.x - 15} y={t.position.y - 20}>{props.id}</text>
    }

    return (
        <g
            onClick={onClick}
            onDoubleClick={(evt) => evt.preventDefault()}
            onContextMenu={(evt) => {
                evt.preventDefault();
                evt.stopPropagation();
            }}>
            <circle id={props.id + "_handle"} cx={t.position.x} cy={t.position.y} r={getHandleWidth()}
                    fill="transparent" stroke="transparent"/>
            <rect
                className="transition" width="30" height="30" rx={4} fill={getFill()} stroke={getStroke()}
                id={props.id} x={t.position.x - 15} y={t.position.y - 15}
            />
            <TextLabel/>
        </g>
    );
};

