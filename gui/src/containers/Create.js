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
    padding: 24
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
    value: 0
  };

  componentDidMount() {
    if (this.props.location && this.props.location.state)
      this.setState({ value: this.props.location.state.tab })
  }

  handleChange = (event, value) => {
    this.setState({ value })
  };

  render () {
    const { classes } = this.props
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
          {value === 0 && <TabContainer><BridgeForm /></TabContainer>}
          {value === 1 && <TabContainer><JobForm /></TabContainer>}
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
