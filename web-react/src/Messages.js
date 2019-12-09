import React, {useEffect, useRef, useState} from 'react';
import DownloadFileModal from "./DownloadFileModal";

const Messages = ({instanceAddress, currentContact, messages, name}) => {
    const [currentMessage, setCurrentMessage] = useState("");
    const [showModal, setShowModal] = useState(false);

    const hideModal = () => {
        setShowModal(false);
    };

    const generalName = 'General';

    const messagesEndRef = useRef(null);

    const scrollToBottom = () => {
        messagesEndRef.current.scrollIntoView({ behavior: "smooth" })
    };

    useEffect(scrollToBottom, [messages]);

    const updateMessage = e => {
        setCurrentMessage(e.target.value);
    };

    const postMessage = e => {
        e.preventDefault();
        let message = {"Text": currentMessage};
        let route = 'rumor';
        if (currentContact !== generalName) {
            message.Destination = currentContact;
            route = 'privateMessage';
        }
        fetch(`http://${instanceAddress}/${route}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(message),
        })
        .then(() => {
            setCurrentMessage("");
        })
        .catch(e => {
            console.error(e);
            alert("Error: see logs for details");
        });
    };

    const messagesHtml = (messages, name) => messages.map((message, i) =>
        <div key={message.Origin+message.ID+'_'+i} className={'niceBox ' + ((message.Origin === name) ? 'ownRumor' : '')}>
            <div className="rumorHeader">{message.Origin}</div>
            {message.Text}
        </div>
    );

    return (
        <div id="rumors">
            <DownloadFileModal
                instanceAddress={instanceAddress}
                show={showModal}
                dismiss={hideModal}
                currentContact={currentContact} />

            <div id="chatContainer">
                <h2 className="justify-content-between align-items-center d-flex">
                    <span className="badge badge-info">{currentContact}</span>
                    {
                        currentContact !== generalName &&
                        <button type="button" className="btn btn-primary" onClick={() => setShowModal(true)}>Browse files</button>
                    }
                </h2>
                <div id="chat">
                    {messages && messagesHtml(messages, name)}
                    <div ref={messagesEndRef} />
                </div>
                <form id="rumorForm" className="input-group" onSubmit={postMessage}>
                    <div className="input-group-prepend">
                        <label htmlFor="newRumor"><span className="input-group-text">Your message:</span></label>
                    </div>
                    <input type="text" className="form-control" value={currentMessage} onChange={updateMessage} />
                    <div className="input-group-append">
                        <input type="submit" className="form-control btn btn-primary" value="Send" />
                    </div>
                </form>
            </div>
        </div>
    );
};

export default Messages;