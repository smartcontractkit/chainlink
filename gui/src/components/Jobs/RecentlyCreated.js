import React from 'react'
import PropTypes from 'prop-types'
import Typography from '@material-ui/core/Typography'
import Grid from '@material-ui/core/Grid'
import CardContent from '@material-ui/core/CardContent'
import TableCell from '@material-ui/core/TableCell'
import TableRow from '@material-ui/core/TableRow'
import TimeAgo from 'components/TimeAgo'
import Link from 'components/Link'
import SimpleListCard from 'components/Cards/SimpleList'
import SimpleListCardItem from 'components/Cards/SimpleListItem'

const RecentlyCreated = ({ jobs }) => {
  let status

  if (!jobs) {
    status = (
      <TableRow>
        <TableCell scope="row">
          <CardContent>
            <Typography variant="body1" color="textSecondary">
              ...
            </Typography>
          </CardContent>
        </TableCell>
      </TableRow>
    )
  } else if (jobs.length === 0) {
    status = (
      <TableRow>
        <TableCell scope="row">
          <CardContent>
            <Typography variant="body1" color="textSecondary">
              No recently created jobs
            </Typography>
          </CardContent>
        </TableCell>
      </TableRow>
    )
  } else {
    status = (
      <React.Fragment>
        {jobs.map(j => (
          <SimpleListCardItem key={j.id}>
            <Grid container spacing={0}>
              <Grid item xs={12}>
                <Link to={`/jobs/${j.id}`}>
                  <Typography variant="body1" component="span" color="primary">
                    {j.id}
                  </Typography>
                </Link>
              </Grid>
              <Grid item xs={12}>
                <Typography variant="body1" color="textSecondary">
                  Created <TimeAgo>{j.createdAt}</TimeAgo>
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
  jobs: PropTypes.array
}

export default RecentlyCreated
