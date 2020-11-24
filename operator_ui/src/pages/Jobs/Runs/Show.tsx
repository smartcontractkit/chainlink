import React from 'react'
import { v2 } from 'api'
import { Route, RouteComponentProps, Switch } from 'react-router-dom'
import { CardTitle } from '@chainlink/styleguide'
import { Card, Grid, Typography } from '@material-ui/core'
import Content from 'components/Content'
import StatusCard from './StatusCard'
import { useErrorHandler } from 'hooks/useErrorHandler'
import { useLoadingPlaceholder } from 'hooks/useLoadingPlaceholder'
import RegionalNav from './RegionalNav'
import { DirectRequestJobRun } from '../sharedTypes'
import { Overview } from './Overview/Overview'
import { Json } from './Json'

type Props = RouteComponentProps<{
  jobSpecId: string
  jobRunId: string
}>

export const Show = ({ match }: Props) => {
  const [jobRun, setState] = React.useState<DirectRequestJobRun>()

  const { jobSpecId, jobRunId } = match.params

  const { error, ErrorComponent, setError } = useErrorHandler()
  const { LoadingPlaceholder } = useLoadingPlaceholder(!error && !jobRun)

  React.useEffect(() => {
    document.title = 'Job run details'
  }, [])

  const getJobRun = React.useCallback(async () => {
    return v2.runs
      .getJobSpecRun(jobRunId)
      .then((jobSpecRunsResponse) => {
        const jobRun = jobSpecRunsResponse.data
        setState({
          ...jobSpecRunsResponse.data.attributes,
          id: jobRun.id,
          jobId: jobSpecId,
        })
      })
      .catch(setError)
  }, [jobRunId, jobSpecId, setError])

  React.useEffect(() => {
    getJobRun()
  }, [getJobRun])

  const erroredTasks = jobRun
    ? jobRun.taskRuns.filter((tr) => tr.status === 'errored')
    : []

  return (
    <>
      <RegionalNav {...match.params} jobRun={jobRun} />
      <Content>
        <ErrorComponent />
        <LoadingPlaceholder />
        {jobRun && (
          <Grid container spacing={40}>
            <Grid item xs={8}>
              <Grid container spacing={40}>
                {erroredTasks.length > 0 && (
                  <Grid item xs={12}>
                    <Card>
                      <CardTitle divider>Errors</CardTitle>
                      <ul>
                        {erroredTasks.map((tr) => (
                          <li key={tr.id}>
                            <Typography variant="body1">
                              {tr.result.error}
                            </Typography>
                          </li>
                        ))}
                      </ul>
                    </Card>
                  </Grid>
                )}
                <Grid item xs={12}>
                  <Switch>
                    <Route
                      path={`${match.path}/json`}
                      render={() => <Json jobRun={jobRun} />}
                    />
                    <Route render={() => <Overview jobRun={jobRun} />} />
                  </Switch>
                </Grid>
              </Grid>
            </Grid>
            <Grid item xs={4}>
              <StatusCard {...jobRun} title={jobRun.status} />
            </Grid>
          </Grid>
        )}
      </Content>
    </>
  )
}

export default Show
