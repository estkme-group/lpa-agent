import { createAsyncThunk, createSlice, PayloadAction } from '@reduxjs/toolkit'
import { RootState } from '../store'
import { Information, Profile } from './euicc-types'

export interface EUICCState {
  information: Information
  profiles: Profile[]
}

const initialState: EUICCState = {
  information: { eid: '', default_smds: '' },
  profiles: [],
}

export const fetchInformation = createAsyncThunk('lpa/information', async (_, { dispatch }) => {
  const response = await fetch('/api/lpa/')
  if (!response.ok) return
  dispatch(setInformation(await response.json()))
})

export const fetchProfileList = createAsyncThunk('lpa/profile-list', async (_, { dispatch }) => {
  const response = await fetch('/api/lpa/profile/')
  if (!response.ok) return
  dispatch(setProfileList(await response.json()))
})

export const putProfile = createAsyncThunk('lpa/putProfile', async (profile: Profile, { dispatch }) => {
  const response = await fetch(`/api/lpa/profile/${profile.iccid}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(profile),
  })
  if (!response.ok) return
  await dispatch(fetchProfileList())
})

export const purge = createAsyncThunk('lpa/purge', async (_, { dispatch }) => {
  const response = await fetch('/api/lpa/', {
    method: 'DELETE',
  })
  if (!response.ok) return
  await dispatch(fetchProfileList())
})

export const deleteProfile = createAsyncThunk(
  'lpa/delete-profile',
  async (profile: Profile, { dispatch }) => {
    const response = await fetch(`/api/lpa/profile/${profile.iccid}`, {
      method: 'DELETE',
    })
    if (!response.ok) return
    await dispatch(fetchProfileList())
  },
)

interface RequestDownloadProfile {
  readonly smdp?: string
  readonly matchingId?: string
  readonly imei?: string
  readonly confirmCode?: string
}

export const downloadProfile = createAsyncThunk(
  'lpa/download-profile',
  async (cfg: RequestDownloadProfile, { dispatch }) => {
    const response = await fetch('/api/lpa/profile/download', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        smdp: cfg.smdp,
        matching_id: cfg.matchingId,
        imei: cfg.imei,
        confirm_code: cfg.confirmCode,
      }),
    })
    if (!response.ok) return null
    dispatch(fetchProfileList())
    return response.json()
  },
)

export const euiccSlice = createSlice({
  name: 'euicc',
  initialState,
  reducers: {
    setInformation(state, action: PayloadAction<Information>) {
      state.information = action.payload
    },
    setProfileList(state, action: PayloadAction<Profile[]>) {
      state.profiles = action.payload
    },
  },
})

export const { setInformation, setProfileList } = euiccSlice.actions

export const selectInformation = (state: RootState) => state.euicc.information

export const selectProfileList = (state: RootState) => state.euicc.profiles
