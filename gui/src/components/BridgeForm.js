import axios from 'axios'
import React, { Fragment } from 'react'
import { withFormik, Form } from 'formik'
import * as Yup from 'yup'
import { withStyles } from '@material-ui/core/styles'
import Button from '@material-ui/core/Button'
import { TextField, Typography } from '@material-ui/core'
import { ToastContainer, toast } from 'react-toastify'
import 'react-toastify/dist/ReactToastify.css'

const styles = theme => ({
  textfield: {
    paddingTop: theme.spacing.unit * 1.25,
    width: '270px'
  },
  card: {
    paddingBottom: theme.spacing.unit * 2
  },
  form: {
    position: 'relative',
    textAlign: 'center'
  }
})

const App = ({ errors, touched, isSubmitting, classes, handleChange }) => (
  <Fragment>
    <br />
    <Form className={classes.form} noValidate>
      <div>
        {touched.name && errors.name && <Typography color='error'>{errors.name}</Typography>}
        <TextField
          onChange={handleChange}
          className={classes.textfield}
          label='Type Bridge Name'
          type='name'
          name='name'
          placeholder='name'
        />
      </div>
      <div>
        {touched.url && errors.url && <Typography color='error'>{errors.url}</Typography>}
        <TextField
          onChange={handleChange}
          className={classes.textfield}
          label='Type Bridge URL'
          type='url'
          name='url'
          placeholder='url'
        />
      </div>
      <div>
        {touched.confirmations && errors.confirmations && <Typography color='error'>{errors.confirmations}</Typography>}
        <TextField
          onChange={handleChange}
          className={classes.textfield}
          label='Type Confirmations'
          type='confirmations'
          name='confirmations'
          placeholder='confirmations'
        />
      </div>
      <Button color='primary' type='submit' disabled={isSubmitting}>
        Build Bridge
      </Button>
      <ToastContainer />
    </Form>
  </Fragment>
)

const BridgeForm = withFormik({
  mapPropsToValues ({ name, url, confirmations }) {
    return {
      name: name || '',
      url: url || '',
      confirmations: confirmations || ''
    }
  },
  validationSchema: Yup.object().shape({
    name: Yup.string().required('Name is required'),
    url: Yup.string()
      .required('URL is required'),
    confirmations: Yup.number()
      .positive('Should be a positive number')
      .typeError('Should be a number')
  }),
  handleSubmit (values) {
    const formattedValues = JSON.parse(JSON.stringify(values).replace('confirmations', 'defaultConfirmations'))
    formattedValues.defaultConfirmations = parseInt(formattedValues.defaultConfirmations) || 0
    axios
      .post('/v2/bridge_types', formattedValues, {
        headers: {
          'Content-Type': 'application/json'
        },
        auth: {
          username: 'chainlink',
          password: 'twochains'
        }
      })
      .then(res =>
        toast.success(`Bridge ${res.data.name} created`, {
          position: toast.POSITION.BOTTOM_RIGHT
        })
      )
  }
})(App)

export default withStyles(styles)(BridgeForm)
