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
  textfield: {
    paddingTop: theme.spacing.unit * 1.25
  },
  form: {
    paddingTop: theme.spacing.unit * 4
  },
  card: {
    paddingBottom: theme.spacing.unit * 2
  },
  button: {
    marginTop: theme.spacing.unit * 3
  },
  flash: {
    textAlign: 'center'
  }
})

const FormLayout = ({ isSubmitting, classes, handleChange, errors, success, networkError }) => (
  <Fragment>
    {
      errors.length > 0 &&
      <Flash error className={classes.flash}>
        {errors.map((msg, i) => <p key={i}>{msg}</p>)}
      </Flash>
    }
    {
      !(errors.length > 0) && networkError &&
      <Flash error className={classes.flash}>
        Received a Network Error.
      </Flash>
    }
    {
      JSON.stringify(success) !== '{}' &&
      <Flash success className={classes.flash}>
        Bridge <Link to={`/bridges/${success.name}`}>{success.name}</Link> was successfully created.
      </Flash>
    }
    <Form className={classes.form} noValidate>
      <Grid container justify='center' spacing={0}>
        <Grid item xs={2}>
          <TextField
            fullWidth
            onChange={handleChange}
            className={classes.textfield}
            label='Type Bridge Name'
            type='name'
            name='name'
            placeholder='name'
          />
        </Grid>
      </Grid>
      <Grid container justify='center' spacing={0}>
        <Grid item xs={2}>
          <TextField
            label='Type Bridge URL'
            type='url'
            name='url'
            placeholder='url'
            fullWidth
            onChange={handleChange}
            className={classes.textfield}
          />
        </Grid>
      </Grid>
      <Grid container justify='center' spacing={0}>
        <Grid item xs={2}>
          <TextField
            onChange={handleChange}
            className={classes.textfield}
            fullWidth
            type='confirmations'
            name='confirmations'
            placeholder='confirmations'
            label='Type Confirmations'
          />
        </Grid>
      </Grid>
      <Grid container justify='center'>
        <Button color='primary' type='submit' className={classes.button} disabled={isSubmitting}>
          Build Bridge
        </Button>
      </Grid>
    </Form>
  </Fragment>
)

const BridgeForm = withFormik({
  mapPropsToValues (props) {
    const { name, url, confirmations } = props
    return {
      name: name || '',
      url: url || '',
      confirmations: confirmations || ''
    }
  },
  handleSubmit (values, { props }) {
    const formattedValues = JSON.parse(JSON.stringify(values).replace('confirmations', 'defaultConfirmations'))
    formattedValues.defaultConfirmations = parseInt(formattedValues.defaultConfirmations) || 0
    props.submitCreate('v2/bridge_types', formattedValues, true)
  }
})(FormLayout)

const mapStateToProps = state => ({
  success: state.create.successMessage,
  errors: state.create.errors.messages,
  networkError: state.create.networkError
})

export const ConnectedBridgeForm = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({submitCreate})
)(BridgeForm)

export default withStyles(styles)(ConnectedBridgeForm)
