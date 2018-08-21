import React from 'react'
import PropTypes from 'prop-types'
import { withStyles } from '@material-ui/core/styles'
import AppBar from '@material-ui/core/AppBar'
import Tabs from '@material-ui/core/Tabs'
import Tab from '@material-ui/core/Tab'
import Typography from '@material-ui/core/Typography'
import JobForm from 'components/JobForm'
import { Card } from '@material-ui/core'

const styles = theme => ({
  root: {
    flexGrow: 1,
    backgroundColor: 'transparent',
    paddingTop: theme.spacing.unit * 2
  },
  card: {
    paddingBottom: theme.spacing.unit * 2
  }
})

const TabContainer = (props, classes) => {
  return (
    <Typography component='div' className={classes.padding}>
      {props.children}
    </Typography>
  )
}

TabContainer.propTypes = {
  children: PropTypes.node.isRequired
}

const CreateJobSpec = props => {
  const { classes } = props
  const passJsonState = props.location && props.location.state
  return (
    <div className={classes.root}>
      <Card className={classes.card}>
        <AppBar position='static'>
          <Tabs value={0}>
            <Tab label='Create Job' />
          </Tabs>
        </AppBar>
        <TabContainer>
          <JobForm {...passJsonState} />
        </TabContainer>
      </Card>
    </div>
  )
}

CreateJobSpec.propTypes = {
  classes: PropTypes.object.isRequired
}

export default withStyles(styles)(CreateJobSpec)
