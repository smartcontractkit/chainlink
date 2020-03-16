import { connect } from 'react-redux'
import OracleTable from './OracleTable.component'
import {
  aggregationSelectors,
  aggregationOperations,
} from 'state/ducks/aggregation'
import { Ducks } from 'state/ducks'

const mapStateToProps = (state: Ducks) => ({
  networkGraphNodes: aggregationSelectors.networkGraphNodes(state),
  networkGraphState: aggregationSelectors.networkGraphState(state),
  ethGasPrice: state.aggregation.ethGasPrice,
})

const mapDispatchToProps = {
  fetchEthGasPrice: aggregationOperations.fetchEthGasPrice,
}

export default connect(mapStateToProps, mapDispatchToProps)(OracleTable)
