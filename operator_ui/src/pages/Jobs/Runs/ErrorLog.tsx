import Grid from '@material-ui/core/Grid'
import { Theme, withStyles, WithStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import { fetchJobRun } from 'actionCreators'
import Content from 'components/Content'
import StatusCard from 'components/StatusCard/StatusCard'
import { AppState } from 'src/reducers'
import { JobRun, TaskRun } from 'operator_ui'
import React, { useEffect } from 'react'
import { connect } from 'react-redux'
import jobRunSelector from 'selectors/jobRun'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import RegionalNav from './RegionalNav'

const filterErrorTaskRuns = (jobRun: JobRun) => {
  return jobRun.taskRuns.filter((tr: TaskRun) => {
    return tr.status === 'errored'
  })
}

const detailsStyles = ({ spacing }: Theme) => ({
  list: {
    marginTop: spacing.unit * 4,
  },
})

interface DetailsProps extends WithStyles<typeof detailsStyles> {
  jobRun?: JobRun
}

const Details = withStyles(detailsStyles)(
  ({ jobRun, classes }: DetailsProps) => {
    if (!jobRun) {
      return <div>Fetching job run...</div>
    }

    const errorTaskRuns = filterErrorTaskRuns(jobRun)

    return (
      <Grid container spacing={0}>
        <Grid item xs={12}>
          <StatusCard title={jobRun.status}>
            <ul className={classes.list}>
              {errorTaskRuns.map((tr: TaskRun) => (
                <li key={tr.id}>
                  <Typography variant="body1">{tr.result.error}</Typography>
                </li>
              ))}
            </ul>
          </StatusCard>
        </Grid>
      </Grid>
    )
  },
)

interface Props {
  jobSpecId: string
  jobRunId: string
  jobRun?: JobRun
  fetchJobRun: (id: string) => Promise<any>
}

const ShowErrorLog: React.FC<Props> = ({
  jobRunId,
  jobSpecId,
  jobRun,
  fetchJobRun,
}) => {
  useEffect(() => {
    fetchJobRun(jobRunId)
  }, [fetchJobRun, jobRunId])

  return (
    <div>
      <RegionalNav jobSpecId={jobSpecId} jobRunId={jobRunId} jobRun={jobRun} />

      <Content>
        <Details jobRun={jobRun} />
      </Content>
    </div>
  )
}

interface Match {
  params: {
    jobSpecId: string
    jobRunId: string
  }
}

const mapStateToProps = (state: AppState, ownProps: { match: Match }) => {
  const { jobSpecId, jobRunId } = ownProps.match.params
  const jobRun = jobRunSelector(state, jobRunId)

  return {
    jobSpecId,
    jobRunId,
    jobRun,
  }
}

export const ConnectedShowErrorLog = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({ fetchJobRun }),
)(ShowErrorLog)

export default ConnectedShowErrorLog
