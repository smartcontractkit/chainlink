import React from 'react'
import Grid from '@material-ui/core/Grid'
import Typography from '@material-ui/core/Typography'
import Title from 'components/Title'
import PaddedCard from 'components/PaddedCard'
import Content from 'components/Content'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import { withRouteData } from 'react-static'
import { connect } from 'react-redux'

const About = ({version, sha}) => (
  <Content>
    <Title>About</Title>

    <Grid container spacing={40}>
      <Grid item xs={12}>
        <PaddedCard>
          <Typography variant='subtitle1' color='textSecondary'>Version</Typography>
          <Typography variant='body1' color='inherit'>
            {version}
          </Typography>
          <Typography variant='subtitle1' color='textSecondary'>SHA</Typography>
          <Typography variant='body1' color='inherit'>
            {sha}
          </Typography>
        </PaddedCard>
      </Grid>
    </Grid>
  </Content>
)

export const ConnectedAbout = connect(
  null,
  matchRouteAndMapDispatchToProps({})
)(About)

export default withRouteData(ConnectedAbout)
