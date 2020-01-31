import * as actions from './actions'

const setTooltip = payload => {
  return async dispatch => {
    dispatch(actions.setTooltip(payload))
  }
}

const setDrawer = payload => {
  return async dispatch => {
    dispatch(actions.setDrawer(payload))
  }
}

export { setTooltip, setDrawer }
