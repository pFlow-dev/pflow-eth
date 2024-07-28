import React, {Fragment} from 'react';
import PflowExplorer from "./components/explorer/PflowExplorer";
import styles from "./App.module.css";

function App() {
    return (<Fragment>
        <div className={styles.appContainer}>
            <svg className={styles.svgHeader}>
            </svg>
            <PflowExplorer/>
            <svg className={styles.svgFooter}>
            </svg>
        </div>
    </Fragment>
);
}

export default App;
