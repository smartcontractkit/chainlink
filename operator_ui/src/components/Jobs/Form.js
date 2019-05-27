import React from 'react'
import PropTypes from 'prop-types'
import * as formik from 'formik'
import { withStyles } from '@material-ui/core/styles'
import Button from 'components/Button'
import { TextField, Grid } from '@material-ui/core'
import { Prompt } from 'react-router-dom'
import { set, get } from 'utils/storage'

const styles = theme => ({
  card: {
    paddingBottom: theme.spacing(2)
  },
  flash: {
    textAlign: 'center',
    paddingTop: theme.spacing(1),
    paddingBottom: theme.spacing(1)
  },
  button: {
    marginTop: theme.spacing(2)
  }
})

const Form = ({
  actionText,
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
        message="You have not submitted the form, are you sure you want to leave?"
      />
      <formik.Form noValidate>
        <Grid container>
          <Grid item xs={12}>
            <TextField
              value={values.json}
              onChange={handleChange}
              error={errors.json && touched.json}
              helperText={errors.json}
              autoComplete="off"
              fullWidth
              label="JSON Blob"
              rows={10}
              rowsMax={25}
              placeholder="Paste JSON"
              multiline
              margin="normal"
              name="json"
              id="json"
              variant="outlined"
            />
          </Grid>
          <Grid item xs={12}>
            <Button
              variant="primary"
              type="submit"
              disabled={isSubmitting}
              className={classes.button}
            >
              {actionText}
            </Button>
          </Grid>
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
  mapPropsToValues({ definition }) {
    const shouldPersist = Object.keys(get('persistSpec')).length !== 0
    let persistedJSON = shouldPersist && get('persistSpec')
    if (shouldPersist) set('persistSpec', {})
    const json =
      JSON.stringify(definition, null, '\t') ||
      (shouldPersist && persistedJSON) ||
      ''
    return { json }
  },

  validate(values) {
    const errors = {}

    try {
      JSON.parse(values.json, null, '\t')
    } catch (e) {
      errors.json = 'Invalid JSON'
    }

    return errors
  },

  handleSubmit(values, { props, setSubmitting }) {
    const definition = JSON.parse(values.json)
    set('persistSpec', values.json)
    props.onSubmit(definition, props.onSuccess, props.onError)
    setTimeout(() => {
      setSubmitting(false)
    }, 1000)
  }
}

const FormikForm = formik.withFormik(formikOpts)(Form)

export default withStyles(styles)(FormikForm)
