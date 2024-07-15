import React, {Fragment} from 'react';
import './App.css';
import StateMachine from "./components/StateMachine";

function App() {
    //return <StateMachine/>;

    return (
        <Fragment>
            <div className="appContainer">
                <StateMachine/>
            </div>
        </Fragment>
    );
}

export default App;
