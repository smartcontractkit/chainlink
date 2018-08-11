import React, { Fragment } from 'react'
import { withFormik, Form } from 'formik'
import { withStyles } from '@material-ui/core/styles'
import Button from '@material-ui/core/Button'
import { TextField, Grid } from '@material-ui/core'
import { connect } from 'react-redux'
import { submitCreate } from 'actions'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import { Prompt } from 'react-static'
import { BridgeAndJobNotifications } from './FormNotifications';
const styles = theme => ({
  textfield: {
    paddingTop: theme.spacing.unit * 1.25,
    width: theme.spacing.unit * 50
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
    textAlign: 'center',
    paddingTop: theme.spacing.unit,
    paddingBottom: theme.spacing.unit
  }
})

const BridgeFormLayout = ({
  isSubmitting,
  classes,
  handleChange,
  error,
  success,
  authenticated,
  networkError,
  values
}) => (
  <Fragment>
    <Prompt
      when={(values.name !== '' || values.url !== '' || values.confirmations !== '') && !isSubmitting}
      message='You have not submitted the form, are you sure you want to leave?'
    />
    <BridgeAndJobNotifications
      error={error}
      success={success}
      networkError={networkError}
      authenticated={authenticated}
      jobOrBridge='Bridge'
      classes={classes}
    />
    <Form className={classes.form} noValidate>
      <Grid container direction='column' alignItems='center'>
        <TextField
          onChange={handleChange}
          className={classes.textfield}
          label='Type Bridge Name'
          name='name'
          id='name'
          placeholder='name'
        />
        <TextField
          label='Type Bridge URL'
          name='url'
          id='url'
          placeholder='url'
          onChange={handleChange}
          className={classes.textfield}
        />
        <TextField
          onChange={handleChange}
          className={classes.textfield}
          name='confirmations'
          placeholder='confirmations'
          id='confirmations'
          label='Type Confirmations'
        />
        <Button
          color='primary'
          type='submit'
          className={classes.button}
          disabled={isSubmitting || !values.name || !values.url}>
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
})(BridgeFormLayout)

const mapStateToProps = state => ({
  success: state.create.successMessage,
  error: state.create.errors,
  networkError: state.create.networkError,
  authenticated: state.authentication.allowed
})

export const ConnectedBridgeForm = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({ submitCreate })
)(BridgeForm)

export default withStyles(styles)(ConnectedBridgeForm)
