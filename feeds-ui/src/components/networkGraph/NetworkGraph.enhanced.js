import NetworkGraph from './NetworkGraph.component'
import { compose } from 'recompose'
import { connect } from 'react-redux'

import {
  aggregationSelectors,
  aggregationOperations,
} from 'state/ducks/aggregation'

const mapStateToProps = state => ({
  networkGraphNodes: aggregationSelectors.networkGraphNodes(state),
  networkGraphState: aggregationSelectors.networkGraphState(state),
  pendingAnswerId: state.aggregation.pendingAnswerId,
  updateHeight: state.aggregation.updateHeight,
})

const mapDispatchToProps = {
  fetchJobId: aggregationOperations.fetchJobId,
}

export default compose(connect(mapStateToProps, mapDispatchToProps))(
  NetworkGraph,
)
