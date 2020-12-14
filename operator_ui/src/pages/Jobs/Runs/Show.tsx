import React from 'react'
import { v2 } from 'api'
import { Route, RouteComponentProps, Switch } from 'react-router-dom'
import { CardTitle } from '@chainlink/styleguide'
import { Card, Grid, Typography } from '@material-ui/core'
import Content from 'components/Content'
import { useErrorHandler } from 'hooks/useErrorHandler'
import StatusCard from './StatusCard'
import { useLoadingPlaceholder } from 'hooks/useLoadingPlaceholder'
import RegionalNav from './RegionalNav'
import { DirectRequestJobRun, PipelineJobRun } from '../sharedTypes'
import { Overview } from './Overview/Overview'
import { PipelineJobRunOverview } from './Overview/PipelineJobRunOverview'
import { Json } from './Json'
import { isOcrJob } from '../utils'
import { TaskList } from '../TaskListDag'
import { augmentOcrTasksList } from './augmentOcrTasksList'
import {
  transformDirectRequestJobRun,
  transformPipelineJobRun,
} from '../transformJobRuns'

function getErrorsList(
  jobRun: DirectRequestJobRun | PipelineJobRun | undefined,
): string[] {
  if (jobRun?.type === 'Direct request job run') {
    return jobRun.taskRuns
      .filter(({ status }) => status === 'errored')
      .map((tr) => tr.result.error)
      .filter((error): error is string => error !== null)
  }

  if (jobRun?.type === 'Off-chain reporting job run' && jobRun.errors) {
    return jobRun.errors.filter((error): error is string => error !== null)
  }

  return []
}

type Props = RouteComponentProps<{
  jobSpecId: string
  jobRunId: string
}>

export const Show = ({ match }: Props) => {
  const [jobRun, setState] = React.useState<
    DirectRequestJobRun | PipelineJobRun
  >()

  const { jobSpecId, jobRunId } = match.params
  const { error, ErrorComponent, setError } = useErrorHandler()
  const { LoadingPlaceholder } = useLoadingPlaceholder(!error && !jobRun)

  React.useEffect(() => {
    document.title = 'Job run details'
  }, [])

  const getJobRun = React.useCallback(async () => {
    if (isOcrJob(jobSpecId)) {
      v2.ocrRuns
        .getJobSpecRun({ jobSpecId, runId: jobRunId })
        .then((jobSpecRunResponse) => jobSpecRunResponse.data)
        .then(transformPipelineJobRun(jobSpecId))
        .then(setState)
        .catch(setError)
    } else {
      return v2.runs
        .getJobSpecRun(jobRunId)
        .then((jobSpecRunResponse) => jobSpecRunResponse.data)
        .then(transformDirectRequestJobRun(jobSpecId))
        .then(setState)
        .catch(setError)
    }
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
                    <Route
                      path={`${match.path}/json`}
                      render={() => <Json jobRun={jobRun} />}
                    />
                    {jobRun.type === 'Direct request job run' && (
                      <Route render={() => <Overview jobRun={jobRun} />} />
                    )}
                    {jobRun.type === 'Off-chain reporting job run' && (
                      <Route
                        render={() => (
                          <PipelineJobRunOverview jobRun={jobRun} />
                        )}
                      />
                    )}
                  </Switch>
                </Grid>
              </Grid>
            </Grid>
            <Grid item xs={4}>
              <Grid container spacing={40}>
                <Grid item xs={12}>
                  <StatusCard {...jobRun} title={jobRun.status} />
                </Grid>
                {jobRun.type === 'Off-chain reporting job run' && (
                  <Grid item xs={12}>
                    <Card style={{ overflow: 'visible' }}>
                      <CardTitle divider>Task list</CardTitle>
                      <TaskList stratify={augmentOcrTasksList({ jobRun })} />
                    </Card>
                  </Grid>
                )}
              </Grid>
            </Grid>
          </Grid>
        )}
      </Content>
    </>
  )
}

export default Show
