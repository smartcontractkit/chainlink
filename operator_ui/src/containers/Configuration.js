import React from 'react'
import PropTypes from 'prop-types'
import { useHooks, useEffect } from 'use-react-hooks'
import { connect } from 'react-redux'
import Grid from '@material-ui/core/Grid'
import Typography from '@material-ui/core/Typography'
import PaddedCard from '@chainlink/styleguide/components/PaddedCard'
import KeyValueList from '@chainlink/styleguide/components/KeyValueList'
import Content from 'components/Content'
import DeleteJobRuns from 'containers/Configuration/DeleteJobRuns'
import { fetchConfiguration } from 'actions'
import configurationSelector from 'selectors/configuration'
import extractBuildInfo from 'utils/extractBuildInfo'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'

const buildInfo = extractBuildInfo()

export const Configuration = useHooks(props => {
  useEffect(() => {
    document.title = 'Configuration'
    props.fetchConfiguration()
  }, [])

  return (
    <Content>
      <Grid container>
        <Grid item sm={12} md={8}>
          <KeyValueList title="Configuration" entries={props.data} showHead />
        </Grid>
        <Grid item sm={12} md={4}>
          <Grid container>
            <Grid item xs={12}>
              <PaddedCard>
                <Typography variant="h5" color="secondary">
                  Version
                </Typography>
                <Typography variant="body1" color="textSecondary">
                  {buildInfo.version}
                </Typography>
              </PaddedCard>
            </Grid>
            <Grid item xs={12}>
              <PaddedCard>
                <Typography variant="h5" color="secondary">
                  SHA
                </Typography>
                <Typography variant="body1" color="textSecondary">
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
  data: PropTypes.array.isRequired
}

const mapStateToProps = state => {
  const data = configurationSelector(state)
  return { data }
}

export const ConnectedConfiguration = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({ fetchConfiguration })
)(Configuration)

export default ConnectedConfiguration
