import React from "react";
import {MetaModel} from "../../lib/pflow";

interface PlaceProps {
    id: string;
    metaModel: MetaModel;
}

// export default class Place extends React.Component<PlaceProps, NodeState> {
export default function Place(props: PlaceProps) {

    const {metaModel} = props;

    function renderTokens(p: { x: number; y: number; }) {
        let tokens = metaModel.getTokenCount(props.id);

        let tokenColor = "#0000007x";

        if (tokens === 0) {
            return; // don"t show zeros
        }
        if (tokens === 1) {
            return (<circle cx={p.x} cy={p.y} r="2" id={props.id + "_tokens"} fill={tokenColor} stroke={tokenColor}
                            orient="0" className="tokens"/>);
        }
        if (tokens < 10) {
            return (<text id={props.id + "_tokens"} x={p.x - 4} y={p.y + 5} fill={tokenColor} stroke={tokenColor}
                          className="large">{tokens}</text>);
        }
        if (tokens < 100) {
            return (<text id={props.id + "_tokens"} x={p.x - 7} y={p.y + 5} fill={tokenColor} stroke={tokenColor}
                          className="small">{tokens}</text>);
        }
        if (tokens < 1_000) {
            return (<text id={props.id + "_tokens"} x={p.x - 10} y={p.y + 5} fill={tokenColor} stroke={tokenColor}
                          className="small">{tokens}</text>);
        }
        if (tokens < 10_000) {
            return (<text id={props.id + "_tokens"} x={p.x - 14} y={p.y + 5} fill={tokenColor} stroke={tokenColor}
                          className="small">{tokens}</text>);
        }
        return (<g transform="">
            <text id={props.id + "_tokens"} x={p.x - 14} y={p.y + 5}
                  fill={tokenColor} stroke={tokenColor}
                  style={{
                      fontSize: "10px",
                      fontFamily: "Arial",
                  }}>{tokens}</text>
        </g>);
    }

    function getHandleWidth() {
        return 36;
    }

    function getStroke() {
        return "#000000";
    }

    async function onClick(evt: React.MouseEvent) {
        await metaModel.placeClick(props.id);
        evt.stopPropagation();
    }

    async function onAltClick(evt: React.MouseEvent) {
        await metaModel.placeAltClick(props.id);
        evt.preventDefault();
        evt.stopPropagation();
    }

    const p = metaModel.getPlace(props.id).position;
    let fill = "#FFFFFF";

    function TextLabel() {
        return <text id={props.id + "[label]"} x={p.x - 18} y={p.y - 20}>{props.id}</text>
    }

    return (
        <g
            onClick={onClick}
            onContextMenu={onAltClick}>

            <circle id={props.id + "_handle"} cx={p.x} cy={p.y} r={getHandleWidth()} fill="transparent"
                    stroke="transparent"/>
            <circle cx={p.x} cy={p.y} r="16" id={props.id}
                    strokeWidth="1.5" fill={fill} stroke={getStroke()} orient="0"
                    className="place"
                    shapeRendering="auto"
            />
            {renderTokens(p)}
            <TextLabel/>
        </g>
    );
};
