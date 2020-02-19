import OracleTable from './OracleTable.component'
import { connect } from 'react-redux'

import {
  aggregationSelectors,
  aggregationOperations,
} from 'state/ducks/aggregation'

const mapStateToProps = state => ({
  networkGraphNodes: aggregationSelectors.networkGraphNodes(state),
  networkGraphState: aggregationSelectors.networkGraphState(state),
  ethGasPrice: state.aggregation.ethGasPrice,
})

const mapDispatchToProps = {
  fetchEthGasPrice: aggregationOperations.fetchEthGasPrice,
}

export default connect(mapStateToProps, mapDispatchToProps)(OracleTable)
