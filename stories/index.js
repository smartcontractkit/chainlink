import React from 'react'
import { storiesOf } from '@storybook/react'
import { Router } from 'react-static'
import { createMuiTheme } from '@material-ui/core/styles'
import Button from '@material-ui/core/Button'
import {muiTheme} from 'storybook-addon-material-ui'
import Grid from '@material-ui/core/Grid'
import Typography from '@material-ui/core/Typography'
import JavascriptTimeAgo from 'javascript-time-ago'
import en from 'javascript-time-ago/locale/en'
import theme from '../gui/src/theme'
import TokenBalanceCard from 'components/Cards/TokenBalance'
import TimeAgo from 'components/TimeAgo'
import SimpleListCard from 'components/Cards/SimpleList'
import SimpleListCardItem from 'components/Cards/SimpleListItem'
import Logo from 'components/Logo'
import JobRunsList from 'components/JobRuns/List'

window.JavascriptTimeAgo = JavascriptTimeAgo
JavascriptTimeAgo.locale(en)

const customTheme = createMuiTheme(theme)

storiesOf('Typography', module)
  .addDecorator(muiTheme([customTheme]))
  .add('V2 Variants', () => (
    <React.Fragment>
      <Typography component='h2' variant='h1' gutterBottom>
        h1. Heading
      </Typography>
      <Typography variant='h2' gutterBottom>
        h2. Heading
      </Typography>
      <Typography variant='h3' gutterBottom>
        h3. Heading
      </Typography>
      <Typography variant='h4' gutterBottom>
        h4. Heading
      </Typography>
      <Typography variant='h5' gutterBottom>
        h5. Heading
      </Typography>
      <Typography variant='h6' gutterBottom>
        h6. Heading
      </Typography>
      <Typography variant='subtitle1' gutterBottom>
        subtitle1. Lorem ipsum dolor sit amet, consectetur adipisicing elit. Quos blanditiis tenetur
      </Typography>
      <Typography variant='subtitle2' gutterBottom>
        subtitle2. Lorem ipsum dolor sit amet, consectetur adipisicing elit. Quos blanditiis tenetur
      </Typography>
      <Typography variant='body1' gutterBottom>
        body1. Lorem ipsum dolor sit amet, consectetur adipisicing elit. Quos blanditiis tenetur
        unde suscipit, quam beatae rerum inventore consectetur, neque doloribus, cupiditate numquam
        dignissimos laborum fugiat deleniti? Eum quasi quidem quibusdam.
      </Typography>
      <Typography variant='body2' gutterBottom>
        body2. Lorem ipsum dolor sit amet, consectetur adipisicing elit. Quos blanditiis tenetur
        unde suscipit, quam beatae rerum inventore consectetur, neque doloribus, cupiditate numquam
        dignissimos laborum fugiat deleniti? Eum quasi quidem quibusdam.
      </Typography>
      <Typography variant='button' gutterBottom>
        button text
      </Typography>
      <Typography variant='caption' gutterBottom>
        caption text
      </Typography>
      <Typography variant='overline' gutterBottom>
        overline text
      </Typography>
    </React.Fragment>
  ))

storiesOf('Buttons', module)
  .addDecorator(muiTheme([customTheme]))
  .add('Contained', () => (
    <React.Fragment>
      <Button variant='contained'>Default Button</Button>
      <Button variant='contained' color='primary'>Primary Button</Button>
      <Button variant='contained' color='secondary'>Secondary</Button>
      <Button variant='contained' color='secondary' disabled>Disabled</Button>
      <Button variant='contained' href='#contained-buttons'>Link</Button>
      <input
        accept='image/*'
        id='contained-button-file'
        multiple
        type='file'
      />
      <label htmlFor='contained-button-file'>
        <Button variant='contained' component='span'>
          Upload
        </Button>
      </label>
    </React.Fragment>
  ))
  .add('Outlined', () => (
    <React.Fragment>
      <Button variant='outlined'>Default Button</Button>
      <Button variant='outlined' color='primary'>Primary Button</Button>
      <Button variant='outlined' color='secondary'>Secondary</Button>
      <Button variant='outlined' color='secondary' disabled>Disabled</Button>
      <Button variant='outlined' href='#outlined-buttons'>Link</Button>
      <input
        accept='image/*'
        id='outlined-button-file'
        multiple
        type='file'
      />
      <label htmlFor='outlined-button-file'>
        <Button variant='outlined' component='span'>
          Upload
        </Button>
      </label>
    </React.Fragment>
  ))

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

storiesOf('Custom', module)
  .addDecorator(muiTheme([customTheme]))
  .add('Logo', () => (
    <Logo width={40} height={50} />
  ))
  .add('TimeAgo', () => (
    <Typography>
      <TimeAgo>2018-11-27T02:26:42.014852Z</TimeAgo>
    </Typography>
  ))
