import React from 'react'
import { connect } from 'react-redux'
import { useHooks, useEffect } from 'use-react-hooks'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles
} from '@material-ui/core/styles'
import Grid from '@material-ui/core/Grid'
import Typography from '@material-ui/core/Typography'
import { fetchJobRun } from '../../../actions'
import jobRunSelector from '../../../selectors/jobRun'
import PaddedCard from '../../../components/PaddedCard'
import PrettyJson from '../../../components/PrettyJson'
import matchRouteAndMapDispatchToProps from '../../../utils/matchRouteAndMapDispatchToProps'
import Content from '../../../components/Content'
import RegionalNav from '../../../components/JobRuns/RegionalNav'
import StatusCard from '../../../components/JobRuns/StatusCard'

const filterErrorTaskRuns = jobRun => {
  return jobRun.taskRuns.filter(tr => {
    return tr.status === 'errored'
  })
}

const detailsStyles = ({ spacing }: Theme) => ({
  list: {
    marginTop: spacing.unit * 4
  }
})

interface IDetailsProps extends WithStyles<typeof detailsStyles> {
  fetching: boolean
  jobRun?: any
}

const Details = withStyles(detailsStyles)(
  ({ fetching, jobRun, classes }: IDetailsProps) => {
    if (fetching || !jobRun) {
      return <div>Fetching job run...</div>
    }

    const errorTaskRuns = filterErrorTaskRuns(jobRun)

    return (
      <Grid container spacing={0}>
        <Grid item xs={12}>
          <StatusCard title={jobRun.status}>
            <ul className={classes.list}>
              {errorTaskRuns.map(tr => (
                <li key={tr.id}>
                  <Typography variant="body1">{tr.result.error}</Typography>
                </li>
              ))}
            </ul>
          </StatusCard>
        </Grid>
      </Grid>
    )
  }
)

const styles = (theme: Theme) => ({})

interface IProps extends WithStyles<typeof styles> {
  fetching: boolean
  jobSpecId: string
  jobRunId: string
  jobRun?: any
  fetchJobRun: (string) => Promise<any>
}

const ShowErrorLog = useHooks(
  ({ fetching, jobRunId, jobSpecId, jobRun, fetchJobRun }: IProps) => {
    useEffect(() => {
      fetchJobRun(jobRunId)
    }, [jobRunId])

    return (
      <div>
        <RegionalNav
          jobSpecId={jobSpecId}
          jobRunId={jobRunId}
          jobRun={jobRun}
        />

        <Content>
          <Details fetching={fetching} jobRun={jobRun} />
        </Content>
      </div>
    )
  }
)

const mapStateToProps = (state, ownProps) => {
  const { jobSpecId, jobRunId } = ownProps.match.params
  const jobRun = jobRunSelector(state, jobRunId)
  const fetching = state.jobRuns.fetching

  return {
    jobSpecId,
    jobRunId,
    jobRun,
    fetching
  }
}

export const ConnectedShowErrorLog = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({ fetchJobRun })
)(ShowErrorLog)

export default withStyles(styles)(ConnectedShowErrorLog)
