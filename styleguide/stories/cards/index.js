import React from 'react'
import { storiesOf } from '@storybook/react'
import { Router } from 'react-static'
import { muiTheme } from 'storybook-addon-material-ui'
import { createMuiTheme } from '@material-ui/core/styles'
import Grid from '@material-ui/core/Grid'
import Typography from '@material-ui/core/Typography'
import SimpleList from 'components/Cards/SimpleList'
import SimpleListItem from 'components/Cards/SimpleListItem'
import KeyValueList from 'components/KeyValueList'
import TokenBalance from 'components/Cards/TokenBalance'
import StatusCard from 'components/JobRuns/StatusCard'
import JobRunsList from 'components/JobRuns/List'
import CardTitle from 'components/Cards/Title'
import theme from '../../../operator_ui/src/theme'

const customTheme = createMuiTheme(theme)

storiesOf('Cards', module)
  .addDecorator(muiTheme([customTheme]))
  .add('SimpleList', () => (
    <Grid container>
      <Grid xs={4}>
        <SimpleList title="Recently Created">
          {['jobs', 'distribution', 'jump'].map(text => (
            <SimpleListItem>
              <Typography>{text}</Typography>
            </SimpleListItem>
          ))}
        </SimpleList>
      </Grid>
    </Grid>
  ))
  .add('KeyValueList', () => (
    <Grid container spacing={40}>
      <Grid item xs={12}>
        <Grid container>
          <Grid item xs={4}>
            <KeyValueList title="Loading" entries={[]} />
          </Grid>
        </Grid>
      </Grid>
      <Grid item xs={12}>
        <Grid container>
          <Grid item xs={4}>
            <KeyValueList entries={[['WITHOUT_TITLE', 'true']]} />
          </Grid>
        </Grid>
      </Grid>
      <Grid item xs={12}>
        <Grid container>
          <Grid item xs={4}>
            <KeyValueList
              title="With Title"
              entries={[['WITHOUT_TITLE', 'false']]}
            />
          </Grid>
        </Grid>
      </Grid>
    </Grid>
  ))
  .add('TokenBalance', () => (
    <Grid container>
      <Grid xs={4}>
        <TokenBalance title="Ether Balance" value={'10000000000000000000000'} />
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
  .add('CardTitle', () => (
    <React.Fragment>
      <CardTitle>Without Divider</CardTitle>
      <CardTitle divider>With Divider</CardTitle>
    </React.Fragment>
  ))

storiesOf('Tabular Data', module)
  .addDecorator(muiTheme([customTheme]))
  .add('Job Runs', () => (
    <Router>
      <JobRunsList
        runs={[
          {
            id: 'f5b5c848b8154d5eab8cd9a36fe1d506',
            status: 'errored',
            createdAt: '2018-11-26T18:26:42.133809-08:00',
            result: { data: {}, error: 'server not responding' }
          },
          {
            id: 'c1aeec88e8104424aa69deb383e76695',
            status: 'completed',
            createdAt: '2018-11-23T09:18:14.120683-08:00',
            result: { data: { price: 123.45 } }
          }
        ]}
      />
    </Router>
  ))
