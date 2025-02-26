import React, { useState } from 'react';
import '../css/Header.css';

const Header: React.FC = () => {
  const [isMenuOpen, setIsMenuOpen] = useState(false);

  return (
    <header className="header">
      <div className="header-container">
        <div className="logo">
          {/* <img src="/logo.svg" alt="Web3 Wallet" /> */}
          <span>ZK Wallet</span>
        </div>

        <nav className={`nav-menu ${isMenuOpen ? 'active' : ''}`}>
          <ul>
            <li><a href="#features">Features</a></li>
            <li><a href="#security">Security</a></li>
            <li><a href="#how-it-works">How it Works</a></li>
            <li><a href="#faq">FAQ</a></li>
          </ul>
        </nav>

        <div className="header-actions">
          <button className="connect-wallet">Connect Wallet</button>
          <button 
            className="mobile-menu-btn"
            onClick={() => setIsMenuOpen(!isMenuOpen)}
          >
            <span></span>
            <span></span>
            <span></span>
          </button>
        </div>
      </div>
    </header>
  );
};

export default Header;