import { createSlice,PayloadAction } from "@reduxjs/toolkit";

interface UserInfoState {
  walletAddress:string,
  email:string,
}
const initialState:UserInfoState = {
    walletAddress:'',
    email:'',
}

  const userInfoSlice = createSlice({
    name: 'userInfo',
    initialState,
    reducers: {
      setWalletAddress(state, action:PayloadAction<string>) {
        state.walletAddress = action.payload;
      },
      setEmail(state, action:PayloadAction<string>) {
        state.email = action.payload;
      },
      resetAcInfo(state) {
        state.walletAddress = '';
        state.email = '';
      },
    },
  });
  
  export const { setWalletAddress, setEmail, resetAcInfo } = userInfoSlice.actions;
  export default userInfoSlice.reducer;