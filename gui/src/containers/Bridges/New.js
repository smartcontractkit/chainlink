import React from 'react'
import PropTypes from 'prop-types'
import { withStyles } from '@material-ui/core/styles'
import Grid from '@material-ui/core/Grid'
import PaddedCard from 'components/PaddedCard'
import Typography from '@material-ui/core/Typography'
import Breadcrumb from 'components/Breadcrumb'
import BreadcrumbItem from 'components/BreadcrumbItem'
import BridgeForm from 'components/BridgeForm'

const styles = theme => ({
  title: {
    marginTop: theme.spacing.unit * 5,
    marginBottom: theme.spacing.unit * 5
  },
  breadcrumb: {
    marginTop: theme.spacing.unit * 5,
    marginBottom: theme.spacing.unit * 5
  }
})

const New = props => {
  return (
    <React.Fragment>
      <Breadcrumb className={props.classes.breadcrumb}>
        <BreadcrumbItem href='/'>Dashboard</BreadcrumbItem>
        <BreadcrumbItem>></BreadcrumbItem>
        <BreadcrumbItem href='/bridges'>Bridges</BreadcrumbItem>
        <BreadcrumbItem>></BreadcrumbItem>
        <BreadcrumbItem>New</BreadcrumbItem>
      </Breadcrumb>
      <Typography variant='display2' color='inherit' className={props.classes.title}>
        New Bridge
      </Typography>

      <Grid container spacing={40}>
        <Grid item xs={8}>
          <PaddedCard>
            <BridgeForm />
          </PaddedCard>
        </Grid>
      </Grid>
    </React.Fragment>
  )
}

New.propTypes = {
  classes: PropTypes.object.isRequired
}

export default withStyles(styles)(New)
