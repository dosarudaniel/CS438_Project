import React, {useState} from 'react';
import {Modal} from "react-bootstrap";
import Button from "react-bootstrap/Button";

const DownloadFileModal = ({instanceAddress, show, dismiss, currentContact}) => {
    const [targetFilename, setTargetFilename] = useState("");
    const [targetHash, setTargetHash] = useState("");

    const postNewFile = e => {
        e.preventDefault();
        fetch(`http://${instanceAddress}/downloadFile/${currentContact}/${targetFilename}/${targetHash}`, {
            method: 'POST'
        }).then(onCancel).catch(e => {
            console.error(e);
            alert("Error: see logs for details");
        });
    };

    const onCancel = () => {
        setTargetFilename("");
        setTargetHash("");
        dismiss();
    };

    const updateFilename = e => {
        setTargetFilename(e.target.value);
    };
    const updateHash = e => {
        setTargetHash(e.target.value);
    };

    const downloadableFilesSelection = () => {
        return (
            <form>
                <div className="form-group">
                    <label htmlFor="filename">File name</label>
                    <input id="filename" type="text" className="form-control"
                           value={targetFilename} onChange={updateFilename}/>
                </div>
                <div className="form-group">
                    <label htmlFor="hash">Hash</label>
                    <input id="hash" type="text" className="form-control"
                           value={targetHash} onChange={updateHash}/>
                </div>
            </form>
        )
    };

    return (
        <Modal show={show} onHide={dismiss}>
            <Modal.Header closeButton>
                <Modal.Title>Share a new file</Modal.Title>
            </Modal.Header>
            <Modal.Body>
                <div className="list-group">
                    {downloadableFilesSelection()}
                </div>
            </Modal.Body>
            <Modal.Footer>
                <Button variant="secondary" onClick={onCancel}>Cancel</Button>
                <Button variant="primary" onClick={postNewFile}>Download file</Button>
            </Modal.Footer>
        </Modal>
    );
};

export default DownloadFileModal;
