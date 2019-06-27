import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'react-redux'
import { useHooks, useEffect } from 'use-react-hooks'
import KeyValueList from '@chainlink/styleguide/components/KeyValueList'
import Content from 'components/Content'
import { fetchTransaction } from 'actions'
import transactionSelector from 'selectors/transaction'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'

export const Show = useHooks(props => {
  useEffect(() => {
    props.fetchTransaction(props.transactionId)
  }, [])
  const { transaction } = props

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
})

Show.propTypes = {
  classes: PropTypes.object.isRequired,
  transaction: PropTypes.object
}

const mapStateToProps = (state, ownProps) => {
  const transactionId = ownProps.match.params.transactionId
  const transaction = transactionSelector(state, transactionId)

  return { transactionId, transaction }
}

export const ConnectedShow = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({ fetchTransaction })
)(Show)

export default ConnectedShow
