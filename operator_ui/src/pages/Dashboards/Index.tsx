import React from 'react'
import { useDispatch, useSelector } from 'react-redux'
import Grid from '@material-ui/core/Grid'
import Activity from 'components/Dashboards/Activity'
import TokenBalanceCard from 'components/Cards/TokenBalance'
import Footer from 'components/Footer'
import Content from 'components/Content'
import { fetchAccountBalance } from 'actionCreators'
import accountBalanceSelector from 'selectors/accountBalance'

type Props = {
  recentJobRunsCount: number
  recentlyCreatedPageSize: number
}

// NOTE - The Recent job runs shown in activity and recent jobs do not have a
// JPV2 equivalent. These have been removed for now.
export const Index = ({ recentJobRunsCount = 2 }: Props) => {
  const dispatch = useDispatch()
  const accountBalance = useSelector(accountBalanceSelector)

  React.useEffect(() => {
    document.title = 'Dashboard'
  }, [])

  React.useEffect(() => {
    dispatch(fetchAccountBalance())
  }, [dispatch])

  return (
    <Content>
      <Grid container>
        <Grid item xs={9}>
          {/* We don't have a JPV2 equivalent for this yet */}
          <Activity runs={[]} pageSize={recentJobRunsCount} count={0} />
        </Grid>
        <Grid item xs={3}>
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

            {/* We don't have a JPV2 equivalent for this yet */}
            {/* <Grid item xs={12}>
              <RecentlyCreatedJobs jobs={recentlyCreatedJobs} />
            </Grid> */}
          </Grid>
        </Grid>
      </Grid>
      <Footer />
    </Content>
  )
}

export default Index
