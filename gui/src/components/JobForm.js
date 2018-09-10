import React, { Fragment } from 'react'
import { withFormik, Form } from 'formik'
import { withStyles } from '@material-ui/core/styles'
import Button from '@material-ui/core/Button'
import { TextField, Grid } from '@material-ui/core'
import { connect } from 'react-redux'
import { submitJobSpec } from 'actions'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
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

const JobFormLayout = ({
  isSubmitting,
  classes,
  handleChange,
  values,
  submitCount
}) => {
  return (
    <Fragment>
      <Prompt
        when={values.json !== '' && submitCount === 0}
        message='You have not submitted the form, are you sure you want to leave?'
      />
      <Form noValidate>
        <Grid container justify='center'>
          <Grid container justify='center'>
            <Grid item sm={8}>
              <TextField value={values.json} onChange={handleChange} fullWidth label='Paste JSON' rows={10} rowsMax={25} placeholder='Paste JSON' multiline margin='normal' name='json' id='json' />
            </Grid>
          </Grid>
          <Button className={classes.button} variant='contained' color='primary' type='submit' disabled={isSubmitting || !values.json}>
            Build Job
          </Button>
        </Grid>
      </Form>
    </Fragment>
  )
}

const JobForm = withFormik({
  mapPropsToValues ({ fromJson }) {
    return {
      json: JSON.stringify(fromJson, null, '\t') || ''
    }
  },
  handleSubmit (values, { props, setSubmitting }) {
    props.submitJobSpec(values.json.trim(), false)
    setTimeout(() => { setSubmitting(false) }, 1000)
  }
})(JobFormLayout)

const mapStateToProps = state => {
  return {
    networkError: state.create.networkError
  }
}

export const ConnectedJobForm = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({ submitJobSpec })
)(JobForm)

export default withStyles(styles)(ConnectedJobForm)
