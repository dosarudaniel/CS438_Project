import React, {useState} from 'react';
import ShareFileModal from "./ShareFileModal";

const SharedFiles = ({instanceAddress, sharedFiles}) => {
    const [showModal, setShowModal] = useState(false);

    const hideModal = () => {
        setShowModal(false);
    };

    const sharedFilesList = () => {
        return sharedFiles.map(file =>
            <li className="list-group-item d-flex align-items-center justify-content-between" key={file.Name+file.Hash}>
                {file.Name}
                <button className="btn btn-secondary" onClick={()=>alert(file.Hash)}>Hash</button>
            </li>
        )
    };

    return (
        <div id="files">
            <ShareFileModal instanceAddress={instanceAddress} show={showModal} dismiss={hideModal} />
            <div className="niceBox">
                <h3>Shared Files</h3>
                <ul className="list-group" id="sharedFiles">
                    { sharedFilesList() }
                </ul>
                <button className="btn btn-primary" onClick={() => setShowModal(true)}>Share a new file</button>
            </div>
        </div>
    );
};

export default SharedFiles;