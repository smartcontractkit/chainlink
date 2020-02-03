import NetworkGraphInfo from './NetworkGraphInfo.component'
import { compose } from 'recompose'
import { connect } from 'react-redux'

const mapStateToProps = state => ({
  currentAnswer: state.aggregation.currentAnswer,
  pendingAnswerId: state.aggregation.pendingAnswerId,
  oracleResponse: state.aggregation.oracleResponse,
  oracles: state.aggregation.oracles,
  requestTime: state.aggregation.requestTime,
  minimumResponses: state.aggregation.minimumResponses,
  updateHeight: state.aggregation.updateHeight,
})

export default compose(connect(mapStateToProps))(NetworkGraphInfo)
