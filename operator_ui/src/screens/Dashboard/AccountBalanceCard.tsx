import React from 'react'

import { gql } from '@apollo/client'

import Card from '@material-ui/core/Card'
import CardActions from '@material-ui/core/CardActions'
import CardContent from '@material-ui/core/CardContent'
import CardHeader from '@material-ui/core/CardHeader'
import Grid from '@material-ui/core/Grid'
import CircularProgress from '@material-ui/core/CircularProgress'

import { fromJuels } from 'src/utils/tokens/link'
import {
  DetailsCardItemTitle,
  DetailsCardItemValue,
} from 'src/components/Cards/DetailsCard'
import Link from 'src/components/Link'

export const ACCOUNT_BALANCES_PAYLOAD__RESULTS_FIELDS = gql`
  fragment AccountBalancesPayload_ResultsFields on EthKey {
    address
    chain {
      id
    }
    ethBalance
    isFunding
    linkBalance
  }
`

export interface Props {
  data?: FetchAccountBalances
  loading: boolean
  errorMsg?: string
}

export const AccountBalanceCard: React.FC<Props> = ({
  data,
  errorMsg,
  loading,
}) => {
  const results = data?.ethKeys.results
  const ethKey = results && results.length > 0 ? results[0] : undefined

  return (
    <Card>
      <CardHeader title="Account Balance" />

      <CardContent>
        <Grid container spacing={16}>
          {loading && (
            <Grid
              item
              xs={12}
              style={{ display: 'flex', justifyContent: 'center' }}
            >
              <CircularProgress data-testid="loading-spinner" size={24} />
            </Grid>
          )}

          {errorMsg && (
            <Grid item xs={12}>
              <DetailsCardItemValue value={errorMsg} />
            </Grid>
          )}

          {ethKey && (
            <>
              <Grid item xs={12}>
                <DetailsCardItemTitle title="Address" />
                <DetailsCardItemValue value={ethKey.address} />
              </Grid>

              <Grid item xs={6}>
                <DetailsCardItemTitle title="ETH Balance" />
                <DetailsCardItemValue value={ethKey.ethBalance || '--'} />
              </Grid>

              <Grid item xs={6}>
                <DetailsCardItemTitle title="LINK Balance" />
                <DetailsCardItemValue
                  value={
                    ethKey.linkBalance ? fromJuels(ethKey.linkBalance) : '--'
                  }
                />
              </Grid>
            </>
          )}

          {!ethKey && !loading && !errorMsg && (
            <Grid item xs={12}>
              <DetailsCardItemValue value="No account available" />
            </Grid>
          )}
        </Grid>
      </CardContent>
      {results && results.length > 1 && (
        <CardActions style={{ marginLeft: 8 }}>
          <Link href="/keys" color="primary">
            View more accounts
          </Link>
        </CardActions>
      )}
    </Card>
  )
}
