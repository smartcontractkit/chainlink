export const SET_TOOLTIP = 'networkGraph/SET_TOOLTIP'
export const SET_DRAWER = 'networkGraph/SET_DRAWER'

export interface SetTooltipAction {
  type: typeof SET_TOOLTIP
  payload: any
}

export interface SetDrawerAction {
  type: typeof SET_DRAWER
  payload: any
}

export type NetworkGraphActionTypes = SetTooltipAction | SetDrawerAction
