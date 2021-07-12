import { KeyValueList } from 'components/KeyValueList'
import { fetchTransaction } from 'actionCreators'
import Content from 'components/Content'
import PropTypes from 'prop-types'
import React, { useEffect } from 'react'
import { connect } from 'react-redux'
import transactionSelector from 'selectors/transaction'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'

export const Show = ({ fetchTransaction, transaction, transactionId }) => {
  useEffect(() => {
    fetchTransaction(transactionId)
  }, [fetchTransaction, transactionId])

  return (
    <Content>
      {transaction && (
        <KeyValueList
          title={transaction.id}
          entries={Object.entries(transaction)}
          titleize
        />
      )}
    </Content>
  )
}

Show.propTypes = {
  classes: PropTypes.object.isRequired,
  transaction: PropTypes.object,
}

const mapStateToProps = (state, ownProps) => {
  const transactionId = ownProps.match.params.transactionId
  const transaction = transactionSelector(state, transactionId)

  return { transactionId, transaction }
}

export const ConnectedShow = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({ fetchTransaction }),
)(Show)

export default ConnectedShow
