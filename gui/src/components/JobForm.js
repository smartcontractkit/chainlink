import React from 'react'
import { withFormik, Form } from 'formik'
import { withStyles } from '@material-ui/core/styles'
import Button from '@material-ui/core/Button'
import { TextField, Grid, Typography } from '@material-ui/core'
import { connect } from 'react-redux'
import { submitCreate } from 'actions'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import Flash from './Flash'

const styles = theme => ({
  jsonfield: {
    paddingTop: theme.spacing.unit * 1.25
  },
  card: {
    paddingBottom: theme.spacing.unit * 2
  },
  flash: {
    textAlign: 'center'
  }
})

const FormLayout = ({ isSubmitting, classes, handleChange, creating, errors }) => (
  <Grid justify='center' container spacing={24}>
    {errors.length > 0 && (
      <Flash error className={classes.flash}>
        {errors.map((msg, i) => <p key={i}>{msg}</p>)}
      </Flash>
    )}
    <Grid item xs={5}>
      <Form noValidate>
        <TextField
          fullWidth
          onChange={handleChange}
          label='Paste JSON'
          placeholder='Paste JSON'
          multiline
          className={classes.jsonfield}
          margin='normal'
          type='json'
          name='json'
        />
        <Grid container alignContent='center' direction='column'>
          <Grid item>
            <Button color='primary' type='submit' disabled={isSubmitting}>
              Build Job
            </Button>
          </Grid>
          {
            creating && (
              <Grid item xs>
                <Typography variant='body1' color='textSecondary' align='center'>
                Creating...
                </Typography>
              </Grid>
            )}
        </Grid>
      </Form>
    </Grid>
  </Grid>
)

const JobForm = withFormik({
  mapPropsToValues ({ json }) {
    return {
      json: json || ''
    }
  },
  handleSubmit (values, { props }) {
    props.submitCreate('v2/specs', values.json).then(e => console.log(e))
  }
})(FormLayout)

const mapStateToProps = state => ({
  creating: state.session.fetching
})

export const ConnectedJobForm = connect(mapStateToProps, matchRouteAndMapDispatchToProps({ submitCreate }))(JobForm)

export default withStyles(styles)(ConnectedJobForm)
