import React, { useState,useEffect } from 'react';
import '../css/AuthModal.css';
import axios from 'axios';
import { useAppDispatch,useAppSelector } from '../redux/hooks';
import {setWalletAddress,setEmail,resetAcInfo} from '../redux/slices/index'

interface AuthModalProps {
  isOpen: boolean;
  onSuccess: () => void;
  onClose: () => void;
}
interface AuthResponse{
  credential:string
}
interface UserInfoState {
  walletAddress:string,
  email:string,
}

const AuthModal: React.FC<AuthModalProps> = ({ isOpen, onSuccess, onClose }) => {

  const [error, setError] = useState<string>('');

  const obj:UserInfoState = useAppSelector((state) => state.userInfo);
  const dispatch = useAppDispatch();

  const [email, setemail] = useState<string>(obj.email);
  const [address,setaddress]=useState<string>(obj.walletAddress);

  if (!isOpen) return null;

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    try {
      // call the post method: POST/credentials
      const res = await axios.post<AuthResponse>('/credentials', { email });
      setaddress(res.data.credential);
      dispatch(setWalletAddress('0xeubfdn'))
      dispatch(setEmail("keo"));
      
      onSuccess();
    } catch (error) {
      setError('Email verification failed');
      console.error('Verification failed:', error);
    }
  };

  return (
    <div className="modal-overlay">
      <div className="modal-content">
        <button className="close-button" onClick={onClose}>&times;</button>
        
        <h2>Welcome</h2>
        <p>Enter your email to continue</p>
        
        {error && <div className="error-message">{error}</div>}
        
        <form className="auth-form" onSubmit={handleSubmit}>
          <input 
            type="email" 
            value={email}
            onChange={(e) => setemail(e.target.value)}
            placeholder="Email address" 
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