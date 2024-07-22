import * as mm from "./model";
import {FlowBuilder, ModelDeclaration, ModelType, newModel, Version} from "./model";
import {ethers, toBigInt} from "ethers";
import {MyStateMachine, MyStateMachine__factory} from "../typechain-types";
import {ModelContext} from "./api";


const initialModel = mm.newModel({
    type: mm.ModelType.petriNet,
    declaration: () => {}
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

const ScaleX = 80;
const ScaleY = 80;
const Margin = 22;

function scaleX(x: number) {
    return x * ScaleX;
}

function scaleY(y: number) {
    return y * ScaleY + Margin;
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
        return this.status;
    }

    async loadFromServer(): Promise<ModelDeclaration> {
        return this.getModelFromServer().then((modelContext) => {
            const ctx = modelContext.context;
            this.petriNet = newModel({
                declaration: (dsl: mm.Dsl) => {
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
                        place(p.Label, p.Initial, p.Capacity, scaleX(p.Position.X), scaleY(p.Position.Y));
                    });
                    ctx.Transitions.forEach((t) => {
                        transition(t.Label, t.Role, scaleX(t.Position.X), scaleY(t.Position.Y));
                        t.Delta.forEach((d, i) => {
                            if (d < 0) {
                                arc(places[i], t.Label, d);
                            }
                            if (d > 0) {
                                arc(t.Label, places[i], d);
                            }
                        })
                        t.Guard.forEach((g, i) => {
                            if (g < 0) {
                                guard(places[i], t.Label, g);
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
            this.petriNet = newModel({
                declaration: (dsl: mm.Dsl) => {
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
                        place(p.label, convert(p.initial), convert(p.capacity), scaleX(convert(p.position.x)), scaleY(convert(p.position.y)));
                    });
                    ctx.transitions.forEach((t) => {
                        transition(t.label, convert(t.role), scaleX(convert(t.position.x)), scaleY(convert(t.position.y)));
                        t.delta.forEach((d, i) => {
                            if (d < 0) {
                                arc(places[i], t.label, convert(d));
                            }
                            if (d > 0) {
                                arc(t.label, places[i], convert(d));
                            }
                        })
                        t.guard.forEach((g, i) => {
                            if (g < 0) {
                                guard(places[i], t.label, convert(g));
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