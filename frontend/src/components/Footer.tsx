import React from 'react';
import "../css/Footer.css"
const Footer: React.FC = () => {
  return (
    <footer id="footer" className="footer">
      <div className="footer-content">
        <div className="footer-links">
          <a href="/privacy">Privacy Policy</a>
          <a href="/terms">Terms of Service</a>
        </div>
        <div className="social-links">
          <a href="https://twitter.com">Twitter</a>
          <a href="https://discord.com">Discord</a>
          <a href="https://github.com">GitHub</a>
        </div>
      </div>
    </footer>
  );
};

export default Footer;