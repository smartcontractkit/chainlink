import Grid from '@material-ui/core/Grid'
import { createMuiTheme } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
import { storiesOf } from '@storybook/react'
import React from 'react'
import { muiTheme } from 'storybook-addon-material-ui'
import {
  CardTitle,
  KeyValueList,
  SimpleListCard,
  SimpleListCardItem,
  theme
} from '../src'

const customTheme = createMuiTheme(theme)

storiesOf('Cards', module)
  .addDecorator(muiTheme([customTheme]))
  .add('CardTitle', () => (
    <React.Fragment>
      <CardTitle>Without Divider</CardTitle>
      <CardTitle divider>With Divider</CardTitle>
    </React.Fragment>
  ))
  .add('SimpleListCard', () => (
    <Grid container>
      <Grid xs={4}>
        <SimpleListCard title="Recently Created">
          {['jobs', 'distribution', 'jump'].map(text => (
            <SimpleListCardItem>
              <Typography>{text}</Typography>
            </SimpleListCardItem>
          ))}
        </SimpleListCard>
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
