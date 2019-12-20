import * as types from './types'

export const setTooltip = payload => ({
  type: types.SET_TOOLTIP,
  payload,
})

export const setDrawer = payload => ({
  type: types.SET_DRAWER,
  payload,
})
