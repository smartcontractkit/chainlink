import React from 'react'
import { v2 } from 'api'
import { Route, RouteComponentProps, Switch } from 'react-router-dom'
import Grid from '@material-ui/core/Grid'
import Content from 'components/Content'
import StatusCard from 'components/StatusCard/StatusCard'
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
    document.title = 'Show job run'
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

  return (
    <>
      <RegionalNav {...match.params} jobRun={jobRun} />
      <Content>
        <ErrorComponent />
        <LoadingPlaceholder />
        {jobRun && (
          <Grid container spacing={40}>
            <Grid item xs={8}>
              <Switch>
                <Route
                  path={`${match.path}/json`}
                  render={() => <Json jobRun={jobRun} />}
                />
                <Route render={() => <Overview jobRun={jobRun} />} />
              </Switch>
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
