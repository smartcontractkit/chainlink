import React from 'react'
import { storiesOf } from '@storybook/react'
import { Router } from 'react-static'
import {muiTheme} from 'storybook-addon-material-ui'
import { createMuiTheme } from '@material-ui/core/styles'
import Grid from '@material-ui/core/Grid'
import Typography from '@material-ui/core/Typography'
import SimpleListCard from 'components/Cards/SimpleList'
import SimpleListCardItem from 'components/Cards/SimpleListItem'
import TokenBalanceCard from 'components/Cards/TokenBalance'
import StatusCard from 'components/JobRuns/StatusCard'
import JobRunsList from 'components/JobRuns/List'
import theme from '../../gui/src/theme'

const customTheme = createMuiTheme(theme)

storiesOf('Cards', module)
  .addDecorator(muiTheme([customTheme]))
  .add('SimpleList', () => (
    <Grid container>
      <Grid xs={4}>
        <SimpleListCard title='Recently Created'>
          {['jobs', 'distribution', 'jump'].map(text => (
            <SimpleListCardItem>
              <Typography>{text}</Typography>
            </SimpleListCardItem>
          ))}
        </SimpleListCard>
      </Grid>
    </Grid>
  ))
  .add('TokenBalance', () => (
    <Grid container>
      <Grid xs={4}>
        <TokenBalanceCard title='Ether Balance' value={'10000000000000000000000'} />
      </Grid>
    </Grid>
  ))
  .add('Status', () => (
    <React.Fragment>
      <StatusCard>unstarted</StatusCard>
      <StatusCard>completed</StatusCard>
      <StatusCard>errored</StatusCard>
    </React.Fragment>
  ))

storiesOf('Tabular Data', module)
  .addDecorator(muiTheme([customTheme]))
  .add('Job Runs', () => (
    <Router>
      <JobRunsList runs={[
        {id: 'f5b5c848b8154d5eab8cd9a36fe1d506', status: 'errored', createdAt: '2018-11-26T18:26:42.133809-08:00', result: {data: {}, error: 'server not responding'}},
        {id: 'c1aeec88e8104424aa69deb383e76695', status: 'completed', createdAt: '2018-11-23T09:18:14.120683-08:00', result: {data: {price: 123.45}}}
      ]} />
    </Router>
  ))
