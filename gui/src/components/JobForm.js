import React, { Fragment } from 'react'
import { withFormik, Form } from 'formik'
import { withStyles } from '@material-ui/core/styles'
import Button from '@material-ui/core/Button'
import { TextField, Grid } from '@material-ui/core'
import { connect } from 'react-redux'
import { submitJobSpec } from 'actions'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import { Prompt } from 'react-static'
import { BridgeAndJobNotifications } from './FormNotifications'

const styles = theme => ({
  card: {
    paddingBottom: theme.spacing.unit * 2
  },
  flash: {
    textAlign: 'center',
    paddingTop: theme.spacing.unit,
    paddingBottom: theme.spacing.unit
  }
})

const JobFormLayout = ({
  isSubmitting,
  classes,
  handleChange,
  error,
  success,
  authenticated,
  networkError,
  values,
  submitCount
}) => (
  <Fragment>
    <Prompt
      when={values.json !== '' && submitCount === 0}
      message='You have not submitted the form, are you sure you want to leave?'
    />
    <BridgeAndJobNotifications
      error={error}
      success={success}
      networkError={networkError}
      authenticated={authenticated}
      classes={classes}
      jobOrBridge='Job'
    />
    <Form noValidate>
      <Grid container alignItems='center' justify='center'>
        <Grid item xl={8}>
          <TextField
            onChange={handleChange}
            label='Paste JSON'
            placeholder='Paste JSON'
            multiline
            rows={10}
            fullWidth
            margin='normal'
            name='json'
            id='json'
          />
        </Grid>
        <Grid container justify='center'>
          <Button variant='contained' color='primary' type='submit' disabled={isSubmitting || !values.json}>
              Build Job
          </Button>
        </Grid>
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
  handleSubmit (values, { props, setSubmitting }) {
    props.submitJobSpec(values.json.trim(), false)
    setTimeout(() => { setSubmitting(false) }, 1000)
  }
})(JobFormLayout)

const mapStateToProps = state => ({
  success: state.create.successMessage,
  error: state.create.errors,
  networkError: state.create.networkError,
  authenticated: state.authentication.allowed
})

export const ConnectedJobForm = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({ submitJobSpec })
)(JobForm)

export default withStyles(styles)(ConnectedJobForm)
