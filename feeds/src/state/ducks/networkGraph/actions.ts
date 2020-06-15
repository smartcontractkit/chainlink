import { SET_TOOLTIP, SET_DRAWER, NetworkGraphActionTypes } from './types'

export function setTooltip(payload: any): NetworkGraphActionTypes {
  return {
    type: SET_TOOLTIP,
    payload,
  }
}

export function setDrawer(payload: any): NetworkGraphActionTypes {
  return {
    type: SET_DRAWER,
    payload,
  }
}
