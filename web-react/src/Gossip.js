import React, {useEffect, useState} from 'react'
import Contacts from './Contacts';
import Peers from './Peers';
import Messages from './Messages';
import Helpers from './Helpers';
import {FaWindowClose} from "react-icons/fa";
import SharedFiles from "./SharedFiles";

const Gossip = ({instanceAddress, quitGossiperInstance}) => {
    const generalName = 'General';

    const [notResponding, setNotResponding] = useState(true);
    const [settings, setSettings] = useState({});
    const [contacts, setContacts] = useState(new Map());
    const [currentContact, setCurrentContact] = useState(generalName);
    const [peers, setPeers] = useState([]);
    const [messages, setMessages] = useState(new Map([[generalName, []]]));
    const [currentMessages, setCurrentMessages] = useState([]);
    const [generalUnread, setGeneralUnread] = useState(0);
    const [privateMessageOffset, setPrivateMessageOffset] = useState(0);
    const [sharedFiles, setSharedFiles] = useState([]);

    useEffect(() => {
        setNotResponding(true);
        setSettings({});
        setContacts(new Map());
        setCurrentContact(generalName);
        setPeers([]);
        setMessages(new Map([[generalName, []]]));
        setCurrentMessages([]);
        setGeneralUnread(0);
        setPrivateMessageOffset(0);
        setSharedFiles([]);
    }, [instanceAddress]);

    useEffect(() => {
        console.log("Current contact: ", currentContact);
        console.log("Contacts state: ", contacts);
        if (messages.has(currentContact)) {
            setCurrentMessages(messages.get(currentContact));
            if (currentContact === generalName) {
                setGeneralUnread(0);
            } else {
                setContacts(state => {
                    state.set(currentContact, {name: contacts.get(currentContact).name, unread: 0});
                    return state;
                })
            }
        } else {
            setCurrentMessages([]);
        }
    }, [currentContact]);

    useEffect(() => {
        console.log("Current contact: ", currentContact);
        if (messages.has(currentContact)) {
            setCurrentMessages(messages.get(currentContact));
        } else {
            setCurrentMessages([]);
        }
    }, [messages]);

    Helpers.useInterval(() => {
        if (notResponding) {
            getSettings();
        } else {
            refreshPeers();
            refreshMessages();
            refreshContacts();
            refreshPrivateMessages();
            refreshSharedFiles();
        }
    }, 1000);

    const notRespondingAlert = () => {
        return (
            <div className="alert alert-danger" id="alertServerNotResponding">
                <strong>Error: </strong>Server is not responding
            </div>
        );
    };

    const getSettings = () => {
        fetch(`http://${instanceAddress}/settings`)
            .then(response => response.json())
            .then(data => {
                console.log("Settings: ", data);
                if (notResponding) {
                    setMessages(new Map([[generalName, []]]));
                    setPeers([]);
                    setPrivateMessageOffset(0);
                }
                setSettings(data);
                setNotResponding(false);
            })
            .catch(e => {
                console.error(e);
                setNotResponding(true);
            });
    };

    const refreshPeers = () => {
        fetch(`http://${instanceAddress}/peers`)
            .then(response => response.json())
            .then(data => {
                console.log("Peers: ", data);
                setPeers(data);
                setNotResponding(false);
            })
            .catch(e => {
                console.error(e);
                setNotResponding(true);
            });
    };

    const refreshContacts = () => {
        fetch(`http://${instanceAddress}/contacts`)
            .then(response => response.json())
            .then(data => {
                console.log("Contacts: ", data);
                setContacts(state => {
                    data.forEach(contact => {
                        if (contact !== settings.Name && !state.has('user' + contact)) {
                            state.set('user' + contact, {name: contact, unread: 0});
                        }
                    });
                    return state;
                });
                setNotResponding(false);
            })
            .catch(e => {
                console.error(e);
                setNotResponding(true);
            });
    };

    const setNewMessages = (contactName, newMessages) => {
        if (currentContact !== contactName) {
            console.log(`Adding messages for non-active contact ${contactName}: `, newMessages);
            if (contactName === generalName) {
                setGeneralUnread(state => state + newMessages.length)
            } else {
                setContacts(state => {
                    if (!state.has(contactName)) {
                        const message = newMessages[0];
                        const name = (message.Origin === settings.Name ? message.Destination : message.Origin);
                        state.set(contactName, {name: name, unread: newMessages.length});
                    } else {
                        const contact = state.get(contactName);
                        state.set(contactName, {name: contact.name, unread: contact.unread + newMessages.length});
                    }
                    return state;
                });
            }
        }

        setMessages(state => {
            if (state.has(contactName)) {
                state.set(contactName, state.get(contactName).concat(newMessages));
            } else {
                state.set(contactName, newMessages);
            }
            newMessages = []; // FIXME Ugly hack because of mysterious double call
            if (currentContact === contactName) {
                setCurrentMessages(state.get(contactName));
            }
            return state;
        });
    };

    const refreshMessages = () => {
        console.log("message state: ", messages);
        fetch(`http://${instanceAddress}/rumors/${messages.get(generalName).length}`)
            .then(response => response.json())
            .then(newMessages => {
                console.log("Messages: ", newMessages);
                if (Array.isArray(newMessages) && newMessages.length > 0) {
                    console.log("messages state: ", messages);
                    setNewMessages(generalName, newMessages);
                }
                setNotResponding(false);
            })
            .catch(e => {
                console.error(e);
                setNotResponding(true);
            });
    };

    const refreshPrivateMessages = () => {
        fetch(`http://${instanceAddress}/privateMessages/${privateMessageOffset}`)
            .then(response => response.json())
            .then(newMessages => {
                console.log("New private messages: ", newMessages);
                if (Array.isArray(newMessages) && newMessages.length > 0) {
                    console.log("messages state: ", messages);
                    setPrivateMessageOffset(state => state + newMessages.length);
                    newMessages.forEach(message => {
                        const contactName = 'user' + (message.Origin === settings.Name ? message.Destination : message.Origin);
                        setNewMessages(contactName, [message]);
                    });
                }
                setNotResponding(false);
            })
            .catch(e => {
                console.error(e);
                setNotResponding(true);
            });
    };

    const refreshSharedFiles = () => {
        fetch(`http://${instanceAddress}/sharedFiles`)
            .then(response => response.json())
            .then(files => {
                console.log("Shared files: ", files);
                setSharedFiles(files);
                setNotResponding(false);
            })
            .catch(e => {
                console.error(e);
                setNotResponding(true);
            });
    };

    return (
        <div id="mainContainer" className="container-fluid">
            <div className="header">
                <h1>
                    Peerster web-UI
                    <span id="nameGossiper" className="badge badge-secondary ml-2">{settings.Name}</span>
                    <button type="button" className="btn btn-danger ml-2" onClick={quitGossiperInstance}>
                        <FaWindowClose/></button>
                </h1>
            </div>
            {notResponding && notRespondingAlert()}
            <div id="main" className="container-fluid">
                <div className="row">
                    <div className="col-md-4">
                        <Contacts instanceAddress={instanceAddress} currentContact={currentContact}
                                  generalUnread={generalUnread} contacts={contacts} contactSelected={c => {
                            console.log("contactSelected", c);
                            console.log("Contacts state: ", contacts);
                            setCurrentContact(c)
                        }}/>
                        <Peers instanceAddress={instanceAddress} peers={peers}/>
                        <SharedFiles instanceAddress={instanceAddress} sharedFiles={sharedFiles} />
                    </div>

                    <div className="col-md-8" style={{height: "100%"}}>
                        <Messages
                            instanceAddress={instanceAddress}
                            currentContact={currentContact === generalName ? currentContact : contacts.get(currentContact).name}
                            messages={currentMessages}
                            name={settings.Name}/>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Gossip;