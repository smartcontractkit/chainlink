import React from 'react'
import PropTypes from 'prop-types'
import { withStyles } from '@material-ui/core/styles'
import AppBar from '@material-ui/core/AppBar'
import Tabs from '@material-ui/core/Tabs'
import Tab from '@material-ui/core/Tab'
import Typography from '@material-ui/core/Typography'
import BridgeForm from 'components/BridgeForm'
import JobForm from 'components/JobForm'
import { Card } from '@material-ui/core'

const TabContainer = props => {
  return (
    <Typography component='div' style={{ padding: 8 * 3 }}>
      {props.children}
    </Typography>
  )
}

TabContainer.propTypes = {
  children: PropTypes.node.isRequired
}

const styles = theme => ({
  root: {
    flexGrow: 1,
    backgroundColor: theme.palette.background.paper
  },
  card: {
    paddingBottom: theme.spacing.unit * 2
  }
})

class WrappedTabs extends React.Component {
  state = {
    value: 0
  };

  handleChange = (event, value) => {
    this.setState({ value })
  };

  render () {
    const { classes } = this.props
    const { value } = this.state

    return (
      <div className={classes.root}>
        <br />
        <Card className={classes.card}>
          <AppBar position='static'>
            <Tabs value={value} onChange={this.handleChange}>
              <Tab label='Create Bridge' />
              <Tab label='Create Job' />
            </Tabs>
          </AppBar>
          {value === 0 && <TabContainer><BridgeForm /></TabContainer>}
          {value === 1 && <TabContainer><JobForm /></TabContainer>}
        </Card>
      </div>
    )
  }
}

WrappedTabs.propTypes = {
  classes: PropTypes.object.isRequired
}

export default withStyles(styles)(WrappedTabs)
