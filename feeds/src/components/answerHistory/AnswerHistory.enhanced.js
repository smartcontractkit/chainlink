import { compose } from 'recompose'
import { connect } from 'react-redux'
import AnswerHistory from './AnswerHistory.component'

const mapStateToProps = state => ({
  answerHistory: state.aggregation.answerHistory,
})

export default compose(connect(mapStateToProps))(AnswerHistory)
