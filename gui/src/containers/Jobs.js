import React from 'react'
import PropType from 'prop-types'
import Grid from '@material-ui/core/Grid'
import Typography from '@material-ui/core/Typography'
import JobList from 'containers/JobList'
import AccountBalance from 'containers/AccountBalance'
import { withSiteData } from 'react-static'
import { withStyles } from '@material-ui/core/styles'

const styles = theme => ({
  title: {
    marginTop: theme.spacing.unit * 5,
    marginBottom: theme.spacing.unit * 5
  }
})

const Jobs = ({classes}) => (
  <div>
    <Typography variant='display2' color='inherit' className={classes.title}>
      Jobs
    </Typography>

    <Grid container spacing={40}>
      <Grid item xs={9}>
        <JobList />
      </Grid>
      <Grid item xs={3}>
        <AccountBalance />
      </Grid>
    </Grid>
  </div>
)

Jobs.propTypes = {
  classes: PropType.object.isRequired
}

export default withSiteData(withStyles(styles)(Jobs))
