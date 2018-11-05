import React from 'react'
import PropTypes from 'prop-types'
import * as formik from 'formik'
import { withStyles } from '@material-ui/core/styles'
import Button from '@material-ui/core/Button'
import { TextField, Grid } from '@material-ui/core'
import { Prompt } from 'react-static'

const styles = theme => ({
  card: {
    paddingBottom: theme.spacing.unit * 2
  },
  flash: {
    textAlign: 'center',
    paddingTop: theme.spacing.unit,
    paddingBottom: theme.spacing.unit
  },
  button: {
    marginTop: theme.spacing.unit * 2
  }
})

const Form = ({
  actionText,
  dirty,
  isSubmitting,
  classes,
  handleChange,
  values,
  touched,
  errors,
  submitCount
}) => {
  return (
    <React.Fragment>
      <Prompt
        when={values.json !== '' && submitCount === 0}
        message='You have not submitted the form, are you sure you want to leave?'
      />
      <formik.Form noValidate>
        <Grid container justify='center'>
          <Grid container justify='center'>
            <Grid item sm={8}>
              <TextField
                value={values.json}
                onChange={handleChange}
                error={errors.json && touched.json}
                fullWidth
                label='Paste JSON'
                rows={10}
                rowsMax={25}
                placeholder='Paste JSON'
                multiline margin='normal'
                name='json'
                id='json'
              />
            </Grid>
          </Grid>
          <Button
            className={classes.button}
            variant='contained'
            color='primary'
            type='submit'
            disabled={isSubmitting || !values.json}
          >
            {actionText}
          </Button>
        </Grid>
      </formik.Form>
    </React.Fragment>
  )
}

Form.propTypes = {
  actionText: PropTypes.string.isRequired,
  onSubmit: PropTypes.func.isRequired
}

const formikOpts = {
  mapPropsToValues ({definition}) {
    const json = JSON.stringify(definition, null, '\t') || ''
    return {json}
  },

  validate (values) {
    const errors = {}

    try {
      JSON.parse(values.json, null, '\t')
    } catch (e) {
      errors.json = 'Invalid JSON'
    }

    return errors
  },

  handleSubmit (values, { props, setSubmitting }) {
    const definition = JSON.parse(values.json)
    props.onSubmit(definition, props.onSuccess, props.onError)
    setTimeout(() => { setSubmitting(false) }, 1000)
  }
}

const FormikForm = formik.withFormik(formikOpts)(Form)

export default withStyles(styles)(FormikForm)
