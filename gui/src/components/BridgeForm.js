import React, { Fragment } from 'react'
import { withFormik, Form } from 'formik'
import { withStyles } from '@material-ui/core/styles'
import Button from '@material-ui/core/Button'
import { TextField, Grid } from '@material-ui/core'
import { connect } from 'react-redux'
import { submitBridgeType } from 'actions'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import { Prompt } from 'react-static'

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
    textAlign: 'center',
    paddingTop: theme.spacing.unit,
    paddingBottom: theme.spacing.unit
  }
})

const BridgeFormLayout = ({
  isSubmitting,
  classes,
  handleChange,
  values,
  submitCount
}) => (
  <Fragment>
    <Prompt
      when={(values.name !== '' || values.url !== '' || values.confirmations !== '') && submitCount === 0}
      message='You have not submitted the form, are you sure you want to leave?'
    />
    <Form className={classes.form} noValidate>
      <Grid container justify='center'>
        <Grid item sm={2}>
          <TextField
            className={classes.textfield}
            onChange={handleChange}
            fullWidth
            label='Bridge Name'
            name='name'
            id='name'
            placeholder='name'
          />
        </Grid>
        <Grid container justify='center'>
          <Grid item sm={2}>
            <TextField
              className={classes.textfield}
              label='Bridge URL'
              name='url'
              fullWidth
              placeholder='url'
              id='url'
              onChange={handleChange}
            />
          </Grid>
        </Grid>
        <Grid container justify='center'>
          <Grid item sm={2}>
            <TextField
              className={classes.textfield}
              onChange={handleChange}
              name='minimumContractPayment'
              type='number'
              inputProps={{min: 0}}
              fullWidth
              placeholder='0'
              id='minimumContractPayment'
              label='Minimum Contract Payment'
            />
          </Grid>
        </Grid>
        <Grid container justify='center'>
          <Grid item sm={2}>
            <TextField
              className={classes.textfield}
              onChange={handleChange}
              name='confirmations'
              type='number'
              inputProps={{min: 0}}
              fullWidth
              placeholder='0'
              id='confirmations'
              label='Confirmations'
            />
          </Grid>
        </Grid>
        <Button variant='contained' color='primary' type='submit' className={classes.button} disabled={isSubmitting || !values.name || !values.url}>
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
  handleSubmit (values, { props, setSubmitting }) {
    values.confirmations = parseInt(values.confirmations) || 0
    props.submitBridgeType(values, true)
    setTimeout(() => {
      setSubmitting(false)
    }, 1000)
  }
})(BridgeFormLayout)

const mapStateToProps = state => ({
  networkError: state.create.networkError,
  fetching: state.fetching.count
})

export const ConnectedBridgeForm = connect(mapStateToProps, matchRouteAndMapDispatchToProps({ submitBridgeType }))(
  BridgeForm
)

export default withStyles(styles)(ConnectedBridgeForm)
