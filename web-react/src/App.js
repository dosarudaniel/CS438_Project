import React, { useState } from 'react';
import GossiperSelector from './GossiperSelector';
import Gossip from './Gossip';

function App() {

  const [loginState, setLoginState] = useState(true);
  const [instanceAddress, setInstanceAddress] = useState("");

  const setGossiperInstance = (addess) => {
    console.log(addess);
    setInstanceAddress(addess);
    setLoginState(false);
  };

  const quitGossiperInstance = () => {
    setInstanceAddress("");
    setLoginState(true);
  };

  return (
    <div className="App">
      <div id="welcomeScreen" className={"container" + (loginState ? "" : " hideAnimation")}>
        <GossiperSelector loginState={loginState} setGossiperInstance={setGossiperInstance} />
      </div>
      {!loginState && <Gossip instanceAddress={instanceAddress} quitGossiperInstance={quitGossiperInstance} />}
    </div>
  );
}

export default App;
