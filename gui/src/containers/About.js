import React from 'react'
import Grid from '@material-ui/core/Grid'
import Typography from '@material-ui/core/Typography'
import PaddedCard from 'components/PaddedCard'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import { withStyles } from '@material-ui/core/styles'
import { withRouteData } from 'react-static'
import { connect } from 'react-redux'

const styles = theme => ({
  title: {
    marginTop: theme.spacing.unit * 5,
    marginBottom: theme.spacing.unit * 5
  }
})

const About = ({classes, version, sha}) => {
  return (
    <div>
      <Typography variant='display2' color='inherit' className={classes.title}>
        About
      </Typography>

      <Grid container spacing={40}>
        <Grid item xs={12}>
          <PaddedCard>
            <Typography variant='subheading' color='textSecondary'>Version</Typography>
            <Typography variant='body1' color='inherit'>
              {version}
            </Typography>
            <Typography variant='subheading' color='textSecondary'>SHA</Typography>
            <Typography variant='body1' color='inherit'>
              {sha}
            </Typography>
          </PaddedCard>
        </Grid>
      </Grid>
    </div>
  )
}

const mapStateToProps = state => ({})

export const ConnectedAbout = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({})
)(About)

export default withRouteData(withStyles(styles)(ConnectedAbout))
