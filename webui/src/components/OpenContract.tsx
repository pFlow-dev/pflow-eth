import {Fragment} from 'react';

interface ContractQueryProps {
}

export default function OpenContract(props: ContractQueryProps) {

    async function onSubmit(evt: any) {
        evt.preventDefault();

        try {
            // const def = await metaModel.loadFromAddress({address});
        }
        catch (e) {
            const err = e as Error;
            alert("Error importing contract: "+err.message+" -- Connected to the right network?")
            console.error(e)
        }
    }

    return <Fragment>
        <form onSubmit={onSubmit} >
            <input type="text" defaultValue={"0x5FbDB2315678afecb367f032d93F642f64180aa3"} style={{width: "30em"}}/>
            <input type="submit" value="Open ETH Contract"/>
        </form>
    </Fragment>
}
