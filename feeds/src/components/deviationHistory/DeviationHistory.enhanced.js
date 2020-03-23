import { compose } from 'recompose'
import { connect } from 'react-redux'
import DeviationHistory from './DeviationHistory.component'

const mapStateToProps = state => ({
  answerHistory: state.aggregator.answerHistory,
})

export default compose(connect(mapStateToProps))(DeviationHistory)
