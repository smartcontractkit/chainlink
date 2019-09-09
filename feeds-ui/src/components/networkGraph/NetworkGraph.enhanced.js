import NetworkGraph from './NetworkGraph.component'
import { compose } from 'recompose'
import { connect } from 'react-redux'

import { aggregationSelectors } from 'state/ducks/aggregation'

const mapStateToProps = state => ({
  networkGraphNodes: aggregationSelectors.networkGraphNodes(state),
  networkGraphLinks: aggregationSelectors.networkGraphLinks(state),
  networkGraphData: aggregationSelectors.networkGraphData(state)
})

const mapDispatchToProps = {}

export default compose(
  connect(
    mapStateToProps,
    mapDispatchToProps
  )
)(NetworkGraph)
