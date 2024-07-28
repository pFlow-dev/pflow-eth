import * as mm from "./model";
import {FlowBuilder, ModelDeclaration, ModelType, newModel } from "./model";
import {ethers, toBigInt} from "ethers";
import {MyStateMachine, MyStateMachine__factory} from "../typechain-types";
import {ModelContext} from "./api";
import React, {SetStateAction} from "react";


const initialModel = mm.newModel({
    type: mm.ModelType.petriNet,
    declaration: () => {}
});

type ModelOptions = {
    address?: string, // falls back to hardhat default address
    initialModel?: mm.Model,
    updateHook?: React.Dispatch<SetStateAction<MetaModel>> | null,
}

export type NodeStatus = {
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

const xScale = 80;
const yScale = 80;
const positionMargin = 22;

function scaleX(x: number) {
    return x * xScale + positionMargin;
}

function scaleY(y: number) {
    return y * yScale ;
}

export class MetaModel {
    petriNet: mm.Model = initialModel;
    height: number = 600;
    stateMachine: MyStateMachine = {} as MyStateMachine;
    address: string = '';
    provider: ethers.BrowserProvider;
    status: NodeStatus | null = null;
    updateHook: React.Dispatch<SetStateAction<MetaModel>> | null = null;

    constructor(opts?: ModelOptions) {
        if (opts && opts.updateHook) {
            this.updateHook = opts.updateHook;
        }
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

    update(): void {
        if (this.updateHook) {
            this.updateHook(new MetaModel({
                address: this.address,
                initialModel: this.petriNet,
                updateHook: this.updateHook,
            }));
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

    getSession(): string {
        return localStorage.getItem('session') || '';
    }

    setSession(session: string): void {
        localStorage.setItem('session', session)
    }

    async pollServer(): Promise<NodeStatus | null> {
        const response = await fetch('/v0/ping?session=' + this.getSession());
        if (response.ok) {
            this.status = await response.json() as NodeStatus;
        }
        return this.status;
    }

    async reloadModel(): Promise<void> {
        if (localStorage.getItem('mode') === 'ForceApi') {
            return this.loadFromServer().then(() => this.update());
        } else {
            return this.loadFromContract().then(() => this.update());
        }
    }

    async loadFromServer(): Promise<ModelDeclaration> {
        return this.getModelFromServer().then((modelContext) => {
            const ctx = modelContext.context;
            this.petriNet = newModel({
                declaration: (dsl: mm.ModelBuilder) => {
                    const places: string[] = []
                    const transitions: string[] = []
                    ctx.Places.forEach((p) => {
                        places[p.Offset] = p.Label;
                    })
                    ctx.Transitions.forEach((t) => {
                        transitions[t.Offset] = t.Label;
                    })
                    const {place, transition, arc, guard} = FlowBuilder({
                        modelDsl: dsl,
                        places,
                        transitions
                    });
                    ctx.Places.forEach((p) => {
                        // NOTE: 'initial' value in this context is set based on  latest on-chain state
                        place(p.Label, ctx.State[p.Offset], p.Capacity, scaleX(p.Position.X), scaleY(p.Position.Y));
                    });
                    ctx.Transitions.forEach((t) => {
                        transition(t.Label, t.Role, scaleX(t.Position.X), scaleY(t.Position.Y));
                        t.Delta.forEach((d, i) => {
                            if (d < 0) {
                                arc(places[i], t.Label, 0-d);
                            }
                            if (d > 0) {
                                arc(t.Label, places[i], d);
                            }
                        })
                        t.Guard.forEach((g, i) => {
                            if (g < 0) {
                                guard(places[i], t.Label, 0-g);
                            }
                            if (g > 0) {
                                guard(t.Label, places[i], g);
                            }
                        });
                    });
                },
                type: ModelType.petriNet
            });
            return this.petriNet.toObject('sparse');
        });
    }

    async signal(action: string, scalar: string): Promise<any> {
        if (localStorage.getItem('mode') === 'ForceApi') {
            return this.signalWallet(action, scalar);
        }
        return this.signalWallet(action, scalar);
    }

    async signalServer(action: string, scalar: string): Promise<any> {
        const url = `/v0/signal?addr=${this.address}&session=${this.getSession()}&action=${action}&scalar=${scalar}`;
        const response = await fetch(url);
        return await response.json();
    }

    async signalWallet(action: string, scalar: string): Promise<any> {
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

    async loadFromContract(): Promise<ModelDeclaration> {
        if (!this.stateMachine) {
            console.error('state machine contract not initialized');
            return {} as ModelDeclaration;
        }
        const convert = (n: bigint) => parseInt(n.toString());
        return this.stateMachine.context().then((ctx) => {
            console.log('context', ctx);
            this.petriNet = newModel({
                declaration: (dsl: mm.ModelBuilder) => {
                    const places:  string[] = []
                    const transitions: string[] = []
                    ctx.places.forEach((p) => {
                        places[convert(p.offset)] = p.label
                    })
                    ctx.transitions.forEach((t) => {
                        transitions[convert(t.offset)] = t.label
                    })
                    const { place, transition, arc, guard } = FlowBuilder({
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
                                arc(places[i], t.label, 0-convert(d));
                            }
                            if (d > 0) {
                                arc(t.label, places[i], convert(d));
                            }
                        })
                        t.guard.forEach((g, i) => {
                            if (g < 0) {
                                guard(places[i], t.label,0-convert(g));
                            }
                            if (g > 0) {
                                guard(t.label, places[i], convert(g));
                            }
                        });
                    });
                },
                type: ModelType.petriNet
            })

            return this.petriNet.toObject('sparse')
        })
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
        console.log('placeClick', id);
        return Promise.resolve();
    }

    placeAltClick(id: string): Promise<void> {
        console.log('placeAltClick', id);
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
        const t = this.petriNet.def.transitions.get(id)
        if (!t) {
            throw new Error('transition not found: ' + id);
        }
        return this.signal(`${t.offset}`, '1').then(() => {
            this.reloadModel().then(() => this.update());
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