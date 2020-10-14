import React from 'react'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import { TextField, Grid } from '@material-ui/core'
import { withFormik, FormikProps, Form as FormikForm } from 'formik'
import * as storage from '@chainlink/local-storage'
import { Prompt } from 'react-router-dom'
import isEqual from 'lodash/isEqual'
import Button from 'components/Button'

const styles = (theme: Theme) =>
  createStyles({
    card: {
      paddingBottom: theme.spacing.unit * 2,
    },
    flash: {
      textAlign: 'center',
      paddingTop: theme.spacing.unit,
      paddingBottom: theme.spacing.unit,
    },
    button: {
      marginTop: theme.spacing.unit * 2,
    },
  })

const SUBMITTING_TIMEOUT_MS = 1000
const UNSAVED_JOB_SPEC = 'persistSpec'

interface FormValues {
  json: string
}

interface OwnProps extends Partial<FormValues>, WithStyles<typeof styles> {
  definition: string
  actionText: string
  isSubmitting: boolean
  handleChange: Function
  errors: any
  onSubmit: (values: FormValues, onSuccess: Function, onError: Function) => void
  onSuccess: Function
  onError: Function
}

type Props = FormikProps<FormValues> & OwnProps

function submitSuccess(callback: Function) {
  return (response: object) => {
    storage.remove(UNSAVED_JOB_SPEC)
    return callback(response)
  }
}

function submitError(callback: Function, values: FormValues) {
  return (error: object) => {
    storage.set(UNSAVED_JOB_SPEC, values.json)
    return callback(error)
  }
}

function initialValues({ json }: OwnProps): FormValues {
  const unsavedJobSpec = storage.get(UNSAVED_JOB_SPEC)
  return unsavedJobSpec ? { json: unsavedJobSpec } : { json: json || '' }
}

function isDirty(props: Props): boolean {
  const initial = initialValues(props)
  return !isEqual(props.values, initial) && props.submitCount === 0
}

const Form: React.FC<Props> = (props) => {
  return (
    <>
      <Prompt
        when={isDirty(props)}
        message="You have not submitted the form, are you sure you want to leave?"
      />
      <FormikForm noValidate>
        <Grid container>
          <Grid item xs={12}>
            <TextField
              value={props.values.json}
              onChange={props.handleChange}
              error={props.errors.json && props.touched.json}
              helperText={props.errors.json}
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
              disabled={props.isSubmitting}
              className={props.classes.button}
              size="large"
            >
              {props.actionText}
            </Button>
          </Grid>
        </Grid>
      </FormikForm>
    </>
  )
}

const WithFormikForm = withFormik<OwnProps, FormValues>({
  mapPropsToValues({ definition }) {
    const json =
      JSON.stringify(definition, null, '\t') ||
      storage.get(UNSAVED_JOB_SPEC) ||
      ''
    return { json }
  },

  validate(values) {
    try {
      JSON.parse(values.json)
      return {}
    } catch {
      return { json: 'Invalid JSON' }
    }
  },

  handleSubmit(values, { props, setSubmitting }) {
    const definition = JSON.parse(values.json)
    props.onSubmit(
      definition,
      submitSuccess(props.onSuccess),
      submitError(props.onError, values),
    )
    setTimeout(() => {
      setSubmitting(false)
    }, SUBMITTING_TIMEOUT_MS)
  },
})(Form)

export default withStyles(styles)(WithFormikForm)
