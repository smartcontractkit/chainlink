import React from 'react'
import PropTypes from 'prop-types'
import { Prompt } from 'react-router-dom'
import * as formik from 'formik'
import { withStyles } from '@material-ui/core/styles'
import { TextField, Grid } from '@material-ui/core'
import Button from '@material-ui/core/Button'

const styles = theme => ({
  textfield: {
    paddingTop: theme.spacing.unit * 1.25
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

const isDirty = ({ values, name, url, minimumContractPayment, confirmations, submitCount }) => {
  return (
    values.name !== name ||
    values.url !== url ||
    values.minimumContractPayment.toString() !== minimumContractPayment.toString() ||
    values.confirmations !== confirmations
  ) && submitCount === 0
}

const Form = props => (
  <React.Fragment>
    <Prompt
      when={isDirty(props)}
      message='You have not submitted the form, are you sure you want to leave?'
    />
    <formik.Form noValidate>
      <Grid container spacing={8}>
        <Grid item xs={12} md={7}>
          <TextField
            label='Bridge Name'
            name='name'
            placeholder='name'
            value={props.values.name}
            disabled={props.nameDisabled}
            onChange={props.handleChange}
            className={props.classes.textfield}
            fullWidth
          />
        </Grid>
        <Grid item xs={12} md={7}>
          <TextField
            label='Bridge URL'
            name='url'
            placeholder='https://'
            value={props.values.url}
            onChange={props.handleChange}
            className={props.classes.textfield}
            fullWidth
          />
        </Grid>
        <Grid item xs={12} md={7}>
          <Grid container spacing={8}>
            <Grid item xs={7}>
              <TextField
                label='Minimum Contract Payment'
                name='minimumContractPayment'
                placeholder='0'
                value={props.values.minimumContractPayment}
                type='number'
                inputProps={{ min: 0 }}
                onChange={props.handleChange}
                className={props.classes.textfield}
                fullWidth
              />
            </Grid>
            <Grid item xs={7}>
              <TextField
                label='Confirmations'
                name='confirmations'
                placeholder='0'
                value={props.values.confirmations}
                type='number'
                inputProps={{ min: 0 }}
                onChange={props.handleChange}
                className={props.classes.textfield}
                fullWidth
              />
            </Grid>
          </Grid>
        </Grid>
        <Grid item xs={12} md={7}>
          <Button
            variant='contained'
            color='primary'
            type='submit'
            className={props.classes.button}
            disabled={props.isSubmitting || !props.values.name || !props.values.url}
          >
            {props.actionText}
          </Button>
        </Grid>
      </Grid>
    </formik.Form>
  </React.Fragment>
)

Form.defaultPropTypes = {
  nameDisabled: false
}

Form.propTypes = {
  actionText: PropTypes.string.isRequired,
  onSubmit: PropTypes.func.isRequired,
  name: PropTypes.string,
  nameDisabled: PropTypes.bool,
  url: PropTypes.string,
  minimumContractPayment: PropTypes.string,
  confirmations: PropTypes.number,
  onSuccess: PropTypes.func.isRequired,
  onError: PropTypes.func.isRequired
}

const formikOpts = {
  mapPropsToValues ({ name, url, minimumContractPayment, confirmations }) {
    return {
      name: name || '',
      url: url || '',
      minimumContractPayment: minimumContractPayment || 0,
      confirmations: confirmations || 0
    }
  },

  handleSubmit (values, { props, setSubmitting }) {
    props.onSubmit(values, props.onSuccess, props.onError)
    setTimeout(() => { setSubmitting(false) }, 1000)
  }
}

const FormikForm = formik.withFormik(formikOpts)(Form)

export default withStyles(styles)(FormikForm)
