import React from 'react'
import { v2 } from 'api'
import { Route, Switch, useParams, useRouteMatch } from 'react-router-dom'
import { CardTitle } from 'components/CardTitle'
import { Card, Grid, Typography } from '@material-ui/core'
import Content from 'components/Content'
import { useErrorHandler } from 'hooks/useErrorHandler'
import StatusCard from './StatusCard'
import { useLoadingPlaceholder } from 'hooks/useLoadingPlaceholder'
import RegionalNav from './RegionalNav'
import { PipelineJobRun } from '../sharedTypes'
import { PipelineJobRunOverview } from './Overview/PipelineJobRunOverview'
import { Json } from './Json'
import { TaskList } from '../TaskListDag'
import { augmentOcrTasksList } from './augmentOcrTasksList'
import { transformPipelineJobRun } from '../transformJobRuns'

interface RouteParams {
  jobId: string
  jobRunId: string
}

function getErrorsList(jobRun?: PipelineJobRun): string[] {
  if (jobRun?.errors) {
    return jobRun.errors.filter((error): error is string => error !== null)
  }

  return []
}

export const Show = () => {
  const { path } = useRouteMatch()
  const { jobId, jobRunId } = useParams<RouteParams>()
  const [jobRun, setState] = React.useState<PipelineJobRun>()

  const { error, ErrorComponent, setError } = useErrorHandler()
  const { LoadingPlaceholder } = useLoadingPlaceholder(!error && !jobRun)

  React.useEffect(() => {
    document.title = 'Job run details'
  }, [])

  const getJobRun = React.useCallback(async () => {
    v2.runs
      .getJobRun({ jobId, runId: jobRunId })
      .then((res) => res.data)
      .then(transformPipelineJobRun(jobId))
      .then(setState)
      .catch(setError)
  }, [jobRunId, jobId, setError])

  React.useEffect(() => {
    getJobRun()
  }, [getJobRun])

  return (
    <>
      <RegionalNav jobId={jobId} jobRunId={jobRunId} jobRun={jobRun} />

      <Content>
        <ErrorComponent />
        <LoadingPlaceholder />

        {jobRun && (
          <Grid container spacing={40}>
            <Grid item xs={4}>
              <Grid container spacing={40}>
                <Grid item xs={12}>
                  <StatusCard {...jobRun} title={jobRun.status} />
                </Grid>
                {jobRun.status == 'errored' && (
                  <Grid item xs={12}>
                    <Card style={{ overflow: 'visible' }}>
                      <CardTitle divider>Task list</CardTitle>
                      <TaskList stratify={augmentOcrTasksList({ jobRun })} />
                    </Card>
                  </Grid>
                )}
              </Grid>
            </Grid>
            <Grid item xs={8}>
              <Grid container spacing={40}>
                {getErrorsList(jobRun).length > 0 && (
                  <Grid item xs={12}>
                    <Card>
                      <CardTitle divider>Errors</CardTitle>
                      <ul>
                        {getErrorsList(jobRun).map((error, index) => (
                          <li key={error + index}>
                            <Typography variant="body1">{error}</Typography>
                          </li>
                        ))}
                      </ul>
                    </Card>
                  </Grid>
                )}
                <Grid item xs={12}>
                  <Switch>
                    <Route path={`${path}/json`}>
                      <Json jobRun={jobRun} />
                    </Route>
                    {jobRun.status == 'errored' && (
                      <Route path={path}>
                        <PipelineJobRunOverview jobRun={jobRun} />
                      </Route>
                    )}
                  </Switch>
                </Grid>
              </Grid>
            </Grid>
          </Grid>
        )}
      </Content>
    </>
  )
}

export default Show
