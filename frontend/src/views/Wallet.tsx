import React, { useState } from 'react';
import Header from '../components/Header';
import '../css/Wallet.css';
import {NetworkIcon,TokenIcon} from '@web3icons/react'

const Wallet: React.FC = () => {
  const [activeTab, setActiveTab] = useState('tokens');
  const [showNetworkDropdown, setShowNetworkDropdown] = useState(false);
  const [selectedNetwork, setSelectedNetwork] = useState({
    name: 'Ethereum',
    logo: <NetworkIcon id='etereum' size={32} variant="branded"/>,
    id: 'ethereum'
  });

  const networks = [
    { name: 'Ethereum', logo: <NetworkIcon id='ethereum' size={32} variant="branded" />, id: 'ethereum' },
    { name: 'Polygon', logo: <NetworkIcon id='polygon' size={32} variant="branded" />, id: 'polygon' },   
    { name: 'Optimism', logo: <NetworkIcon id='optimism' size={32} variant="branded" />, id: 'optimism' },
    { name: 'Base', logo: <NetworkIcon id='base' size={32} variant="branded" />, id: 'base' },
  ];

  const handleNetworkChange = (network: typeof networks[0]) => {
    setSelectedNetwork(network);
    setShowNetworkDropdown(false);
  };

  return (
    <>
    <Header/>
    <div className="wallet-container">
      <div className="wallet-header">
      <div className="network-info">
            <div 
              className="network-selector" 
              onClick={() => setShowNetworkDropdown(!showNetworkDropdown)}
            >
              <span className="network-logo" ><NetworkIcon id={selectedNetwork.id} size={32} variant="branded" /></span>
              <span className="network-name">{selectedNetwork.name}</span>
              <span className="dropdown-arrow">▼</span>
            </div>
            {showNetworkDropdown && (
              <div className="network-dropdown">
                {networks.map((network) => (
                  <div
                    key={network.name}
                    className="network-option"
                    onClick={() => handleNetworkChange(network)}
                  >
                   <span className="network-logo" ><NetworkIcon id={network.id} size={32} variant="branded" /></span>
                    <span>{network.name}</span>
                  </div>
                ))}
              </div>
            )}
          </div>
        <div className="address-container">
          <span className="address">0x1234...5678</span>
          <button className="copy-button">Copy</button>
        </div>
        <button className="send-button">Send</button>
      </div>

      <div className="wallet-tabs">
        <button 
          className={`tab-button ${activeTab === 'tokens' ? 'active' : ''}`}
          onClick={() => setActiveTab('tokens')}
        >
          Tokens
        </button>
        <button 
          className={`tab-button ${activeTab === 'nfts' ? 'active' : ''}`}
          onClick={() => setActiveTab('nfts')}
        >
          NFTs
        </button>
        <button 
          className={`tab-button ${activeTab === 'activity' ? 'active' : ''}`}
          onClick={() => setActiveTab('activity')}
        >
          Activity
        </button>
      </div>

      <div className="wallet-content">
        {activeTab === 'tokens' && (
          <div className="tokens-list">
            <div className="token-item">
              <TokenIcon address="0xc944e90c64b2c07662a292be6244bdf05cda44a7" network={selectedNetwork.id} size="2.5rem" className='token-logo' />
              <div className="token-info">
                <span className="token-name">Ethereum</span>
                <span className="token-balance">2.5 ETH</span>
              </div>
              <span className="token-value">$4,250.00</span>
            </div>
            {/* Add more token items as needed */}
          </div>
        )}

        {activeTab === 'nfts' && (
          <div className="nfts-grid">
            {/* Add NFT cards here */}
          </div>
        )}

        {activeTab === 'activity' && (
          <div className="activity-list">
            {/* Add activity items here */}
          </div>
        )}
      </div>
    </div>
    </>
  );
};

export default Wallet;