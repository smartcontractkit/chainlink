import React, { Fragment } from 'react'
import { withFormik, Form } from 'formik'
import { withStyles } from '@material-ui/core/styles'
import Button from '@material-ui/core/Button'
import { TextField, Grid } from '@material-ui/core'
import { connect } from 'react-redux'
import { submitCreate } from 'actions'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import Flash from './Flash'
import { Link } from 'react-static'

const styles = theme => ({
  jsonfield: {
    paddingTop: theme.spacing.unit * 1.25
  },
  card: {
    paddingBottom: theme.spacing.unit * 2
  },
  flash: {
    textAlign: 'center'
  }
})

const FormLayout = ({ isSubmitting, classes, handleChange, errors, success, networkError }) => (
  <Fragment>
    {errors.length > 0 &&
      <Flash error className={classes.flash}>
        {errors.map((msg, i) => <p key={i}>{msg}</p>)}
      </Flash>
    }
    {!(errors.length > 0) && networkError &&
      <Flash error className={classes.flash}>
        Received a Network Error.
      </Flash>}
    {JSON.stringify(success) !== '{}' && (
      <Flash success className={classes.flash}>
        Job <Link to={`/job_specs/${success.id}`}>{success.id}</Link> was successfully created.
      </Flash>
    )}
    <Grid justify='center' container spacing={24}>
      <Grid item xs={5}>
        <Form noValidate>
          <TextField
            fullWidth
            onChange={handleChange}
            label='Paste JSON'
            placeholder='Paste JSON'
            multiline
            className={classes.jsonfield}
            margin='normal'
            type='json'
            name='json'
          />
          <Grid container justify='center'>
            <Grid item>
              <Button color='primary' type='submit' disabled={isSubmitting}>
                Build Job
              </Button>
            </Grid>
          </Grid>
        </Form>
      </Grid>
    </Grid>
  </Fragment>
)

const JobForm = withFormik({
  mapPropsToValues ({ json }) {
    return {
      json: json || ''
    }
  },
  handleSubmit (values, { props }) {
    props.submitCreate('v2/specs', values.json, false)
  }
})(FormLayout)

const mapStateToProps = state => ({
  success: state.create.successMessage,
  errors: state.create.errors.messages,
  networkError: state.create.networkError
})

export const ConnectedJobForm = connect(mapStateToProps, matchRouteAndMapDispatchToProps({ submitCreate }))(JobForm)

export default withStyles(styles)(ConnectedJobForm)
