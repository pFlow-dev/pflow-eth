import React from 'react';
import styles from './App.module.css';
import StateMachine from "./components/StateMachine";

function App() {
    return (
        <div className={styles.appContainer}>
            <StateMachine/>
        </div>
    );
}

export default App;
