import {Fragment, useEffect, useState} from 'react';

interface ContractQueryProps {
}

export default function OpenContract(props: ContractQueryProps) {
    const [address, setAddress] = useState('');

    useEffect(() => {
        const params = new URLSearchParams(window.location.search);
        const addressFromUrl = params.get('address') || '';
        setAddress(addressFromUrl);
    }, [window.location]);

    async function onSubmit(evt: any) {
        evt.preventDefault();

        if (!address || !/^0x[a-fA-F0-9]{40}$/.test(address)) {
            alert("Invalid address format");
            return;
        }

        try {
            // const def = await metaModel.loadFromAddress({address});
            window.location.href = `?address=${address}`;
        } catch (e) {
            const err = e as Error;
            alert("Error importing contract: " + err.message + " -- Connected to the right network?");
            console.error(e);
        }
    }

    return (
        <Fragment>
            <form onSubmit={onSubmit}>
                <input
                    type="text"
                    value={address}
                    onChange={(e) => setAddress(e.target.value)}
                    style={{width: '30em'}}
                />
                <input type="submit" value="Open ETH Contract"/>
            </form>
        </Fragment>
    );
}