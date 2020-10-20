// react-static does not support react-router-redux, connected-react-router or
// provide an escape hatch.
// https://github.com/nozzle/react-static/issues/211#issuecomment-389695521
//
// We need to handle this manually. mapDispatchToProps provides the props of
// the matched URL. We can dispatch this url and parse it in a reducer.
import { bindActionCreators } from 'redux'
import { RouterActionType } from '../reducers/actions'

const matchRouteAndMapDispatchToProps = (actionCreators) => (
  dispatch,
  ownProps,
) => {
  dispatch({
    type: RouterActionType.MATCH_ROUTE,
    match: ownProps.match,
  })

  return bindActionCreators(actionCreators, dispatch)
}

export default matchRouteAndMapDispatchToProps
