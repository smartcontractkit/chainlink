import React from 'react'
import { useDispatch, useSelector } from 'react-redux'
import Grid from '@material-ui/core/Grid'

import TokenBalanceCard from 'components/Cards/TokenBalance'
import Footer from 'components/Footer'
import Content from 'components/Content'
import { fetchAccountBalance } from 'actionCreators'
import accountBalanceSelector from 'selectors/accountBalance'

import { Activity } from 'src/screens/Dashboard/Activity'
import { RecentJobs } from 'src/screens/Dashboard/RecentJobs'

type Props = {
  recentJobRunsCount: number
  recentlyCreatedPageSize: number
}

export const Index = ({
  recentJobRunsCount = 2,
  recentlyCreatedPageSize,
}: Props) => {
  const dispatch = useDispatch()
  const accountBalance = useSelector(accountBalanceSelector)

  React.useEffect(() => {
    dispatch(fetchAccountBalance())
  }, [dispatch, recentlyCreatedPageSize, recentJobRunsCount])

  return (
    <Content>
      <Grid container>
        <Grid item xs={8}>
          <Activity />
        </Grid>
        <Grid item xs={4}>
          <Grid container>
            <Grid item xs={12}>
              <TokenBalanceCard
                title="Link Balance"
                value={accountBalance?.linkBalance}
              />
            </Grid>
            <Grid item xs={12}>
              <TokenBalanceCard
                title="Ether Balance"
                value={accountBalance?.ethBalance}
              />
            </Grid>
            <Grid item xs={12}>
              <RecentJobs />
            </Grid>
          </Grid>
        </Grid>
      </Grid>
      <Footer />
    </Content>
  )
}

export default Index
