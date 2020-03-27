import { connect } from 'react-redux'
import AnswerHistory from './AnswerHistory.component'
import { AppState } from 'state'

const mapStateToProps = (state: AppState) => ({
  answerHistory: state.aggregator.answerHistory,
})

export default connect(mapStateToProps)(AnswerHistory)
