/* Autogenerated file. Do not edit manually. */
/* tslint:disable */
/* eslint-disable */

import {Contract, type ContractRunner, Interface} from "ethers";
import type {Metamodel, MetamodelInterface} from "../Metamodel";

const _abi = [
    {
        anonymous: false,
        inputs: [
            {
                indexed: true,
                internalType: "uint8",
                name: "role",
                type: "uint8",
            },
            {
                indexed: true,
                internalType: "uint8",
                name: "actionId",
                type: "uint8",
            },
            {
                indexed: true,
                internalType: "uint256",
                name: "scalar",
                type: "uint256",
            },
            {
                indexed: false,
                internalType: "uint256",
                name: "sequence",
                type: "uint256",
            },
        ],
        name: "SignaledEvent",
        type: "event",
    },
    {
        inputs: [],
        name: "context",
        outputs: [
            {
                components: [
                    {
                        internalType: "uint256",
                        name: "sequence",
                        type: "uint256",
                    },
                    {
                        internalType: "int256[]",
                        name: "state",
                        type: "int256[]",
                    },
                    {
                        components: [
                            {
                                internalType: "string",
                                name: "label",
                                type: "string",
                            },
                            {
                                internalType: "uint8",
                                name: "offset",
                                type: "uint8",
                            },
                            {
                                components: [
                                    {
                                        internalType: "uint8",
                                        name: "x",
                                        type: "uint8",
                                    },
                                    {
                                        internalType: "uint8",
                                        name: "y",
                                        type: "uint8",
                                    },
                                ],
                                internalType: "struct Model.Position",
                                name: "position",
                                type: "tuple",
                            },
                            {
                                internalType: "uint256",
                                name: "initial",
                                type: "uint256",
                            },
                            {
                                internalType: "uint256",
                                name: "capacity",
                                type: "uint256",
                            },
                        ],
                        internalType: "struct Model.Place[]",
                        name: "places",
                        type: "tuple[]",
                    },
                    {
                        components: [
                            {
                                internalType: "string",
                                name: "label",
                                type: "string",
                            },
                            {
                                internalType: "uint8",
                                name: "offset",
                                type: "uint8",
                            },
                            {
                                components: [
                                    {
                                        internalType: "uint8",
                                        name: "x",
                                        type: "uint8",
                                    },
                                    {
                                        internalType: "uint8",
                                        name: "y",
                                        type: "uint8",
                                    },
                                ],
                                internalType: "struct Model.Position",
                                name: "position",
                                type: "tuple",
                            },
                            {
                                internalType: "uint8",
                                name: "role",
                                type: "uint8",
                            },
                            {
                                internalType: "int256[]",
                                name: "delta",
                                type: "int256[]",
                            },
                            {
                                internalType: "int256[]",
                                name: "guard",
                                type: "int256[]",
                            },
                        ],
                        internalType: "struct Model.Transition[]",
                        name: "transitions",
                        type: "tuple[]",
                    },
                ],
                internalType: "struct Model.Context",
                name: "",
                type: "tuple",
            },
        ],
        stateMutability: "view",
        type: "function",
    },
    {
        inputs: [],
        name: "sequence",
        outputs: [
            {
                internalType: "uint256",
                name: "",
                type: "uint256",
            },
        ],
        stateMutability: "view",
        type: "function",
    },
    {
        inputs: [
            {
                internalType: "uint8",
                name: "action",
                type: "uint8",
            },
            {
                internalType: "uint256",
                name: "scalar",
                type: "uint256",
            },
        ],
        name: "signal",
        outputs: [],
        stateMutability: "nonpayable",
        type: "function",
    },
    {
        inputs: [
            {
                internalType: "uint8[]",
                name: "actions",
                type: "uint8[]",
            },
            {
                internalType: "uint256[]",
                name: "scalars",
                type: "uint256[]",
            },
        ],
        name: "signalMany",
        outputs: [],
        stateMutability: "nonpayable",
        type: "function",
    },
] as const;

export class Metamodel__factory {
    static readonly abi = _abi;

    static createInterface(): MetamodelInterface {
        return new Interface(_abi) as MetamodelInterface;
    }

    static connect(address: string, runner?: ContractRunner | null): Metamodel {
        return new Contract(address, _abi, runner) as unknown as Metamodel;
    }
}
