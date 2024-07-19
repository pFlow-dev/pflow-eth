// result from golang server uses capitalization (and i64 for numbers)
export interface ModelContext {
    nonce: number;
    context: {
        Sequence: number;
        State: number[];
        Places: {
            Label: string;
            Offset: number;
            Position: {
                X: number;
                Y: number;
            };
            Initial: number;
            Capacity: number;
        }[];
        Transitions: {
            Label: string;
            Offset: number;
            Position: {
                X: number;
                Y: number;
            };
            Role: number;
            Delta: number[];
            Guard: number[];
        }[];
    };
}