import { configureStore } from '@reduxjs/toolkit'
import { euiccSlice } from './reducers/euicc'

export const store = configureStore({
  reducer: {
    euicc: euiccSlice.reducer,
  },
})

export type RootState = ReturnType<typeof store.getState>

export type AppDispatch = typeof store.dispatch
