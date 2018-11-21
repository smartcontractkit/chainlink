import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { connect } from 'react-redux'
import { withSiteData } from 'react-static'
import Grid from '@material-ui/core/Grid'
import Title from 'components/Title'
import ConfigList from 'components/ConfigList'
import Content from 'components/Content'
import { fetchConfiguration } from 'actions'
import configsSelector from 'selectors/configs'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'

export class Configuration extends Component {
  componentDidMount () {
    this.props.fetchConfiguration()
  }

  render () {
    const {props} = this

    return (
      <Content>
        <Title>Configuration</Title>

        <Grid container spacing={40}>
          <Grid item xs={12}>
            <ConfigList
              configs={props.configs}
              error={props.error}
            />
          </Grid>
        </Grid>
      </Content>
    )
  }
}

Configuration.propTypes = {
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

export default withSiteData(ConnectedConfiguration)
