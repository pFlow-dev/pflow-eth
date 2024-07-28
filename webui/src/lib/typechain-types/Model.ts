/* Autogenerated file. Do not edit manually. */
/* tslint:disable */
/* eslint-disable */
import type {
    BaseContract,
    BigNumberish,
    ContractMethod,
    ContractRunner,
    EventFragment,
    FunctionFragment,
    Interface,
    Listener,
} from "ethers";
import type {
    TypedContractEvent,
    TypedDeferredTopicFilter,
    TypedEventLog,
    TypedListener,
    TypedLogDescription,
} from "./common";

export interface ModelInterface extends Interface {
    getEvent(nameOrSignatureOrTopic: "SignaledEvent"): EventFragment;
}

export namespace SignaledEventEvent {
    export type InputTuple = [
        role: BigNumberish,
        actionId: BigNumberish,
        scalar: BigNumberish,
        sequence: BigNumberish
    ];
    export type OutputTuple = [
        role: bigint,
        actionId: bigint,
        scalar: bigint,
        sequence: bigint
    ];

    export interface OutputObject {
        role: bigint;
        actionId: bigint;
        scalar: bigint;
        sequence: bigint;
    }

    export type Event = TypedContractEvent<InputTuple, OutputTuple, OutputObject>;
    export type Filter = TypedDeferredTopicFilter<Event>;
    export type Log = TypedEventLog<Event>;
    export type LogDescription = TypedLogDescription<Event>;
}

export interface Model extends BaseContract {
    interface: ModelInterface;
    filters: {
        "SignaledEvent(uint8,uint8,uint256,uint256)": TypedContractEvent<
            SignaledEventEvent.InputTuple,
            SignaledEventEvent.OutputTuple,
            SignaledEventEvent.OutputObject
        >;
        SignaledEvent: TypedContractEvent<
            SignaledEventEvent.InputTuple,
            SignaledEventEvent.OutputTuple,
            SignaledEventEvent.OutputObject
        >;
    };

    connect(runner?: ContractRunner | null): Model;

    waitForDeployment(): Promise<this>;

    queryFilter<TCEvent extends TypedContractEvent>(
        event: TCEvent,
        fromBlockOrBlockhash?: string | number | undefined,
        toBlock?: string | number | undefined
    ): Promise<Array<TypedEventLog<TCEvent>>>;

    queryFilter<TCEvent extends TypedContractEvent>(
        filter: TypedDeferredTopicFilter<TCEvent>,
        fromBlockOrBlockhash?: string | number | undefined,
        toBlock?: string | number | undefined
    ): Promise<Array<TypedEventLog<TCEvent>>>;

    on<TCEvent extends TypedContractEvent>(
        event: TCEvent,
        listener: TypedListener<TCEvent>
    ): Promise<this>;

    on<TCEvent extends TypedContractEvent>(
        filter: TypedDeferredTopicFilter<TCEvent>,
        listener: TypedListener<TCEvent>
    ): Promise<this>;

    once<TCEvent extends TypedContractEvent>(
        event: TCEvent,
        listener: TypedListener<TCEvent>
    ): Promise<this>;

    once<TCEvent extends TypedContractEvent>(
        filter: TypedDeferredTopicFilter<TCEvent>,
        listener: TypedListener<TCEvent>
    ): Promise<this>;

    listeners<TCEvent extends TypedContractEvent>(
        event: TCEvent
    ): Promise<Array<TypedListener<TCEvent>>>;

    listeners(eventName?: string): Promise<Array<Listener>>;

    removeAllListeners<TCEvent extends TypedContractEvent>(
        event?: TCEvent
    ): Promise<this>;

    getFunction<T extends ContractMethod = ContractMethod>(
        key: string | FunctionFragment
    ): T;

    getEvent(
        key: "SignaledEvent"
    ): TypedContractEvent<
        SignaledEventEvent.InputTuple,
        SignaledEventEvent.OutputTuple,
        SignaledEventEvent.OutputObject
    >;
}
