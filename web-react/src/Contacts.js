import React from 'react';

const Contacts = ({instanceAddress, currentContact, generalUnread, contacts, contactSelected}) => {
    const contactClasses = "list-group-item d-flex justify-content-between align-items-center";
    const generalName = 'General';

    const clickCallback = e => {
        console.log("Contact selected event:", e);
        console.log("Contact selected event target:", e.target);
        contactSelected(e.target.value);
    };

    const contactList = contacts => [...contacts.entries()].sort().map(contact => {
            const currentClasses = contactClasses + ((currentContact === contact[0]) ? ' active' : '');
            return (
                <button
                    type="button"
                    onClick={clickCallback}
                    key={contact[0]}
                    value={contact[0]}
                    className={currentClasses}>
                    <span className="d-flex align-items-center">
                        <span className="badge badge-secondary badge-pill mr-1">user</span>
                        {contact[1].name}
                    </span>
                    {contact[1].unread > 0 && <span className="badge badge-primary badge-pill">{contact[1].unread}</span>}
                </button>
            );
        }
    );

    return (
        <div id="contacts">
            <div className="niceBox">
                <h3>Contacts</h3>
                <ul className="list-group" id="contactList">
                    <button
                        type="button"
                        onClick={clickCallback}
                        value={generalName}
                        className={contactClasses + ((currentContact === generalName) ? ' active' : '')}>
                        General
                        {generalUnread > 0 && <span className="badge badge-primary badge-pill">{generalUnread}</span>}
                    </button>
                    {contactList(contacts)}
                </ul>
            </div>
        </div>
    );
};

export default Contacts;