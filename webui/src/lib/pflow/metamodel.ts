import * as mm from "./model";
import {ModelType, Version} from "./model";
import {ethers, toBigInt} from "ethers";
import {MyStateMachine, MyStateMachine__factory} from "../typechain-types";
import {ModelContext} from "./api";

const initialModel = mm.newModel({
    declaration: (dsl: mm.Dsl) => {
        // empty model
    },
    type: mm.ModelType.petriNet
});

type ModelOptions = {
    address?: string, // falls back to hardhat default address
    initialModel?: mm.Model
}

export interface NodeStatus {
    time: string;
    block: number;
    sequence: number;
    job_timer: string;
    blocks_added: number;
    unconfirmed_tx: number;
    js_build: string;
    css_build: string;
}

export class MetaModel {
    petriNet: mm.Model = initialModel;
    height: number = 600;
    stateMachine: MyStateMachine = {} as MyStateMachine;
    address: string = '';
    provider: ethers.BrowserProvider;
    status: NodeStatus | null = null;

    constructor(opts?: ModelOptions) {
        if (opts && opts.initialModel) {
            this.petriNet = opts.initialModel;
        } else {
            this.petriNet = initialModel;
        }
        this.provider = new ethers.BrowserProvider(window.ethereum);
        if (opts?.address) {
            this.address = opts.address;
            this.stateMachine = MyStateMachine__factory.connect(opts.address, this.provider)
        }
    }

    async contract(): Promise<MyStateMachine> {
        // REVIEW can this happen during construction instead?
        return this.stateMachine.connect(await this.provider.getSigner());
    }

    async getModelFromServer(): Promise<ModelContext> {
        const url = `/v0/model?addr=${this.address}`;
        const response = await fetch(url);
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return await response.json() as ModelContext;
    }

    async pollServer(): Promise<NodeStatus | null> {
        const response = await fetch('/v0/ping');
        if (response.ok) {
            this.status = await response.json() as NodeStatus;
        }
        // TODO: add hook here to detect changes in contract status so we can update the UI
        return this.status;
    }

    async loadFromServer() {
        const modelCtx = await this.getModelFromServer();
        const {context} = modelCtx;
        const {Places, Transitions} = context;
        const def = {
            version: "v0" as Version,
            modelType: ModelType.petriNet,
            places: {},
            transitions: {},
            arcs: []
        }
        // TODO: return a declaration object

        console.log({context})
    }

async signal(action: string, scalar: string): Promise<any> {
    if (!this.stateMachine) {
        return {error: 'state machine contract not initialized'};
    }

    try {
        const contract = await this.contract();
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

    async loadFromContract() {
        if (!this.stateMachine) {
            console.error('state machine contract not initialized');
            return
        }
        const modelCtx = await this.stateMachine.context();

        const {places, transitions, state, sequence} = modelCtx;
        const def = {
            version: "v0" as Version,
            modelType: ModelType.petriNet,
            places: {},
            transitions: {},
            arcs: []
        }
        places.forEach((place) => {
            const {label, offset, initial, capacity, position} = place;
            const {x, y} = position;
            def.places = {
                [label]: {
                    label,
                    offset: Number(offset),
                    initial: Number(initial),
                    capacity: Number(capacity),
                    x: Number(x),
                    y: Number(y)
                },
                ...def.places
            };
        })
        transitions.forEach((transition) => {
            const {label, offset, role, position, delta, guard} = transition;
            const {x, y} = position;
            def.transitions = {
                [label]: {
                    label,
                    offset: Number(offset),
                    role: Number(role),
                    x: Number(x),
                    y: Number(y),
                    delta: delta.map((d) => Number(d)),
                    guard: guard.map((g) => Number(g))
                }
            }
        })
        console.log({def})
        // FIXME: add arcs/guards and install the new model
        return Promise.resolve(def)
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


    transitionClick(id: string): Promise<void> {
        return Promise.resolve();
    }

    getTransition(id: string): mm.Transition {
        const obj = this.getObj(id);
        if (obj.metaType === 'transition') {
            return obj;
        }
        throw new Error('not a transition: ' + id)
    }

    arcClick(id: number): Promise<void> {
        return Promise.resolve();
    }

    arcAltClick(id: number): Promise<void> {
        return Promise.resolve();
    }

}