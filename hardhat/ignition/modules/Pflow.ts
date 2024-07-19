import {buildModule} from "@nomicfoundation/hardhat-ignition/modules";

export default buildModule("Pflow", (m) => {
    const myStateMachine = m.contract("MyStateMachine", []);

    m.call(myStateMachine, "context", []);

    return {myStateMachine};
});
