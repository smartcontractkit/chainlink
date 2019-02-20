import React from 'react'
import PropTypes from 'prop-types'
import { withStyles } from '@material-ui/core/styles'
import Table from '@material-ui/core/Table'
import TableHead from '@material-ui/core/TableHead'
import TableBody from '@material-ui/core/TableBody'
import Typography from '@material-ui/core/Typography'
import TableRow from '@material-ui/core/TableRow'
import TableCell from '@material-ui/core/TableCell'
import Card from '@material-ui/core/Card'
import Link from 'components/Link'
import TimeAgo from 'components/TimeAgo'

const styles = () => ({
  jobRunsCard: {
    overflow: 'auto'
  }
})

const statusColor = status => {
  if (status === 'error') {
    return 'error'
  }
}

const renderRuns = runs => {
  if (runs && runs.length === 0) {
    return (
      <TableRow>
        <TableCell colSpan={5}>The job hasn't run yet</TableCell>
      </TableRow>
    )
  } else if (runs) {
    return runs.map(r => (
      <TableRow key={r.id}>
        <TableCell component="th" scope="row">
          <Link to={`/jobs/${r.jobId}/runs/id/${r.id}`}>{r.id}</Link>
        </TableCell>
        <TableCell component="th" scope="row">
          <Typography variant="body1" color={statusColor(r.status)}>
            {r.status}
          </Typography>
        </TableCell>
        <TableCell>
          <Typography variant="body1">
            <TimeAgo>{r.createdAt}</TimeAgo>
          </Typography>
        </TableCell>
        <TableCell>
          <Typography variant="body1">
            {JSON.stringify(r.result.data)}
          </Typography>
        </TableCell>
        <TableCell>
          <Typography variant="body1" color="error">
            {r.result.error && JSON.stringify(r.result.error)}
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

const List = ({ runs, classes }) => (
  <Card className={classes.jobRunsCard}>
    <Table>
      <TableHead>
        <TableRow>
          <TableCell>
            <Typography variant="body1" color="textSecondary">
              ID
            </Typography>
          </TableCell>
          <TableCell>
            <Typography variant="body1" color="textSecondary">
              Status
            </Typography>
          </TableCell>
          <TableCell>
            <Typography variant="body1" color="textSecondary">
              Created
            </Typography>
          </TableCell>
          <TableCell>
            <Typography variant="body1" color="textSecondary">
              Result
            </Typography>
          </TableCell>
          <TableCell>
            <Typography variant="body1" color="textSecondary">
              Error
            </Typography>
          </TableCell>
        </TableRow>
      </TableHead>
      <TableBody>{renderRuns(runs)}</TableBody>
    </Table>
  </Card>
)

List.propTypes = {
  jobSpecId: PropTypes.string.isRequired,
  runs: PropTypes.array.isRequired
}

export default withStyles(styles)(List)
