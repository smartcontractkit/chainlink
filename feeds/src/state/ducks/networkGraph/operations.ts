import * as actions from './actions'

function setTooltip(payload: any) {
  return async (dispatch: any) => {
    dispatch(actions.setTooltip(payload))
  }
}

function setDrawer(payload: any) {
  return async (dispatch: any) => {
    dispatch(actions.setDrawer(payload))
  }
}

export { setTooltip, setDrawer }
