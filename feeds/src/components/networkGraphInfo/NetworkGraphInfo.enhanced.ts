import { connect } from 'react-redux'
import NetworkGraphInfo from './NetworkGraphInfo.component'

const mapStateToProps = (state: any) => ({
  currentAnswer: state.aggregation.currentAnswer,
  pendingAnswerId: state.aggregation.pendingAnswerId,
  oracleResponse: state.aggregation.oracleResponse,
  oracles: state.aggregation.oracles,
  requestTime: state.aggregation.requestTime,
  minimumResponses: state.aggregation.minimumResponses,
  updateHeight: state.aggregation.updateHeight,
})

export default connect(mapStateToProps)(NetworkGraphInfo)
