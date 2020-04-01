import { connect } from 'react-redux'
import DeviationHistory from './DeviationHistory.component'
import { AppState } from 'state'

const mapStateToProps = (state: AppState) => ({
  answerHistory: state.aggregator.answerHistory,
})

export default connect(mapStateToProps)(DeviationHistory)
