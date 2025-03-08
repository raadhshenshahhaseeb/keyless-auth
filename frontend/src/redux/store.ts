import { configureStore } from '@reduxjs/toolkit'
import userInfoSlice from './slices/index'

export const store = configureStore({
  reducer: {
    userInfo: userInfoSlice,
  },
})

export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch