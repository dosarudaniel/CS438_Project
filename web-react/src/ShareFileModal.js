import React, {useEffect, useState} from 'react';
import {Modal} from "react-bootstrap";
import Button from "react-bootstrap/Button";

const ShareFileModal = ({instanceAddress, show, dismiss}) => {
    const [newFile, setNewFile] = useState("");
    const [sharableFiles, setSharableFiles] = useState([]);
    const [pendingGet, setPendingGet] = useState(true);

    const postNewFile = e => {
        e.preventDefault();
        fetch(`http://${instanceAddress}/sharedFile/${newFile}`, {
            method: 'POST'
        }).then(() => {
            setNewFile("");
            dismiss();
        }).catch(e => {
            console.error(e);
            alert("Error: see logs for details");
        });
    };

    const onCancel = () => {
        setNewFile("");
        dismiss();
    };

    const getSharableFiles = () => {
        setPendingGet(true);
        setSharableFiles([]);
        fetch(`http://${instanceAddress}/sharableFiles`)
            .then(response => response.json())
            .then(files => {
                setSharableFiles(files);
                setPendingGet(false);
            }).catch(e => {
            console.error(e);
            setPendingGet(false);
        });
    };
    useEffect(() => {
        if (show) getSharableFiles();
    }, [show]);

    const sharableFilesSelection = () => {
        const sharableFileClasses = "h3 list-group-item list-group-item-action text-center";
        if (sharableFiles.length === 0) {
            return <h5 className="text-danger">No sharable file</h5>
        }
        return sharableFiles.sort().map(file => (
            <button
                type="button"
                key={file}
                className={sharableFileClasses + ((file === newFile) ? ' list-group-item-primary' : '')}
                onClick={() => setNewFile(file)}>
                {file}
            </button>
        ))
    };

    return (
        <Modal show={show} onHide={dismiss}>
            <Modal.Header closeButton>
                <Modal.Title>Share a new file</Modal.Title>
            </Modal.Header>
            <Modal.Body>
                <div className="list-group">
                    {pendingGet && <h1>Please wait...</h1>}
                    {!pendingGet && sharableFilesSelection()}
                </div>
            </Modal.Body>
            <Modal.Footer>
                <Button variant="secondary" onClick={onCancel}>Cancel</Button>
                <Button variant="primary" onClick={postNewFile}>Share file</Button>
            </Modal.Footer>
        </Modal>
    );
};

export default ShareFileModal;
