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
import PaddedCard from '@chainlink/styleguide/src/components/PaddedCard'
import { fetchJobRun } from '../../../actions'
import jobRunSelector from '../../../selectors/jobRun'
import PrettyJson from '../../../components/PrettyJson'
import matchRouteAndMapDispatchToProps from '../../../utils/matchRouteAndMapDispatchToProps'
import Content from '../../../components/Content'
import RegionalNav from './RegionalNav'
import StatusCard from '../../../components/JobRuns/StatusCard'
import { IJobRun, ITaskRun } from '../../../../@types/operator_ui'
import { IState } from '../../../connectors/redux/reducers'

const filterErrorTaskRuns = (jobRun: IJobRun) => {
  return jobRun.taskRuns.filter((tr: ITaskRun) => {
    return tr.status === 'errored'
  })
}

const detailsStyles = ({ spacing }: Theme) => ({
  list: {
    marginTop: spacing.unit * 4
  }
})

interface IDetailsProps extends WithStyles<typeof detailsStyles> {
  jobRun?: IJobRun
}

const Details = withStyles(detailsStyles)(
  ({ jobRun, classes }: IDetailsProps) => {
    if (!jobRun) {
      return <div>Fetching job run...</div>
    }

    const errorTaskRuns = filterErrorTaskRuns(jobRun)

    return (
      <Grid container spacing={0}>
        <Grid item xs={12}>
          <StatusCard title={jobRun.status}>
            <ul className={classes.list}>
              {errorTaskRuns.map((tr: ITaskRun) => (
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
  jobSpecId: string
  jobRunId: string
  jobRun?: IJobRun
  fetchJobRun: (id: string) => Promise<any>
}

const ShowErrorLog = useHooks(
  ({ jobRunId, jobSpecId, jobRun, fetchJobRun }: IProps) => {
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
          <Details jobRun={jobRun} />
        </Content>
      </div>
    )
  }
)

interface Match {
  params: {
    jobSpecId: string
    jobRunId: string
  }
}

const mapStateToProps = (state: IState, ownProps: { match: Match }) => {
  const { jobSpecId, jobRunId } = ownProps.match.params
  const jobRun = jobRunSelector(state, jobRunId)

  return {
    jobSpecId,
    jobRunId,
    jobRun
  }
}

export const ConnectedShowErrorLog = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({ fetchJobRun })
)(ShowErrorLog)

export default withStyles(styles)(ConnectedShowErrorLog)
