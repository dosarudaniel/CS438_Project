import React, {useEffect, useState} from 'react'
import { FaSync } from "react-icons/fa";

const GossiperSelector = ({loginState, setGossiperInstance}) => {
    const gossiperEntryClasses =
        "h3 list-group-item list-group-item-action d-flex justify-content-between align-items-center";

    const [pendingInstances, setPendingInstances] = useState(0);
    const [gossiperInstances, setGossiperInstances] = useState([]);
    const [unreachableInstances, setUnreachableInstances] = useState([]);

    useEffect(() => {
        if (loginState === true) {
            refreshInstances();
        }
    }, [loginState]);

    const refreshInstances = () => searchInstances(8080, 10);

    const searchInstances = (firstPost, qty) => {
        setGossiperInstances([]);
        setUnreachableInstances([]);
        setPendingInstances(qty);
        for (let i = 0; i < qty; i++) {
            searchInstance("127.0.0.1:" + (firstPost + i))
                .then(() => setPendingInstances(state => state - 1));
        }
    };

    const searchInstance = async (address) => {
        console.log('Fetching' + address);
        fetch(`http://${address}/settings`)
            .then(response => response.json())
            .then(data => {
                console.log(data);
                if (data.hasOwnProperty('Name') && data.hasOwnProperty('GossipAddr')) {
                    setGossiperInstances(state => state.concat({
                        name: data.Name,
                        gossipAddr: data.GossipAddr,
                        uiAddress: address,
                    }));
                } else {
                    setUnreachableInstances(state => state.concat({
                        address: address,
                        reason: "--ERROR--",
                    }));
                }
            })
            .catch(() => {
                setUnreachableInstances(state => state.concat({
                    address: address,
                    reason: "DOWN",
                }));
            });
    };

    const searchingMessage = () => {
        return (
            <li className={gossiperEntryClasses + " disabled"}>
                <div className="spinner-border" role="status">
                    <span className="sr-only">Loading...</span>
                </div>
                Searching for instances...
            </li>
        )
    };

    const instances = () => {
        return gossiperInstances.sort().map(instance => (
            <button
                type="button"
                key={instance.gossipAddr}
                className={gossiperEntryClasses}
                onClick={() => setGossiperInstance(instance.uiAddress)}>
                {instance.name} <span className="badge badge-secondary">{instance.gossipAddr}</span>
            </button>
        ))
    };

    const downInstances = () => {
        return unreachableInstances.sort().map(instance => (
            <button
                type="button"
                key={instance.address}
                className={gossiperEntryClasses + " disabled"}>
                {instance.address} <span className="badge badge-danger">{instance.reason}</span>
            </button>
        ))
    };

    const refreshButton = () => (
        <button
            type="button"
            className="btn btn-info"
            onClick={refreshInstances}>
            <FaSync />
        </button>
    );

    return (
        <form className="GossiperSelector">
            <h1 className="header">Select a gossiper {pendingInstances === 0 && refreshButton()}</h1>
            <div className="list-group">
                {pendingInstances > 0 && searchingMessage()}
                {instances()}
                {downInstances()}
            </div>
        </form>
    )
};

export default GossiperSelector