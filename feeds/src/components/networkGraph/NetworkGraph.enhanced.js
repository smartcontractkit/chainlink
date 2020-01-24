import { connect } from 'react-redux'
import {
  aggregationSelectors,
  aggregationOperations,
} from 'state/ducks/aggregation'
import NetworkGraph from './NetworkGraph.component'

const mapStateToProps = state => ({
  updateHeight: state.aggregation.updateHeight,
  oraclesData: aggregationSelectors.oraclesData(state),
  currentAnswer: state.aggregation.currentAnswer,
})

const mapDispatchToProps = {
  fetchJobId: aggregationOperations.fetchJobId,
}

export default connect(mapStateToProps, mapDispatchToProps)(NetworkGraph)
