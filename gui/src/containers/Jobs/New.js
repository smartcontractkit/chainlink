import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'react-redux'
import { Link } from 'react-static'
import { withStyles } from '@material-ui/core/styles'
import Grid from '@material-ui/core/Grid'
import Breadcrumb from 'components/Breadcrumb'
import BreadcrumbItem from 'components/BreadcrumbItem'
import Title from 'components/Title'
import PaddedCard from 'components/PaddedCard'
import Form from 'components/Jobs/Form'
import ErrorMessage from 'components/Errors/Message'
import { createJobSpec } from 'actions'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'

const styles = theme => ({
  breadcrumb: {
    marginTop: theme.spacing.unit * 5,
    marginBottom: theme.spacing.unit * 5
  }
})

const SuccessNotification = ({data}) => (
  <React.Fragment>
    Successfully created job <Link to={`/jobs/${data.id}`}>{data.id}</Link>
  </React.Fragment>
)

const New = props => (
  <React.Fragment>
    <Breadcrumb className={props.classes.breadcrumb}>
      <BreadcrumbItem href='/'>Dashboard</BreadcrumbItem>
      <BreadcrumbItem>></BreadcrumbItem>
      <BreadcrumbItem href='/jobs'>Jobs</BreadcrumbItem>
      <BreadcrumbItem>></BreadcrumbItem>
      <BreadcrumbItem>New</BreadcrumbItem>
    </Breadcrumb>
    <Title>New Job</Title>

    <Grid container spacing={40}>
      <Grid item xs={12}>
        <PaddedCard>
          <Form
            actionText='Create Job'
            onSubmit={props.createJobSpec}
            onSuccess={SuccessNotification}
            onError={ErrorMessage}
            {...(props.location && props.location.state)}
          />
        </PaddedCard>
      </Grid>
    </Grid>
  </React.Fragment>
)

New.propTypes = {
  classes: PropTypes.object.isRequired
}

export const ConnectedNew = connect(
  null,
  matchRouteAndMapDispatchToProps({createJobSpec})
)(New)

export default withStyles(styles)(ConnectedNew)
