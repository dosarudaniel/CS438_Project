import React, { useState } from 'react';

const Peers = ({instanceAddress, peers}) => {
    const [newPeer, setNewPeer] = useState("");

    const updateNewPeer = e => {
        setNewPeer(e.target.value);
    }

    const postNewPeer = e => {
        e.preventDefault();
        fetch(`http://${instanceAddress}/peer`, {
            method: 'POST',
            headers: {
                'Content-Type': 'text/plain',
            },
            body: newPeer,
        })
        .then(() => setNewPeer(""))
        .catch(e => {
            console.error(e);
            alert("Error: see logs for details");
        });
    }

    return (
        <div id="peers">
            <div className="niceBox">
                <h3>Peers</h3>
                <ul className="list-group" id="peerList">
                    {peers.map(peer => <li className="list-group-item" key={peer}>{peer}</li>)}
                </ul>
                <form id="peerForm" className="input-group" onSubmit={postNewPeer}>
                    <label htmlFor="newPeer"></label>
                    <input type="text" className="form-control" id="newPeer" value={newPeer} onChange={updateNewPeer} placeholder="127.0.0.1:5000"/>
                    <div className="input-group-append">
                        <input type="submit" className="form-control btn btn-primary" value="Add" />
                    </div>
                </form>
            </div>
        </div>
    );
}

export default Peers;