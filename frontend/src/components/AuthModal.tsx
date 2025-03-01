import React, { useState } from 'react';
import { FcGoogle } from 'react-icons/fc';
import { FaFacebook } from 'react-icons/fa';
import { NetworkIcon } from '@web3icons/react';
import { authService } from '../services/authService';
import '../css/AuthModal.css';

interface AuthModalProps {
  isOpen: boolean;
  onSuccess: () => void;
  onClose: () => void;
}

const AuthModal: React.FC<AuthModalProps> = ({ isOpen, onSuccess, onClose }) => {
  const [error, setError] = useState<string>('');

  if (!isOpen) return null;

  const handleGoogleLogin = async () => {
    try {
      authService.loginWithGoogle();
      console.log("login started")
    } catch (error) {
      setError('Google authentication failed');
      console.error('Authentication failed:', error);
    }
  };

  const handleMetaMaskLogin = async () => {
    try {
      // const account = await authService.loginWithMetaMask();
      // if (account) {
        onSuccess();
      // }
    } catch (error) {
      setError('MetaMask authentication failed');
      console.error('Authentication failed:', error);
    }
  };

  const handleFacebookLogin = async () => {
    try {
      authService.loginWithFacebook();
    } catch (error) {
      setError('Facebook authentication failed');
      console.error('Authentication failed:', error);
    }
  };

  const handleFormSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    // const formData = new FormData(e.currentTarget as HTMLFormElement);
    // const email = formData.get('email') as string;
    // const password = formData.get('password') as string;

    try {
      // Implement email/password login logic here
      onSuccess();
    } catch (error) {
      setError('Email/password authentication failed');
      console.error('Authentication failed:', error);
    }
  };

  return (
    <div className="modal-overlay">
      <div className="modal-content">
        <button className="close-button" onClick={onClose}>&times;</button>
        
        <h2>Welcome</h2>
        <p>Choose how you'd like to continue</p>
        
        {error && <div className="error-message">{error}</div>}
        
        <div className="auth-buttons">
          <button className="auth-button google" onClick={handleGoogleLogin}>
            <FcGoogle size={24} />
            <span>Continue with Google</span>
          </button>

          <button className="auth-button metamask" onClick={handleMetaMaskLogin}>
            <NetworkIcon id='ethereum' size={32} variant="branded" />
            <span>Continue with MetaMask</span>
          </button>
          
          <button className="auth-button facebook" onClick={handleFacebookLogin}>
            <FaFacebook size={24} />
            <span>Continue with Facebook</span>
          </button>
        </div>

        <div className="divider">
          <span>or</span>
        </div>

        <form className="auth-form" onSubmit={handleFormSubmit}>
          <input 
            type="email" 
            name="email"
            placeholder="Email address" 
            className="auth-input"
            required
          />
          <input 
            type="password" 
            name="password"
            placeholder="Password" 
            className="auth-input"
            required
          />
          <button type="submit" className="submit-button">
            Continue
          </button>
        </form>
      </div>
    </div>
  );
};

export default AuthModal;