import React, {ReactElement} from "react";
import Place from "./Place";
import Arc from "./Arc";
import Transition from "./Transition";
import {MetaModel} from "../../lib/pflow";
import * as mm from "../../lib/pflow/model";

interface ModelProps {
    metaModel: MetaModel;
}

export default function PetriNet(props: ModelProps) {
    const {metaModel} = props;
    const {places, transitions} = metaModel.petriNet.def;

    const placeElements = Array.from(places.keys()).map((label) =>
        <Place key={label} id={label} metaModel={metaModel}/>,
    );

    const transitionElements = Array.from(transitions.keys()).map((label) =>
        <Transition key={label} id={label} metaModel={metaModel}/>,
    );

    const arcs: ReactElement[] = metaModel.petriNet.def.arcs.map((arc: mm.Arc, index: number) => {
        const source = arc.source.place || arc.source.transition;
        const target = arc.target.place || arc.target.transition;
        if (!source || !target) {
            return <React.Fragment key={index}/>;
        }
        const id = index.toString() + "_" + source.label + "_" + target.label;
        return <Arc key={id} id={id} metaModel={metaModel} arc={arc}/>;
    });

    return (
        <svg id="pflow-svg" onContextMenu={(evt) => evt.preventDefault()}>
            <defs>
                <marker id="markerArrow1" markerWidth="23" markerHeight="13" refX="31" refY="6" orient="auto">
                    <rect className="arrowSpace1" width="28" height="3" fill="#ffffff" stroke="#ffffff" x="3"
                          y="5"/>
                    <path d="M2,2 L2,11 L10,6 L2,2"/>
                </marker>
                <marker id="markerInhibit1" markerWidth="23" markerHeight="13" refX="31" refY="6" orient="auto">
                    <rect className="inhibitSpace1" width="28" height="3" fill="#ffffff" stroke="#ffffff" x="3"
                          y="5"/>
                    <circle cx="5" cy="6.5" r={4}/>
                </marker>
            </defs>
            <g id={"petri-net"} key={"petri-net"}>
                {arcs}
                {placeElements}
                {transitionElements}
            </g>
        </svg>
    )
        ;
}