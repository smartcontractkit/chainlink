import AnswerHistory from './AnswerHistory.component'
import { compose } from 'recompose'
import { connect } from 'react-redux'

import {
  aggregationOperations,
  aggregationSelectors
} from 'state/ducks/aggregation'

const mapStateToProps = state => ({
  answerHistory: state.aggregation.answerHistory
})

const mapDispatchToProps = {}

export default compose(
  connect(
    mapStateToProps,
    mapDispatchToProps
  )
)(AnswerHistory)
