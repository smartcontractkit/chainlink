import React from 'react'
import { storiesOf } from '@storybook/react'
import { muiTheme } from 'storybook-addon-material-ui'
import { createMuiTheme } from '@material-ui/core/styles'
import Grid from '@material-ui/core/Grid'
import Typography from '@material-ui/core/Typography'
import CardTitle from '../components/Cards/Title'
import SimpleList from '../components/Cards/SimpleList'
import SimpleListItem from '../components/Cards/SimpleListItem'
import KeyValueList from '../components/KeyValueList'
import theme from '../theme'

const customTheme = createMuiTheme(theme)

storiesOf('Cards', module)
  .addDecorator(muiTheme([customTheme]))
  .add('CardTitle', () => (
    <React.Fragment>
      <CardTitle>Without Divider</CardTitle>
      <CardTitle divider>With Divider</CardTitle>
    </React.Fragment>
  ))
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
