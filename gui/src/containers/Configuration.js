import React from 'react'
import PropTypes from 'prop-types'
import { useHooks, useEffect } from 'use-react-hooks'
import { connect } from 'react-redux'
import Grid from '@material-ui/core/Grid'
import Typography from '@material-ui/core/Typography'
import PaddedCard from 'components/PaddedCard'
import KeyValueList from 'components/KeyValueList'
import Content from 'components/Content'
import DeleteJobRuns from 'containers/Configuration/DeleteJobRuns'
import { fetchConfiguration } from 'actions'
import configsSelector from 'selectors/configs'
import extractBuildInfo from 'utils/extractBuildInfo'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'

const buildInfo = extractBuildInfo()

export const Configuration = useHooks(props => {
  useEffect(() => { props.fetchConfiguration() }, [])

  return (
    <Content>
      <Grid container>
        <Grid item sm={12} md={8}>
          <KeyValueList
            title='Configuration'
            entries={props.configs}
            error={props.error}
            showHead
          />
        </Grid>
        <Grid item sm={12} md={4}>
          <Grid container>
            <Grid item xs={12}>
              <PaddedCard>
                <Typography variant='h5' color='secondary'>
                  Version
                </Typography>
                <Typography variant='body1' color='textSecondary'>
                  {buildInfo.version}
                </Typography>
              </PaddedCard>
            </Grid>
            <Grid item xs={12}>
              <PaddedCard>
                <Typography variant='h5' color='secondary'>
                  SHA
                </Typography>
                <Typography variant='body1' color='textSecondary'>
                  {buildInfo.sha}
                </Typography>
              </PaddedCard>
            </Grid>
            <Grid item xs={12}>
              <DeleteJobRuns />
            </Grid>
          </Grid>
        </Grid>
      </Grid>
    </Content>
  )
})

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
  matchRouteAndMapDispatchToProps({ fetchConfiguration })
)(Configuration)

export default ConnectedConfiguration
