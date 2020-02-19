import DeviationHistory from './DeviationHistory.component'
import { compose } from 'recompose'
import { connect } from 'react-redux'

const mapStateToProps = state => ({
  answerHistory: state.aggregation.answerHistory,
})

export default compose(connect(mapStateToProps))(DeviationHistory)
