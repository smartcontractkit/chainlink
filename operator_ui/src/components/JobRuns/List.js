import React from 'react'
import PropTypes from 'prop-types'
import { withStyles } from '@material-ui/core/styles'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import Typography from '@material-ui/core/Typography'
import TableRow from '@material-ui/core/TableRow'
import TableCell from '@material-ui/core/TableCell'
import Card from '@material-ui/core/Card'
import Link from 'components/Link'
import TimeAgo from 'components/TimeAgo'
import ReactStaticLinkComponent from '../ReactStaticLinkComponent'
import classNames from 'classnames'
import Button from 'components/Button'
import titleize from 'utils/titleize'

const styles = theme => {
  return {
    jobRunsCard: {
      overflow: 'auto'
    },
    idCell: {
      width: '40%'
    },
    stampCell: {
      width: '30%'
    },
    statusCell: {
      textAlign: 'end',
      width: '30%'
    },
    runDetails: {
      paddingTop: theme.spacing.unit * 2,
      paddingBottom: theme.spacing.unit * 2,
      paddingLeft: theme.spacing.unit * 2
    },
    stamp: {
      paddingLeft: theme.spacing.unit
    },
    status: {
      paddingLeft: theme.spacing.unit * 1.5,
      paddingRight: theme.spacing.unit * 1.5,
      paddingTop: theme.spacing.unit / 2,
      paddingBottom: theme.spacing.unit / 2,
      borderRadius: theme.spacing.unit * 2,
      marginRight: theme.spacing.unit,
      width: 'fit-content',
      display: 'inline-block'
    },
    errored: {
      backgroundColor: theme.palette.error.light,
      color: theme.palette.error.main
    },
    pending: {
      backgroundColor: theme.palette.listPendingStatus.background,
      color: theme.palette.listPendingStatus.color
    },
    completed: {
      backgroundColor: theme.palette.listCompletedStatus.background,
      color: theme.palette.listCompletedStatus.color
    },
    noRuns: {
      padding: theme.spacing.unit * 2
    }
  }
}

const classFromStatus = (classes, status) => {
  if (
    !status ||
    status.startsWith('pending') ||
    status.startsWith('in_progress')
  ) {
    return classes['pending']
  }
  return classes[status.toLowerCase()]
}

const renderRuns = (runs, classes) => {
  if (runs && runs.length === 0) {
    return (
      <TableRow>
        <TableCell colSpan={5}>
          <Typography
            variant="body1"
            color="textSecondary"
            className={classes.noRuns}
          >
            The job hasn’t run yet
          </Typography>
        </TableCell>
      </TableRow>
    )
  } else if (runs) {
    return runs.map(r => (
      <TableRow key={r.id}>
        <TableCell className={classes.idCell} scope="row">
          <div className={classes.runDetails}>
            <Link to={`/jobs/${r.jobId}/runs/id/${r.id}`}>
              <Typography variant="h5" color="primary" component="span">
                {r.id}
              </Typography>
            </Link>
          </div>
        </TableCell>
        <TableCell className={classes.stampCell}>
          <Typography
            variant="body1"
            color="textSecondary"
            className={classes.stamp}
          >
            Created <TimeAgo tooltip>{r.createdAt}</TimeAgo>
          </Typography>
        </TableCell>
        <TableCell className={classes.statusCell} scope="row">
          <Typography
            variant="body1"
            className={classNames(
              classes.status,
              classFromStatus(classes, r.status)
            )}
          >
            {titleize(r.status)}
          </Typography>
        </TableCell>
      </TableRow>
    ))
  }

  return (
    <TableRow>
      <TableCell colSpan={5}>...</TableCell>
    </TableRow>
  )
}

const List = ({ jobSpecId, runs, count, showJobRunsCount, classes }) => {
  return (
    <Card className={classes.jobRunsCard}>
      <Table padding="none">
        <TableBody>
          {renderRuns(runs, classes)}
          {runs && count > showJobRunsCount && (
            <TableRow>
              <TableCell>
                <div className={classes.runDetails}>
                  <Button
                    to={`/jobs/${jobSpecId}/runs`}
                    component={ReactStaticLinkComponent}
                  >
                    View More
                  </Button>
                </div>
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </Card>
  )
}

List.propTypes = {
  jobSpecId: PropTypes.string.isRequired,
  runs: PropTypes.array.isRequired,
  count: PropTypes.number.isRequired,
  showJobRunsCount: PropTypes.any.isRequired
}

export default withStyles(styles)(List)
