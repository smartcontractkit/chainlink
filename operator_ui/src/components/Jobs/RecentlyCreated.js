import { SimpleListCard } from 'components/SimpleListCard'
import { SimpleListCardItem } from 'components/SimpleListCardItem'
import { TimeAgo } from 'components/TimeAgo'
import Grid from '@material-ui/core/Grid'
import { withStyles } from '@material-ui/core/styles'
import TableCell from '@material-ui/core/TableCell'
import TableRow from '@material-ui/core/TableRow'
import Typography from '@material-ui/core/Typography'
import Link from 'components/Link'
import PropTypes from 'prop-types'
import React from 'react'

const styles = () => ({
  block: { display: 'block' },
  overflowEllipsis: { textOverflow: 'ellipsis', overflow: 'hidden' },
})

const RecentlyCreated = ({ jobs, classes }) => {
  let status

  if (!jobs) {
    status = (
      <TableRow>
        <TableCell scope="row">
          <Typography variant="body1" color="textSecondary">
            ...
          </Typography>
        </TableCell>
      </TableRow>
    )
  } else if (jobs.length === 0) {
    status = (
      <TableRow>
        <TableCell scope="row">
          <Typography variant="body1" color="textSecondary">
            No recently created jobs
          </Typography>
        </TableCell>
      </TableRow>
    )
  } else {
    status = (
      <React.Fragment>
        {jobs.map((j) => (
          <SimpleListCardItem key={j.id}>
            <Grid container spacing={0}>
              <Grid item xs={12}>
                <Link
                  href={`/jobs/${j.id}`}
                  classes={{ linkContent: classes.block }}
                >
                  <Typography
                    className={classes.overflowEllipsis}
                    variant="body1"
                    component="span"
                    color="primary"
                  >
                    {j.id}
                  </Typography>
                </Link>
              </Grid>
              <Grid item xs={12}>
                <Typography variant="body1" color="textSecondary">
                  Created <TimeAgo tooltip>{j.createdAt}</TimeAgo>
                </Typography>
              </Grid>
            </Grid>
          </SimpleListCardItem>
        ))}
      </React.Fragment>
    )
  }

  return <SimpleListCard title="Recently Created Jobs">{status}</SimpleListCard>
}

RecentlyCreated.propTypes = {
  jobs: PropTypes.array,
}

export default withStyles(styles)(RecentlyCreated)
