import React ,{useState} from 'react';
import { useNavigate} from 'react-router-dom';
import AuthModal from './AuthModal';
import "../css/Hero.css"
const Hero: React.FC = () => {
  const navigate = useNavigate();
  const [showModal, setShowModal] = useState(false);

  const handleGetStarted = () => {
    setShowModal(true);
  };

  const handleModalClose = () => {
    setShowModal(false);
  };

  const handleSuccessfulAuth = () => {
    setShowModal(false);
    navigate('/wallet');
  };
  return (
    <section id="hero" className="hero">
      <div className="hero-content">
        <h1>Experience Keyless Web3 Login</h1>
        <p>Powered by Zero-Knowledge Proofs for Maximum Privacy</p>
        <button className="cta-button" onClick={handleGetStarted}>Get Started</button>
      </div>
      <AuthModal isOpen={showModal} onSuccess={handleSuccessfulAuth} onClose={handleModalClose} />
    </section>
  );
};

export default Hero;


