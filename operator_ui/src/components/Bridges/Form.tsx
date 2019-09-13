import { Grid, TextField } from '@material-ui/core'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import Button from 'components/Button'
import { withFormik, FormikProps, Form as FormikForm } from 'formik'
import normalizeUrl from 'normalize-url'
import React from 'react'
import { Prompt } from 'react-router-dom'
import { get, set } from 'utils/storage'

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

const isDirty = ({ values, submitCount }: Props) => {
  return (
    (values.name !== '' ||
      values.url !== '' ||
      (values.minimumContractPayment !== '0' && values.confirmations !== 0)) &&
    submitCount === 0
  )
}

// CHECKME
interface OwnProps extends Partial<FormValues>, WithStyles<typeof styles> {
  actionText: string
  nameDisabled?: boolean
  onSubmit: any
  onSuccess: any
  onError: any
}

// CHECKME
interface FormValues {
  name: string
  minimumContractPayment: string
  confirmations: number
  url: string
}

type Props = FormikProps<FormValues> & OwnProps

const Form: React.SFC<Props> = props => (
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
  mapPropsToValues({ name, url, minimumContractPayment, confirmations }) {
    const shouldPersist = Object.keys(get('persistBridge')).length !== 0
    const persistedJSON = shouldPersist && get('persistBridge')
    if (shouldPersist) set('persistBridge', {})
    const json: FormValues = {
      name: name || '',
      url: url || '',
      minimumContractPayment: minimumContractPayment || '0',
      confirmations: confirmations || 0,
    }
    return (shouldPersist && persistedJSON) || json
  },

  handleSubmit(values, { props, setSubmitting }) {
    try {
      values.url = normalizeUrl(values.url)
    } catch {
      values.url = ''
    }
    props.onSubmit(values, props.onSuccess, props.onError)
    set('persistBridge', values)
    setTimeout(() => {
      setSubmitting(false)
    }, 1000)
  },
})(Form)

export default withStyles(styles)(WithFormikForm)
