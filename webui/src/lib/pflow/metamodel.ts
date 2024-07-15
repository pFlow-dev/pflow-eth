import * as mm from "./model";

export type MaybeNode = mm.Place | mm.Transition | null

const noOp = () => {
};

const initialModel = mm.newModel({
    declaration: (dsl: mm.Dsl) => {
        const {fn, role, cell} = dsl;

        const p1 = cell('p1', 1, 1, {x: 100, y: 100});
        const p2 = cell('p2', 0, 0, {x: 200, y: 200});

        const r0 = role('default');

        const t1 = fn('t1', r0, {x: 200, y: 100});
        const t2 = fn('t2', r0, {x: 100, y: 200});
        const t3 = fn('t3', r0, {x: 350, y: 100});

        t1.tx(1, p1)
        t2.tx(1, p2)
        p1.tx(1, t2)

    },
    type: mm.ModelType.petriNet
});


interface StreamLog {
    ts: number,
    revision: number,
    action: string,
    href: string,
}

type ModelOptions = {
    editor: boolean
    superModel?: boolean
    schema?: string
}

export class MetaModel {
    m: mm.Model = initialModel;
    height: number = 600;
    logs: Map<number, StreamLog> = new Map<number, StreamLog>();
    protected superModel: boolean = false;
    protected updateHook: () => void = noOp;

    constructor(opts: ModelOptions) {
        if (opts.superModel) {
            this.superModel = true;
        }
        this.m = initialModel;
    }

    getState(): mm.Vector {
        let s = this.m.initialVector();
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

    isSuper(): boolean {
        return this.superModel;
    }

    getObj(id: string): mm.MetaObject {
        let obj: mm.MetaObject | undefined = this.m.def.transitions.get(id);
        if (obj) {
            return obj
        }
        obj = this.m.def.places.get(id);
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
        const res = this.m.fire(state, action, 1);
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