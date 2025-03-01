import React from 'react';
import { TokenIcon } from '@web3icons/react';
import '../css/Tokenlist.css';
interface TokensListProps {
  selectedNetwork: {
    id: string;
  };
}

const TokensList: React.FC<TokensListProps> = ({ selectedNetwork }) => {
  return (
    <div className="tokens-list">
      <div className="token-item">
        <TokenIcon 
          address="0xc944e90c64b2c07662a292be6244bdf05cda44a7" 
          network={selectedNetwork.id} 
          size="2.5rem" 
          className='token-logo' 
        />
        <div className="token-info">
          <span className="token-name">Ethereum</span>
          <span className="token-balance">2.5 ETH</span>
        </div>
        <span className="token-value">$4,250.00</span>
      </div>
      {/* Add more token items as needed */}
    </div>
  );
};

export default TokensList;