import React from 'react'

import Grid from '@material-ui/core/Grid'

import {
  DetailsCard,
  DetailsCardItemTitle,
  DetailsCardItemValue,
  DetailsCardItemBlockValue,
} from 'src/components/Cards/DetailsCard'
import titleize from 'src/utils/titleize'
// import { TimeAgo } from 'src/components/TimeAgo'

interface Props {
  tx: EthTransactionPayloadFields
}

export const TransactionCard: React.FC<Props> = ({ tx }) => {
  return (
    <DetailsCard>
      <Grid container>
        <Grid item xs={12}>
          <DetailsCardItemTitle title="Txn Hash" />
          <DetailsCardItemValue value={tx.hash} />
        </Grid>

        <Grid item xs={12} md={3}>
          <DetailsCardItemTitle title="Chain ID" />
          <DetailsCardItemValue value={tx.chain.id} />
        </Grid>

        <Grid item xs={12} md={3}>
          <DetailsCardItemTitle title="Status" />
          <DetailsCardItemValue value={titleize(tx.state)} />
        </Grid>

        <Grid item xs={12} md={3}>
          <DetailsCardItemTitle title="Nonce" />
          <DetailsCardItemValue value={tx.nonce} />
        </Grid>

        <Grid item xs={12} md={3}>
          <DetailsCardItemTitle title="Sent At (Block)" />
          <DetailsCardItemValue value={tx.sentAt} />
        </Grid>

        <Grid item xs={12} md={6}>
          <DetailsCardItemTitle title="From" />
          <DetailsCardItemValue value={tx.from} />
        </Grid>

        <Grid item xs={12} md={6}>
          <DetailsCardItemTitle title="To" />
          <DetailsCardItemValue value={tx.to} />
        </Grid>

        <Grid item xs={12} md={3}>
          <DetailsCardItemTitle title="Value" />
          <DetailsCardItemValue value={tx.value} />
        </Grid>

        <Grid item xs={false} md={3}></Grid>

        <Grid item xs={12} md={3}>
          <DetailsCardItemTitle title="Gas Price" />
          <DetailsCardItemValue value={tx.gasPrice} />
        </Grid>

        <Grid item xs={12} md={3}>
          <DetailsCardItemTitle title="Gas Limit" />
          <DetailsCardItemValue value={tx.gasLimit} />
        </Grid>

        <Grid item xs={12}>
          <DetailsCardItemTitle title="Raw Hex" />
          <DetailsCardItemBlockValue>{tx.hex}</DetailsCardItemBlockValue>
        </Grid>

        <Grid item xs={12}>
          <DetailsCardItemTitle title="Data" />
          <DetailsCardItemBlockValue>{tx.data}</DetailsCardItemBlockValue>
        </Grid>
      </Grid>
    </DetailsCard>
  )
}
