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
import Paper from '@material-ui/core/Paper'
import Card from '@material-ui/core/Card'
import Typography from '@material-ui/core/Typography'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableRow from '@material-ui/core/TableRow'
import TimeAgo from '../../components/TimeAgo'
import { getJobRun } from '../../actions/jobRuns'
import { IState } from '../../reducers'
import { JobRun } from '../../entities'

const regionalNavStyles = ({ spacing, palette }: Theme) =>
  createStyles({
    container: {
      padding: spacing.unit * 5,
      paddingBottom: 0
    }
  })

interface IRegionalNavProps extends WithStyles<typeof regionalNavStyles> {
  jobRunId?: string
  jobRun?: IJobRun
}

const RegionalNav = withStyles(regionalNavStyles)(
  ({ jobRunId, jobRun, classes }: IRegionalNavProps) => {
    return (
      <Paper square className={classes.container}>
        <Grid container spacing={0}>
          <Grid item xs={12}>
            <Grid container spacing={0} alignItems="center">
              <Grid item xs={7}>
                <Typography variant="h3" color="secondary" gutterBottom>
                  {jobRunId}
                </Typography>
              </Grid>
            </Grid>
          </Grid>
          <Grid item xs={12}>
            <Typography variant="subtitle2" color="textSecondary">
              {jobRun && (
                <>
                  Created <TimeAgo tooltip={false}>{jobRun.createdAt}</TimeAgo>{' '}
                  ({jobRun.createdAt})
                </>
              )}
            </Typography>
          </Grid>
        </Grid>
      </Paper>
    )
  }
)

interface ISpanRowProps {
  children: React.ReactNode
}

const SpanRow = ({ children }: ISpanRowProps) => (
  <TableRow>
    <TableCell component="th" scope="row" colSpan={3}>
      {children}
    </TableCell>
  </TableRow>
)

const FetchingRow = () => <SpanRow>Loading job run...</SpanRow>

interface IColProps {
  children: React.ReactNode
}

const Col = ({ children }: IColProps) => (
  <TableCell>
    <Typography variant="body1">
      <span>{children}</span>
    </Typography>
  </TableCell>
)

const renderBody = (jobRun?: IJobRun) => {
  if (jobRun) {
    return (
      <>
        <TableRow>
          <Col>Job ID</Col>
          <Col>{jobRun.jobId}</Col>
        </TableRow>
      </>
    )
  } else {
    return <FetchingRow />
  }
}

const styles = ({ spacing, palette }: Theme) =>
  createStyles({
    container: {
      marginTop: spacing.unit * 5,
      padding: spacing.unit * 5
    },
    card: {
      paddingTop: spacing.unit * 2,
      paddingBottom: spacing.unit * 2
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
      <div>
        <RegionalNav jobRunId={jobRunId} jobRun={jobRun} />

        <Grid container spacing={0}>
          <Grid item xs={12}>
            <div className={classes.container}>
              <Card className={classes.card}>
                <Table>
                  <TableBody>{renderBody(jobRun)}</TableBody>
                </Table>
              </Card>
            </div>
          </Grid>
        </Grid>
      </div>
    )
  }
)

const jobRunSelector = (
  { jobRuns }: IState,
  jobRunId?: string
): IJobRun | undefined => {
  if (jobRuns.items) {
    return denormalize(jobRunId, JobRun, { jobRuns: jobRuns.items })
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
