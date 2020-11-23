import React from 'react'
import { v2 } from 'api'
import { RouteComponentProps } from 'react-router-dom'
import Grid from '@material-ui/core/Grid'
import Card from '@material-ui/core/Card'
import Content from 'components/Content'
import StatusCard from 'components/StatusCard/StatusCard'
import TaskExpansionPanel from './TaskExpansionPanel/TaskExpansionPanel'
import { useErrorHandler } from 'hooks/useErrorHandler'
import { useLoadingPlaceholder } from 'hooks/useLoadingPlaceholder'
import RegionalNav from './RegionalNav'
import { DirectRequestJobRun } from '../sharedTypes'

type Props = RouteComponentProps<{
  jobSpecId: string
  jobRunId: string
}>

export const Show = (props: Props) => {
  const [jobRun, setState] = React.useState<DirectRequestJobRun>()

  const { jobSpecId, jobRunId } = props.match.params

  const { error, ErrorComponent, setError } = useErrorHandler()
  const { LoadingPlaceholder } = useLoadingPlaceholder(!error && !jobRun)

  React.useEffect(() => {
    document.title = 'Show Job Run'
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
    <div>
      <RegionalNav {...props.match.params} jobRun={jobRun} />
      <Content>
        <ErrorComponent />
        <LoadingPlaceholder />
        {jobRun && (
          <Grid container spacing={40}>
            <Grid item xs={8}>
              <Card>
                <TaskExpansionPanel jobRun={jobRun} />
              </Card>
            </Grid>
            <Grid item xs={4}>
              <StatusCard {...jobRun} title={jobRun.status} />
            </Grid>
          </Grid>
        )}
      </Content>
    </div>
  )
}

export default Show
