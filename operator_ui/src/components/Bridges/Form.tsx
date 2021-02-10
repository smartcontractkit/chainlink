import { Grid, TextField } from '@material-ui/core'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import * as storage from 'utils/local-storage'
import { withFormik, FormikProps, Form as FormikForm } from 'formik'
import normalizeUrl from 'normalize-url'
import React from 'react'
import { Prompt } from 'react-router-dom'
import isEqual from 'lodash/isEqual'
import Button from 'components/Button'

const styles = (theme: Theme) =>
  createStyles({
    textfield: {
      paddingTop: theme.spacing.unit * 1.25,
    },
    card: {
      paddingBottom: theme.spacing.unit * 2,
    },
    button: {
      marginTop: theme.spacing.unit * 3,
    },
    flash: {
      textAlign: 'center',
      paddingTop: theme.spacing.unit,
      paddingBottom: theme.spacing.unit,
    },
  })

const SUBMITTING_TIMEOUT_MS = 1000
const UNSAVED_BRIDGE = 'persistBridge'

interface FormValues {
  name: string
  minimumContractPayment: string
  confirmations: number
  url: string
}

interface OwnProps extends Partial<FormValues>, WithStyles<typeof styles> {
  actionText: string
  nameDisabled?: boolean
  onSubmit: (values: FormValues, onSuccess: Function, onError: Function) => void
  onSuccess: Function
  onError: Function
}

type Props = FormikProps<FormValues> & OwnProps

function submitSuccess(callback: Function) {
  return (response: object) => {
    storage.remove(UNSAVED_BRIDGE)
    return callback(response)
  }
}

function submitError(callback: Function, values: FormValues) {
  return (error: object) => {
    storage.setJson(UNSAVED_BRIDGE, values)
    return callback(error)
  }
}

const DEFAULT_VALUES: FormValues = {
  name: '',
  url: '',
  minimumContractPayment: '0',
  confirmations: 0,
}

function getValue(ownProps: OwnProps, key: keyof FormValues) {
  const unsavedBridge = storage.getJson(UNSAVED_BRIDGE)
  return ownProps[key] || unsavedBridge[key] || DEFAULT_VALUES[key]
}

function initialValues(props: OwnProps): FormValues {
  return {
    name: getValue(props, 'name'),
    url: getValue(props, 'url'),
    minimumContractPayment: getValue(props, 'minimumContractPayment'),
    confirmations: getValue(props, 'confirmations'),
  }
}

function isDirty(props: Props): boolean {
  const initial = initialValues(props)
  return !isEqual(props.values, initial) && props.submitCount === 0
}

const Form: React.SFC<Props> = (props) => (
  <>
    <Prompt
      when={isDirty(props)}
      message="You have not submitted the form, are you sure you want to leave?"
    />
    <FormikForm noValidate>
      <Grid container spacing={8}>
        <Grid item xs={12} md={7}>
          <TextField
            label="Bridge Name"
            name="name"
            placeholder="name"
            value={props.values.name}
            disabled={props.nameDisabled}
            onChange={props.handleChange}
            className={props.classes.textfield}
            fullWidth
          />
        </Grid>
        <Grid item xs={12} md={7}>
          <TextField
            label="Bridge URL"
            name="url"
            placeholder="https://"
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
                label="Minimum Contract Payment"
                name="minimumContractPayment"
                placeholder="0"
                value={props.values.minimumContractPayment}
                inputProps={{ min: 0 }}
                onChange={props.handleChange}
                className={props.classes.textfield}
                fullWidth
              />
            </Grid>
            <Grid item xs={7}>
              <TextField
                label="Confirmations"
                name="confirmations"
                placeholder="0"
                value={props.values.confirmations}
                type="number"
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
            variant="primary"
            type="submit"
            className={props.classes.button}
            disabled={props.isSubmitting}
            size="large"
          >
            {props.actionText}
          </Button>
        </Grid>
      </Grid>
    </FormikForm>
  </>
)

Form.defaultProps = {
  nameDisabled: false,
}

const WithFormikForm = withFormik<OwnProps, FormValues>({
  mapPropsToValues(ownProps) {
    return initialValues(ownProps)
  },

  handleSubmit(values, { props, setSubmitting }) {
    try {
      values.url = normalizeUrl(values.url)
    } catch {
      values.url = ''
    }
    props.onSubmit(
      values,
      submitSuccess(props.onSuccess),
      submitError(props.onError, values),
    )
    setTimeout(() => {
      setSubmitting(false)
    }, SUBMITTING_TIMEOUT_MS)
  },
})(Form)

export default withStyles(styles)(WithFormikForm)
