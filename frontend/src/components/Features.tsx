import React from 'react';
import "../css/Features.css"
const Features: React.FC = () => {
  return (
    <section id="features" className="features">
      <h2>Why Choose Our Wallet?</h2>
      <div className="features-grid">
        <div className="feature-card">
          <h3>Keyless Security</h3>
          <p>No more private keys to manage or worry about</p>
        </div>
        <div className="feature-card">
          <h3>ZK-Powered Privacy</h3>
          <p>Enhanced privacy with zero-knowledge proof technology</p>
        </div>
        <div className="feature-card">
          <h3>User-Friendly</h3>
          <p>Simple and intuitive interface for all users</p>
        </div>
      </div>
    </section>
  );
};

export default Features;