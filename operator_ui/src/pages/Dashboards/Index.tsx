import React, { useState } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import Grid from '@material-ui/core/Grid'
import Activity from 'components/Dashboards/Activity'
import TokenBalanceCard from 'components/Cards/TokenBalance'
import Footer from 'components/Footer'
import Content from 'components/Content'
import { fetchAccountBalance } from 'actionCreators'
import accountBalanceSelector from 'selectors/accountBalance'
import RecentlyCreated from 'components/Jobs/RecentlyCreated'
import { v2 } from 'api'
import { JobRunV2, Resource } from 'core/store/models'

type Props = {
  recentJobRunsCount: number
  recentlyCreatedPageSize: number
}

const fetchJobs = async (pageSize: number) => {
  const jobs = await v2.jobs.getJobSpecs()

  return jobs.data
    .reverse()
    .slice(0, pageSize)
    .map(({ attributes, id }) => {
      let createdAt = ''

      if (attributes.directRequestSpec) {
        createdAt = attributes.directRequestSpec.createdAt
      } else if (attributes.fluxMonitorSpec) {
        createdAt = attributes.fluxMonitorSpec.createdAt
      } else if (attributes.offChainReportingOracleSpec) {
        createdAt = attributes.offChainReportingOracleSpec.createdAt
      } else if (attributes.keeperSpec) {
        createdAt = attributes.keeperSpec.createdAt
      } else if (attributes.cronSpec) {
        createdAt = attributes.cronSpec.createdAt
      } else if (attributes.webhookSpec) {
        createdAt = attributes.webhookSpec.createdAt
      } else if (attributes.vrfSpec) {
        createdAt = attributes.vrfSpec.createdAt
      }

      return { id, createdAt, name: attributes.name }
    })
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
  const [jobs, setJobs] = useState<
    { id: string; name: string | null; createdAt: string }[]
  >([])
  const [runs, setRuns] = useState<Resource<JobRunV2>[]>([])
  const [count, setCount] = useState(0)
  const dispatch = useDispatch()
  const accountBalance = useSelector(accountBalanceSelector)

  React.useEffect(() => {
    document.title = 'Dashboard'
  }, [])

  React.useEffect(() => {
    dispatch(fetchAccountBalance())
    fetchJobs(recentlyCreatedPageSize).then((fetchedJobs) =>
      setJobs(fetchedJobs),
    )
    fetchRuns(recentJobRunsCount).then((data) => {
      setRuns(data.runs)
      setCount(data.count)
    })
  }, [dispatch, recentlyCreatedPageSize, recentJobRunsCount])

  return (
    <Content>
      <Grid container>
        <Grid item xs={9}>
          <Activity runs={runs} pageSize={recentJobRunsCount} count={count} />
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
            <Grid item xs={12}>
              <RecentlyCreated jobs={jobs} />
            </Grid>
          </Grid>
        </Grid>
      </Grid>
      <Footer />
    </Content>
  )
}

export default Index
