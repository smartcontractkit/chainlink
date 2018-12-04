import React from 'react'
import { useHooks, useEffect } from 'use-react-hooks'
import PropTypes from 'prop-types'
import { connect } from 'react-redux'
import { withRouteData } from 'react-static'
import Grid from '@material-ui/core/Grid'
import Typography from '@material-ui/core/Typography'
import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import Divider from '@material-ui/core/Divider'
import PaddedCard from 'components/PaddedCard'
import ConfigList from 'components/ConfigList'
import Content from 'components/Content'
import { fetchConfiguration } from 'actions'
import configsSelector from 'selectors/configs'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'

export const Configuration = useHooks(props => {
  useEffect(() => { props.fetchConfiguration() }, [])

  return (
    <Content>
      <Grid container spacing={40}>
        <Grid item xs={8}>
          <Card>
            <CardContent>
              <Typography variant='h5' color='secondary'>
                Configuration
              </Typography>
            </CardContent>

            <Divider />

            <ConfigList
              configs={props.configs}
              error={props.error}
            />
          </Card>
        </Grid>
        <Grid item xs={4}>
          <Grid container spacing={40}>
            <Grid item xs={12}>
              <PaddedCard>
                <Typography variant='h5' color='secondary'>
                  Version
                </Typography>
                <Typography variant='body1'>
                  {props.version}
                </Typography>
              </PaddedCard>
            </Grid>
            <Grid item xs={12}>
              <PaddedCard>
                <Typography variant='h5' color='secondary'>
                  SHA
                </Typography>
                <Typography variant='body1'>
                  {props.sha}
                </Typography>
              </PaddedCard>
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
  matchRouteAndMapDispatchToProps({fetchConfiguration})
)(Configuration)

export default withRouteData(ConnectedConfiguration)
