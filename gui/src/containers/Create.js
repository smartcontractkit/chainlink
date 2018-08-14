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

const styles = theme => ({
  root: {
    flexGrow: 1,
    backgroundColor: theme.palette.background.paper,
    paddingTop: theme.spacing.unit * 2
  },
  card: {
    paddingBottom: theme.spacing.unit * 2
  },
  tabPadding: {
    padding: theme.spacing.unit * 3
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

class Create extends React.Component {
  state = {
    structure: 'bridge',
    value: 0
  };

  componentDidMount () {
    // Need to set int value because <Tabs/> component
    // will not focus on tab with strings
    if (this.props.match) {
      const value = this.props.match.params.structure === 'job' ? 1 : 0
      this.setState({ value: value })
    }
  }

  handleChange = (event, value) => {
    this.setState({ value })
    if (this.props.history) {
      switch (value) {
        case 0:
          this.props.history.push('/create/bridge')
          break
        case 1:
          this.props.history.push('/create/job')
          break
      }
    }
  };

  render () {
    const { classes } = this.props
    const matchExists = this.props.match && this.props.match.params
    const structure = (matchExists && this.props.match.params.structure) || 'bridge'
    const { value } = this.state
    return (
      <div className={classes.root}>
        <Card className={classes.card}>
          <AppBar position='static'>
            <Tabs value={value} onChange={this.handleChange}>
              <Tab label='Create Bridge' />
              <Tab label='Create Job' />
            </Tabs>
          </AppBar>
          {structure === 'bridge' && <TabContainer><BridgeForm /></TabContainer>}
          {structure === 'job' && <TabContainer><JobForm /></TabContainer>}
        </Card>
      </div>
    )
  }
}

Create.propTypes = {
  classes: PropTypes.object.isRequired
}

export const withoutStyles = Create
export default withStyles(styles)(Create)
