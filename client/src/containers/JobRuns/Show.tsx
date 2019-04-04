import React, { useEffect } from 'react'
import { denormalize } from 'normalizr'
import { connect } from 'react-redux'
import { bindActionCreators, Dispatch } from 'redux'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles
} from '@material-ui/core/styles'
import Grid from '@material-ui/core/Grid'
import Card from '@material-ui/core/Card'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableRow from '@material-ui/core/TableRow'
import Details from '../../components/JobRuns/Details'
import RegionalNav from '../../components/JobRuns/RegionalNav'
import RunStatus from '../../components/JobRuns/RunStatus'
import { getJobRun } from '../../actions/jobRuns'
import { IState } from '../../reducers'
import { JobRun } from '../../entities'

const Loading = () => (
  <Table>
    <TableBody>
      <TableRow>
        <TableCell component="th" scope="row" colSpan={3}>
          Loading job run...
        </TableCell>
      </TableRow>
    </TableBody>
  </Table>
)

const styles = ({ spacing }: Theme) =>
  createStyles({
    container: {
      padding: spacing.unit * 5,
      paddingBottom: 0
    },
    card: {
      paddingTop: spacing.unit,
      paddingBottom: spacing.unit
    }
  })

interface IProps extends WithStyles<typeof styles> {
  jobRunId?: string
  jobRun?: IJobRun
  getJobRun: Function
  path: string
}

const Show = withStyles(styles)(
  ({ jobRunId, jobRun, getJobRun, classes }: IProps) => {
    useEffect(() => {
      getJobRun(jobRunId)
    }, [])

    return (
      <>
        <RegionalNav jobRun={jobRun} />

        <Grid container spacing={0}>
          <Grid item xs={12}>
            {jobRun && (
              <div className={classes.container}>
                <RunStatus jobRun={jobRun} />
              </div>
            )}

            <div className={classes.container}>
              <Card className={classes.card}>
                {jobRun ? <Details jobRun={jobRun} /> : <Loading />}
              </Card>
            </div>
          </Grid>
        </Grid>
      </>
    )
  }
)

const jobRunSelector = (
  { jobRuns, taskRuns }: IState,
  jobRunId?: string
): IJobRun | undefined => {
  if (jobRuns.items) {
    return denormalize(jobRunId, JobRun, {
      jobRuns: jobRuns.items,
      taskRuns: taskRuns.items
    })
  }
}

interface IOwnProps {
  jobRunId?: string
}

const mapStateToProps = (state: IState, { jobRunId }: IOwnProps) => {
  return {
    jobRun: jobRunSelector(state, jobRunId)
  }
}

const mapDispatchToProps = (dispatch: Dispatch<any>) =>
  bindActionCreators({ getJobRun }, dispatch)

const ConnectedShow = connect(
  mapStateToProps,
  mapDispatchToProps
)(Show)

export default withStyles(styles)(ConnectedShow)
