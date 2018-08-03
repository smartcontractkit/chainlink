import React, { Component } from 'react'
import PropTypes from 'prop-types'
import Grid from '@material-ui/core/Grid'
import Typography from '@material-ui/core/Typography'
import ConfigList from 'components/ConfigList'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import { withSiteData } from 'react-static'
import { withStyles } from '@material-ui/core/styles'
import { connect } from 'react-redux'
import { fetchConfiguration } from 'actions'
import { configsSelector } from 'selectors'

const styles = theme => ({
  title: {
    marginTop: theme.spacing.unit * 5,
    marginBottom: theme.spacing.unit * 5
  }
})

const renderConfigList = ({configs, error}) => (
  <ConfigList
    configs={configs}
    error={error}
  />
)

export class Configuration extends Component {
  componentDidMount () {
    this.props.fetchConfiguration()
  }

  render () {
    const { classes } = this.props

    return (
      <div>
        <Typography variant='display2' color='inherit' className={classes.title}>
          Configuration
        </Typography>

        <Grid container spacing={40}>
          <Grid item xs={12}>
            {renderConfigList(this.props)}
          </Grid>
        </Grid>
      </div>
    )
  }
}

Configuration.propTypes = {
  classes: PropTypes.object.isRequired,
  configs: PropTypes.array.isRequired,
  error: PropTypes.string
}

const mapStateToProps = state => {
  let configError
  if (state.configuration.networkError) {
    configError = 'There was an error fetching the configuration. Please reload the page.'
  }

  return {
    configs: configsSelector(state),
    error: configError
  }
}

export const ConnectedConfiguration = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({fetchConfiguration})
)(Configuration)

export default withSiteData(withStyles(styles)(ConnectedConfiguration))
