import React, { Fragment } from 'react'
import { withFormik, Form } from 'formik'
import { withStyles } from '@material-ui/core/styles'
import Button from '@material-ui/core/Button'
import { TextField, Grid } from '@material-ui/core'
import { connect } from 'react-redux'
import { submitCreate } from 'actions'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import Flash from './Flash'
import { Link, Prompt } from 'react-static'

const styles = theme => ({
  jsonfield: {
    paddingTop: theme.spacing.unit * 1.25,
    width: theme.spacing.unit * 150
  },
  card: {
    paddingBottom: theme.spacing.unit * 2
  },
  flash: {
    textAlign: 'center',
    paddingTop: theme.spacing.unit,
    paddingBottom: theme.spacing.unit
  }
})

const JobFormLayout = ({ isSubmitting, classes, handleChange, error, success, authenticated, networkError, values }) => (
  <Fragment>
      <Prompt when={values.json !== '' && !isSubmitting} message='You have not submitted the form, are you sure you want to leave?'/>
    {
      error.length > 0 &&
      <Flash error className={classes.flash}>
        {error.map((msg, i) => <span key={i}>{msg}</span>)}
      </Flash>
    }
    {
      !(error.length > 0) && networkError &&
        <Flash error className={classes.flash}>
          Received a Network Error.
        </Flash>
    }
    {
      !authenticated &&
      <Flash warning className={classes.flash}>
        Session expired. <Link to='/signin'>Please sign back in.</Link>
      </Flash>
    }
    {
      JSON.stringify(success) !== '{}' &&
      <Flash success className={classes.flash}>
        Job <Link to={`/job_specs/${success.id}`}>{success.id}</Link> was successfully created.
      </Flash>
    }
    <Form noValidate>
      <Grid container direction='column' alignItems='center'>
        <TextField
          onChange={handleChange}
          label='Paste JSON'
          placeholder='Paste JSON'
          multiline
          className={classes.jsonfield}
          margin='normal'
          name='json'
        />
        <Button color='primary' type='submit' disabled={isSubmitting || !values.json}>
          Build Job
        </Button>
      </Grid>
    </Form>
  </Fragment>
)

const JobForm = withFormik({
  mapPropsToValues ({ json }) {
    return {
      json: json || ''
    }
  },
  handleSubmit (values, { props }) {
    props.submitCreate('v2/specs', values.json.trim(), false)
  }
})(JobFormLayout)

const mapStateToProps = state => ({
  success: state.create.successMessage,
  error: state.create.errors,
  networkError: state.create.networkError,
  authenticated: state.authentication.allowed
})

export const ConnectedJobForm = connect(mapStateToProps, matchRouteAndMapDispatchToProps({ submitCreate }))(JobForm)

export default withStyles(styles)(ConnectedJobForm)
