const BASE_URL = 'http://localhost:3000/auth';

export const authService = {
  loginWithGoogle: () => {
    window.location.href = `${BASE_URL}/google`;
  },

  loginWithFacebook: () => {
    window.location.href = `${BASE_URL}/facebook`;
  },

//   loginWithMetaMask: async () => {
//     if (!window.ethereum) {
//       throw new Error('MetaMask is not installed');
//     }
//     const accounts = await window.ethereum.request({ method: 'eth_requestAccounts' });
//     return accounts[0];
//   },

  checkAuthStatus: async () => {
    const response = await fetch(`${BASE_URL}/status`, {
      credentials: 'include'
    });
    return response.json();
  }
};