import * as mm from "./model";
import {DeclarationFunction, FlowBuilder, ModelType, newModel} from "./model";
import {BigNumberish, BrowserProvider, ethers, toBigInt} from "ethers";
import {MyStateMachine, MyStateMachine__factory} from "../typechain-types";
import React, {SetStateAction} from "react";
import * as msm from "../typechain-types/MyStateMachine";
import {uint256HexToNumeric} from "./convert";

const defaultAddressesByChainId: { [chainId: number]: string } = {
    1: '0x76a7c863b91e5b6a69e7e94280149ba6737916c5', // Ethereum Mainnet
    10: '0xb746f05058697671d392dd1a4881f67dc5b006c5', // Optimism Mainnet
    1337: '0x5FbDB2315678afecb367f032d93F642f64180aa3', // Hardhat
    11155420: '0xd90e2a6aabb3415f820097be15c0801ac669827f' // Sepolia Optimism Testnet
};

const initialModel = mm.newModel({
    type: mm.ModelType.petriNet,
    declaration: () => {
    }
});

function hexValue(value: BigNumberish): string {
    return "0x" + value.toString(16);
}

function isValidEthAddress(address: string): boolean {
    const ethAddressRegex = /^0x[a-fA-F0-9]{40}$/;
    return ethAddressRegex.test(address);
}

type ModelOptions = {
    initialModel?: mm.Model,
    updateHook?: React.Dispatch<SetStateAction<MetaModel>> | null,
}

export type NodeStatus = {
    address: string;
    build: string;
    network: string;
    session_data: {
        last_ping: string;
        login_at: string;
        session_id: string;
    };
    status: string;
    time: string;
    block: string;
    sequence: string;
};

export interface EthLog {
    sequence: number;
    block_number: number;
    role: number;
    action: string; // REVIEW: should all be string?
    scalar: number;
    from_address: string;
    transaction_hash: string;
    topic_hash: string;
    removed: boolean;
}

const xScale = 80;
const yScale = 80;
const positionMargin = 22;

function scaleX(x: number) {
    return x * xScale + positionMargin;
}

function scaleY(y: number) {
    return y * yScale;
}

function convert(n: BigNumberish): number {
    return parseInt(n.toString());
}

function declarationFactory(ctx: msm.Model.HeadStruct): DeclarationFunction {
    return (dsl: mm.ModelBuilder): void => {
        const places: string[] = []
        const transitions: string[] = []
        ctx.places.forEach((p) => {
            places[convert(p.offset)] = p.label
        })
        ctx.transitions.forEach((t) => {
            transitions[convert(t.offset)] = t.label
        })
        const {place, transition, arc, guard} = FlowBuilder({
            modelDsl: dsl,
            places,
            transitions
        });
        ctx.places.forEach((p) => {
            // NOTE: 'initial' value in this context is set based on  latest on-chain state
            place(p.label, convert(ctx.state[convert(p.offset)]), convert(p.capacity), scaleX(convert(p.position.x)), scaleY(convert(p.position.y)));
        });
        ctx.transitions.forEach((t) => {
            transition(t.label, convert(t.role), scaleX(convert(t.position.x)), scaleY(convert(t.position.y)));
            t.delta.forEach((d, i) => {
                if (d < 0) {
                    arc(places[i], t.label, 0 - convert(d));
                }
                if (d > 0) {
                    arc(t.label, places[i], convert(d));
                }
            })
            t.guard.forEach((g, i) => {
                if (g < 0) {
                    guard(places[i], t.label, 0 - convert(g));
                }
                if (g > 0) {
                    guard(t.label, places[i], convert(g));
                }
            });
        });
    };
}

function getAddressByNetwork(networkId: bigint): string {
    return defaultAddressesByChainId[parseInt(networkId.toString(10))];
}

export class MetaModel {
    petriNet: mm.Model = initialModel;
    height: number = 600;
    status: NodeStatus | null = null;
    updateHook: React.Dispatch<SetStateAction<MetaModel>> | null = null;
    sequence: bigint = -1n;
    contractState: bigint[] = [];
    contractHead: msm.Model.HeadStruct = {} as msm.Model.HeadStruct;

    constructor(opts?: ModelOptions) {
        if (opts && opts.updateHook) {
            this.updateHook = opts.updateHook;
        }
        if (opts && opts.initialModel) {
            this.petriNet = opts.initialModel;
        } else {
            this.petriNet = initialModel;
        }
    }

    async getStateMachine(provider: ethers.BrowserProvider): Promise<MyStateMachine> {
        return provider.getNetwork().then(async (network) => {
            let address = new URLSearchParams(window.location.search).get('address');
            if (address && !isValidEthAddress(address)) {
                throw new Error("Invalid Ethereum address");
            }
            address = address || getAddressByNetwork(network.chainId);
            const signer = await provider.getSigner();
            return MyStateMachine__factory.connect(address, signer);
        });
    }

    update(): void {
        if (this.updateHook) {
            const mm = new MetaModel({
                initialModel: this.petriNet,
                updateHook: this.updateHook,
            })
            // REVIEW: should this be in constructor?
            mm.contractState = this.contractState;
            mm.sequence = this.sequence;
            mm.contractHead = this.contractHead;
            this.updateHook(mm);
        }
    }

    async contract(provider: ethers.BrowserProvider): Promise<MyStateMachine> {
        // REVIEW can this happen during construction instead?
        return this.getStateMachine(provider).then(async (sm) => {
            return sm.connect(await provider.getSigner());
        });
    }

    async signal(provider: BrowserProvider, action: string, scalar: string): Promise<any> {
        try {
            const contract = await this.contract(provider);
            let tx;

            if (action.includes(',') || scalar.includes(',')) {
                const actionsArray = action.split(',').map(a => toBigInt(a.trim()));
                const scalarsArray = scalar.split(',').map(s => toBigInt(s.trim()));

                if (actionsArray.length !== scalarsArray.length) {
                    throw new Error('Actions and scalars arrays must be of the same length');
                }

                tx = await contract.signalMany(actionsArray, scalarsArray);
            } else {
                tx = await contract.signal(toBigInt(action), toBigInt(scalar));
            }

            return await tx.wait();
        } catch (err) {
            return err;
        }
    }

    async loadFromContract(provider: BrowserProvider): Promise<void> {
        const sm = await this.getStateMachine(provider);
        console.info('loadFromContract');
        return sm.context()
            .then((ctx) => {
                const {state, sequence, latestBlocks} = ctx;
                if (this.sequence !== sequence) {
                    this.sequence = sequence;
                    this.contractState = state;
                    console.log(
                        `sequence: ${sequence.toString(10)} \n` +
                        `state: ${state.map((s) => s.toString(10)).join(',')} \n` +
                        `latestBlocks: ${latestBlocks.map((b) => b.toString(10)).join(',')} `
                    );
                    this.contractHead = ctx;
                }

                this.petriNet = newModel({
                    declaration: declarationFactory(ctx),
                    type: ModelType.petriNet
                })
            })
            .catch((err) => {
                console.error('failed to get context', err)
            })

    }

    getEventLogBlocks(): BigNumberish[] {
        if (!this.contractHead.latestBlocks || this.contractHead.latestBlocks.length === 0) {
            console.error('latestBlocks not set');
            return [];
        }

        return this.contractHead.latestBlocks.filter(block => block > 0);
    }

    async fetchEthLogs(provider: BrowserProvider): Promise<EthLog[]> {
        try {
            const blocks = this.getEventLogBlocks();
            let allLogs: EthLog[] = [];
            const networkId = await provider.getNetwork();
            const address = getAddressByNetwork(networkId.chainId);

            for (const block of blocks) {
                const filter = {
                    fromBlock: hexValue(block),
                    toBlock: hexValue(block),
                    address,
                    topics: [] // Add any topics if needed
                };

                const rawLogs = await provider.send("eth_getLogs", [filter]);
                const formattedLogs = rawLogs.map((log: any) => ({
                    sequence: uint256HexToNumeric(log.data).toString(10),
                    block_number: parseInt(log.blockNumber, 16),
                    role: uint256HexToNumeric(log.topics[1]).toString(10),
                    action: uint256HexToNumeric(log.topics[2]).toString(10),
                    scalar: uint256HexToNumeric(log.topics[3]).toString(10),
                    transaction_hash: log.transactionHash,
                    topic_hash: log.topics[0],
                    removed: log.removed
                }));

                allLogs = allLogs.concat(formattedLogs);
            }

            return allLogs;
        } catch (error) {
            console.error("Error fetching logs:", error);
        }
        return [];
    }

    getState(): mm.Vector {
        let s = this.petriNet.initialVector();
        return [...s]
    }

    getTokenCount(id: string): number {
        const state = this.getState();
        const n = this.getNode(id);
        if (n.metaType === 'place') {
            const p = n as mm.Place;
            return state[p.offset];
        }
        return 0;
    }

    getObj(id: string): mm.MetaObject {
        let obj: mm.MetaObject | undefined = this.petriNet.def.transitions.get(id);
        if (obj) {
            return obj
        }
        obj = this.petriNet.def.places.get(id);
        if (obj) {
            return obj
        }
        throw new Error("object not found: " + id);
    }

    getNode(id: string): mm.Place | mm.Transition {
        const obj = this.getObj(id)
        if (!obj) {
            throw new Error('Failed to select node' + id);
        }
        if (obj.metaType === 'arc') {
            throw new Error('cannot select arc as node' + id);
        }
        return obj;
    }

    placeClick(id: string): Promise<void> {
        return Promise.resolve();
    }

    placeAltClick(id: string): Promise<void> {
        return Promise.resolve();
    }

    getPlace(id: string): mm.Place {
        const obj = this.getObj(id);
        if (obj.metaType === 'place') {
            return obj;
        }
        throw new Error('not a place: ' + id)
    }

    testFire(action: string): { ok: boolean; inhibited: boolean } {
        let state = this.getState();
        const res = this.petriNet.fire(state, action, 1);
        return {ok: res.ok, inhibited: !!res.inhibited};
    }


    transitionClick(provider: BrowserProvider, id: string): Promise<void> {
        const t = this.petriNet.def.transitions.get(id)
        if (!t) {
            throw new Error('transition not found: ' + id);
        }
        return this.signal(provider, `${t.offset}`, '1').then(() => {
            return this.loadFromContract(provider).then(() => this.update());
        });
    }

    getTransition(id: string): mm.Transition {
        const obj = this.getObj(id);
        if (obj.metaType === 'transition') {
            return obj;
        }
        throw new Error('not a transition: ' + id)
    }

    setRole(role: string): void {
        localStorage.setItem('role', role);
    }

    getRoles(): string[] {
        const roles: string[] = [];
        this.petriNet.def.transitions.forEach((t) => {
            if (!roles.includes(t.role.label)) {
                roles.push(t.role.label);
            }
        })
        return roles
    }

    arcClick(id: number): Promise<void> {
        return Promise.resolve();
    }

    arcAltClick(id: number): Promise<void> {
        return Promise.resolve();
    }

}