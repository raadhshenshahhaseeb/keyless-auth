import React from 'react';
import { useNavigate } from 'react-router-dom';
import "../css/Hero.css"
const Hero: React.FC = () => {
  const navigate = useNavigate();
  const handleGetStarted = () => {
    navigate('/wallet');
  };
  return (
    <section id="hero" className="hero">
      <div className="hero-content">
        <h1>Experience Keyless Web3 Login</h1>
        <p>Powered by Zero-Knowledge Proofs for Maximum Privacy</p>
        <button className="cta-button" onClick={handleGetStarted}>Get Started</button>
      </div>
    </section>
  );
};

export default Hero;