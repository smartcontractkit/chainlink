import React, { Fragment } from 'react'
import { withFormik, Form } from 'formik'
import { object, string, number } from 'yup'
import { withStyles } from '@material-ui/core/styles'
import Button from '@material-ui/core/Button'
import { TextField, Typography, Grid } from '@material-ui/core'
import { postBridge } from 'api'

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
  }
})

const FormLayout = ({ errors, touched, isSubmitting, classes, handleChange }) => (
  <Fragment>
    <Form className={classes.form} noValidate>
      <Grid container justify='center' spacing={0}>
        <Grid item xs={2}>
          {touched.name && errors.name && <Typography color='error'>{errors.name}</Typography>}
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
          {touched.url && errors.url && <Typography color='error'>{errors.url}</Typography>}
          <TextField
            fullWidth
            onChange={handleChange}
            className={classes.textfield}
            label='Type Bridge URL'
            type='url'
            name='url'
            placeholder='url'
          />
        </Grid>
      </Grid>
      <Grid container justify='center' spacing={0}>
        <Grid item xs={2}>
          {touched.confirmations && errors.confirmations && <Typography color='error'>{errors.confirmations}</Typography>}
          <TextField
            fullWidth
            onChange={handleChange}
            className={classes.textfield}
            label='Type Confirmations'
            type='confirmations'
            name='confirmations'
            placeholder='confirmations'
          />
        </Grid>
      </Grid>
      <Grid container justify='center' spacing={0}>
        <Button color='primary' type='submit' className={classes.button} disabled={isSubmitting}>
          Build Bridge
        </Button>
      </Grid>
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
  validationSchema: object().shape({
    name: string().required('Name is required'),
    url: string()
      .required('URL is required'),
    confirmations: number()
      .positive('Should be a positive number')
      .typeError('Should be a number')
  }),
  handleSubmit (values) {
    const formattedValues = JSON.parse(JSON.stringify(values).replace('confirmations', 'defaultConfirmations'))
    formattedValues.defaultConfirmations = parseInt(formattedValues.defaultConfirmations) || 0
    postBridge(formattedValues)
  }
})(FormLayout)

export default withStyles(styles)(BridgeForm)
