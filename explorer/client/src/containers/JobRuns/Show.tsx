import React, { useEffect } from 'react'
import build from 'redux-object'
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

interface IProps extends WithStyles<typeof styles> {
  jobRunId?: string
  jobRun?: IJobRun
  etherscanHost?: string
  getJobRun: Function
  path: string
}

const Show = withStyles(styles)(
  ({ jobRunId, jobRun, getJobRun, classes, etherscanHost }: IProps) => {
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
  { jobRuns, taskRuns, chainlinkNodes }: IState,
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

const etherscanHostSelector = ({ config }: IState) => {
  return config.etherscanHost
}

interface IOwnProps {
  jobRunId?: string
}

const mapStateToProps = (state: IState, { jobRunId }: IOwnProps) => {
  const jobRun = jobRunSelector(state, jobRunId)
  const etherscanHost = etherscanHostSelector(state)

  return { jobRun, etherscanHost }
}

const mapDispatchToProps = (dispatch: Dispatch<any>) =>
  bindActionCreators({ getJobRun }, dispatch)

const ConnectedShow = connect(
  mapStateToProps,
  mapDispatchToProps
)(Show)

export default withStyles(styles)(ConnectedShow)
