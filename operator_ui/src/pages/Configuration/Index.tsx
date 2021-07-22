import React, { useEffect } from 'react'
import { KeyValueList } from 'components/KeyValueList'
import Grid from '@material-ui/core/Grid'
import { fetchConfiguration } from 'actionCreators'
import Content from 'components/Content'
import { useDispatch, useSelector } from 'react-redux'
import configurationSelector from 'selectors/configuration'
import extractBuildInfo from 'utils/extractBuildInfo'

import { LoggingCard } from './LoggingCard'
import { NodeInformation } from './NodeInformation'
import { JobRuns } from './JobRuns'

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
              <NodeInformation
                version={buildInfo.version}
                sha={buildInfo.sha}
              />
            </Grid>
            <Grid item xs={12}>
              <JobRuns />
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
