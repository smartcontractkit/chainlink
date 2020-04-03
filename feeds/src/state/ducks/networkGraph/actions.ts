/**
 * networkGraph/SET_TOOLTIP
 */
export interface SetTooltipAction {
  type: 'networkGraph/SET_TOOLTIP'
  payload: any
}

export function setTooltip(payload: any) {
  return {
    type: 'networkGraph/SET_TOOLTIP',
    payload,
  }
}

/**
 * networkGraph/SET_DRAWER
 */
export interface SetDrawerAction {
  type: 'networkGraph/SET_DRAWER'
  payload: any
}

export function setDrawer(payload: any) {
  return {
    type: 'networkGraph/SET_DRAWER',
    payload,
  }
}
