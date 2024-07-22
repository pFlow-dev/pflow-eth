import React, {useEffect, useState} from "react";
import PetriNet from "./model/PetriNet";
import ControlPanel from "./ControlPanel"; // Import ControlPanel component
import {MetaModel, newModel} from "../lib/pflow";
import styles from "./PflowExplorer.module.css";
import Signal from "./Signal";
import TransitionList from "./Transactions";
import PingStatus from "./PingStatus";
import ConnectWallet from "./ConnectWallet";
import * as mm from "../lib/pflow/model";
import {FlowBuilder, ModelDeclaration, ModelType} from "../lib/pflow/model";

const defaultAddress = "0x5fbdb2315678afecb367f032d93f642f64180aa3";

function exampleDeclaration(modelDsl: mm.Dsl): void {
    const {place, transition, arc, guard} = FlowBuilder({
        modelDsl,
        places: ["foo", "bar"],
        transitions: ["increment", "decrement"]
    })
    place("foo", 1, 1, 100, 100)
    place("bar", 1, 1, 100, 200)
    transition("increment", 1, 200, 100)
    transition("decrement", 1, 200, 200)
    arc("increment", "foo", 1)
    arc("bar", "decrement", 1)
    guard("foo", "decrement", 1)
    guard("bar", "increment", 1)
}


export default function PflowExplorer() {
    const initialModel = new MetaModel({
        address: defaultAddress,
        initialModel: newModel({
            declaration: () => {},
            type: ModelType.petriNet
        }),
    });
    const [metaModel, setMetaModel] = useState<MetaModel>(initialModel);
    const [isLoading, setIsLoading] = useState<boolean>(true);
    const [address, setAddress] = useState<string>(defaultAddress);
    const [activeTab, setActiveTab] = useState<string>("model");


    useEffect(() => {
        fetchModel();
    }, []);

    const handleAddressChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    };

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        metaModel.loadFromContract().then((def) => {
            setMetaModel(new MetaModel({
                address,
                initialModel: newModel({
                    declaration: def,
                    type: ModelType.petriNet
                }),
            }));
        }).finally(() => {
            console.log("rebuild the model from the contract data");
        });
    };

    const fetchModel = async () => {
        try {
            setIsLoading(true);
            const def = await metaModel.loadFromServer()
            console.log("fetchModel", def);
            setMetaModel(new MetaModel({
                address,
                initialModel: newModel({
                    declaration: def,
                    type: ModelType.petriNet
                }),
            }));
        } catch (error) {
            console.error("Error fetching model data:", error);
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

    if (!metaModel) {
        return <div>No data available</div>;
    }


    return (
        <div className={styles.stateMachineContainer}>
            <PingStatus metaModel={metaModel}/>
            <div className={styles.tabs}>
                <button onClick={() => handleTabChange("model")}>Model</button>
                <button onClick={() => handleTabChange("signal")}>Signal</button>
                <button onClick={() => handleTabChange("transitionList")}>Transitions</button>
                <button onClick={() => handleTabChange("controlPanel")}>Control Panel</button>
                <ConnectWallet metaModel={metaModel}/>
            </div>
            <form onSubmit={handleSubmit}>
                <input type="text" size={46} value={address} onChange={handleAddressChange} placeholder="Address"/>
                <button type="submit">Load Address</button>
            </form>
            {activeTab === "model" && <PetriNet metaModel={metaModel}/>}
            {activeTab === "signal" && <Signal metaModel={metaModel}/>}
            {activeTab === "transitionList" && <TransitionList/>}
            {activeTab === "controlPanel" && <ControlPanel/>}

            <pre>
                {JSON.stringify(metaModel.petriNet.toObject("sparse"), null, 2)}
            </pre>
        </div>
    );
}