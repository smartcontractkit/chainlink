import Card from '@material-ui/core/Card'
import Grid from '@material-ui/core/Grid'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles
} from '@material-ui/core/styles'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableRow from '@material-ui/core/TableRow'
import React, { useEffect } from 'react'
import { connect, MapDispatchToProps, MapStateToProps } from 'react-redux'
import { bindActionCreators } from 'redux'
import build from 'redux-object'
import { getJobRun } from '../../actions/jobRuns'
import Details from '../../components/JobRuns/Details'
import RegionalNav from '../../components/JobRuns/RegionalNav'
import RunStatus from '../../components/JobRuns/RunStatus'
import { IState as State } from '../../reducers'

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

const styles = ({ spacing, breakpoints }: Theme) =>
  createStyles({
    container: {
      paddingTop: spacing.unit * 2,
      paddingBottom: spacing.unit * 2,
      paddingLeft: spacing.unit * 2,
      paddingRight: spacing.unit * 2,
      [breakpoints.up('sm')]: {
        paddingTop: spacing.unit * 3,
        paddingBottom: spacing.unit * 3,
        paddingLeft: spacing.unit * 3,
        paddingRight: spacing.unit * 3
      }
    },
    card: {
      paddingTop: spacing.unit,
      paddingBottom: spacing.unit
    }
  })

interface OwnProps {
  path: string
  jobRunId?: string
}

interface StateProps {
  jobRun?: IJobRun
  etherscanHost?: string
}

interface DispatchProps {
  getJobRun: any
}

interface IProps
  extends WithStyles<typeof styles>,
    OwnProps,
    StateProps,
    DispatchProps {}

const Show = withStyles(styles)(
  ({ jobRunId, jobRun, getJobRun, classes, etherscanHost }: IProps) => {
    useEffect(() => {
      getJobRun(jobRunId)
    }, [getJobRun, jobRunId])

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
                {jobRun && etherscanHost ? (
                  <Details
                    jobRun={jobRun}
                    etherscanHost={etherscanHost.toString()}
                  />
                ) : (
                  <Loading />
                )}
              </Card>
            </div>
          </Grid>
        </Grid>
      </>
    )
  }
)

const jobRunSelector = (
  { jobRuns, taskRuns, chainlinkNodes }: State,
  jobRunId?: string
): IJobRun | undefined => {
  if (jobRuns.items) {
    const document = {
      jobRuns: jobRuns.items,
      taskRuns: taskRuns.items,
      chainlinkNodes: chainlinkNodes.items
    }
    return build(document, 'jobRuns', jobRunId)
  }
}

const mapStateToProps: MapStateToProps<StateProps, OwnProps, State> = (
  state,
  { jobRunId }
) => {
  const jobRun = jobRunSelector(state, jobRunId)
  const etherscanHost = state.config.etherscanHost

  return { jobRun, etherscanHost }
}

const mapDispatchToProps: MapDispatchToProps<
  DispatchProps,
  OwnProps
> = dispatch => bindActionCreators({ getJobRun }, dispatch)

const ConnectedShow = connect(
  mapStateToProps,
  mapDispatchToProps
)(Show)

export default withStyles(styles)(ConnectedShow)
