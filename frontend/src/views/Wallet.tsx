import React, { useState } from 'react';
import Header from '../components/Header';
import TokensList from '../components/TokensList';
import NftGrid from '../components/NftGrid';
import ActivityList from '../components/ActivityList';
import '../css/Wallet.css';
import { NetworkIcon } from '@web3icons/react';

const Wallet: React.FC = () => {
  const [activeTab, setActiveTab] = useState('tokens');
  const [showNetworkDropdown, setShowNetworkDropdown] = useState(false);
  const [selectedNetwork, setSelectedNetwork] = useState({
    name: 'Ethereum',
    logo: <NetworkIcon id='ethereum' size={32} variant="branded"/>,
    id: 'ethereum',
    chainId: '0x1'
  });

  const networks = [
    { 
      name: 'Ethereum', 
      logo: <NetworkIcon id='ethereum' size={32} variant="branded" />, 
      id: 'ethereum',
      chainId: '0x1'
    },
    { 
      name: 'Polygon', 
      logo: <NetworkIcon id='polygon' size={32} variant="branded" />, 
      id: 'polygon',
      chainId: '0x89' 
    },   
    { 
      name: 'Optimism', 
      logo: <NetworkIcon id='optimism' size={32} variant="branded" />, 
      id: 'optimism',
      chainId: '0xa' 
    },
    { 
      name: 'Base', 
      logo: <NetworkIcon id='base' size={32} variant="branded" />, 
      id: 'base',
      chainId: '0x2105'
    },
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
      {activeTab === 'tokens' && <TokensList selectedNetwork={selectedNetwork} />}
      {activeTab === 'nfts' && <NftGrid />}
      {activeTab === 'activity' && <ActivityList selectedNetwork={selectedNetwork} />}
      </div>
    </div>
    </>
  );
};

export default Wallet;