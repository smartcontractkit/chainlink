import { KeyValueList, PaddedCard } from '@chainlink/styleguide'
import Grid from '@material-ui/core/Grid'
import Typography from '@material-ui/core/Typography'
import { fetchConfiguration } from 'actionCreators'
import Content from 'components/Content'
import DeleteJobRuns from 'pages/Configuration/DeleteJobRuns'
import PropTypes from 'prop-types'
import React, { useEffect } from 'react'
import { connect } from 'react-redux'
import configurationSelector from 'selectors/configuration'
import extractBuildInfo from 'utils/extractBuildInfo'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'

const buildInfo = extractBuildInfo()

export const Configuration = ({ fetchConfiguration, data }) => {
  useEffect(() => {
    document.title = 'Configuration'
    fetchConfiguration()
  }, [fetchConfiguration])

  return (
    <Content>
      <Grid container>
        <Grid item sm={12} md={8}>
          <KeyValueList title="Configuration" entries={data} showHead />
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
}

Configuration.propTypes = {
  data: PropTypes.array.isRequired,
}

const mapStateToProps = (state) => {
  const data = configurationSelector(state)
  return { data }
}

export const ConnectedConfiguration = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({ fetchConfiguration }),
)(Configuration)

export default ConnectedConfiguration
