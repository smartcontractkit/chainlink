import React, { useState } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import Grid from '@material-ui/core/Grid'
import Activity from 'components/Dashboards/Activity'
import TokenBalanceCard from 'components/Cards/TokenBalance'
import Footer from 'components/Footer'
import Content from 'components/Content'
import { fetchAccountBalance } from 'actionCreators'
import accountBalanceSelector from 'selectors/accountBalance'
import { v2 } from 'api'
import { JobRunV2, Resource } from 'core/store/models'

import { RecentJobs } from 'src/screens/Dashboard/RecentJobs'

type Props = {
  recentJobRunsCount: number
  recentlyCreatedPageSize: number
}

const fetchRuns = async (
  size: number,
): Promise<{ runs: Resource<JobRunV2>[]; count: number }> => {
  const response = await v2.runs.getAllJobRuns({ page: 1, size })

  return {
    runs: response.data as unknown as Resource<JobRunV2>[],
    count: response.meta.count,
  }
}

export const Index = ({
  recentJobRunsCount = 2,
  recentlyCreatedPageSize,
}: Props) => {
  const [runs, setRuns] = useState<Resource<JobRunV2>[]>([])
  const [count, setCount] = useState(0)
  const dispatch = useDispatch()
  const accountBalance = useSelector(accountBalanceSelector)

  React.useEffect(() => {
    document.title = 'Dashboard'
  }, [])

  React.useEffect(() => {
    dispatch(fetchAccountBalance())
    fetchRuns(recentJobRunsCount).then((data) => {
      setRuns(data.runs)
      setCount(data.count)
    })
  }, [dispatch, recentlyCreatedPageSize, recentJobRunsCount])

  return (
    <Content>
      <Grid container>
        <Grid item xs={8}>
          <Activity runs={runs} pageSize={recentJobRunsCount} count={count} />
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
