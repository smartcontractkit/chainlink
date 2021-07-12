import React, { useEffect } from 'react'
import { PaddedCard } from 'components/PaddedCard'
import { KeyValueList } from 'components/KeyValueList'
import Grid from '@material-ui/core/Grid'
import Typography from '@material-ui/core/Typography'
import { fetchConfiguration } from 'actionCreators'
import Content from 'components/Content'
import DeleteJobRuns from 'pages/Configuration/DeleteJobRuns'
import { useDispatch, useSelector } from 'react-redux'
import configurationSelector from 'selectors/configuration'
import extractBuildInfo from 'utils/extractBuildInfo'
import { LoggingCard } from './LoggingCard'

const buildInfo = extractBuildInfo()

export const Configuration = () => {
  const dispatch = useDispatch()
  const config = useSelector(configurationSelector)

  useEffect(() => {
    document.title = 'Configuration'
  })

  useEffect(() => {
    dispatch(fetchConfiguration())
  }, [dispatch])

  return (
    <Content>
      <Grid container>
        <Grid item sm={12} md={8}>
          <KeyValueList title="Configuration" entries={config} showHead />
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
            <Grid item xs={12}>
              <LoggingCard />
            </Grid>
          </Grid>
        </Grid>
      </Grid>
    </Content>
  )
}

export default Configuration
