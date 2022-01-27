import React from 'react'

import { gql } from '@apollo/client'

import Grid from '@material-ui/core/Grid'

import Content from 'src/components/Content'
import { Heading1 } from 'src/components/Heading/Heading1'
import { TransactionCard } from './TransactionCard'

export const ETH_TRANSACTION_PAYLOAD_FIELDS = gql`
  fragment EthTransactionPayloadFields on EthTransaction {
    chain {
      id
    }
    data
    from
    gasLimit
    gasPrice
    hash
    hex
    nonce
    sentAt
    state
    to
    value
  }
`

interface Props {
  tx: EthTransactionPayloadFields
}

export const TransactionView: React.FC<Props> = ({ tx }) => {
  return (
    <Content>
      <Grid container spacing={16}>
        <Grid item xs={12}>
          <Heading1>Transaction Details</Heading1>
        </Grid>

        <Grid item xs={12}>
          <TransactionCard tx={tx} />
        </Grid>
      </Grid>
    </Content>
  )
}
