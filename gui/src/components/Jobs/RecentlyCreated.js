import React from 'react'
import PropTypes from 'prop-types'
import { withStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import Grid from '@material-ui/core/Grid'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableRow from '@material-ui/core/TableRow'
import CardContent from '@material-ui/core/CardContent'
import PaddedTitleCard from 'components/PaddedTitleCard'
import TimeAgo from 'components/TimeAgo'
import Link from 'components/Link'

const styles = theme => ({
  cell: {
    borderColor: theme.palette.divider,
    borderTop: `1px solid`,
    borderBottom: 'none',
    paddingTop: theme.spacing.unit,
    paddingBottom: theme.spacing.unit * 2
  }
})

const RecentlyCreated = ({classes, jobs}) => {
  let status

  if (!jobs) {
    status = (
      <CardContent>
        <Typography variant='body1' color='textSecondary'>...</Typography>
      </CardContent>
    )
  } else if (jobs.length === 0) {
    status = (
      <CardContent>
        <Typography variant='body1' color='textSecondary'>
          No recently created jobs
        </Typography>
      </CardContent>
    )
  } else {
    status = (
      <Table>
        <TableBody>
          {jobs.map(j => (
            <TableRow key={j.id}>
              <TableCell scope='row' className={classes.cell}>
                <Grid container>
                  <Grid item xs={12}>
                    <Link to={`/jobs/${j.id}`}>{j.id}</Link>
                  </Grid>
                  <Grid item xs={12}>
                    <Typography variant='body1' color='textSecondary'>
                      Created <TimeAgo>{j.createdAt}</TimeAgo>
                    </Typography>
                  </Grid>
                </Grid>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    )
  }

  return (
    <PaddedTitleCard title='Recently Created Jobs'>
      {status}
    </PaddedTitleCard>
  )
}

RecentlyCreated.propTypes = {
  jobs: PropTypes.array
}

export default withStyles(styles)(RecentlyCreated)
