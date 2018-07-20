import React from 'react'
import { withFormik, Form } from 'formik'
import { withStyles } from '@material-ui/core/styles'
import Button from '@material-ui/core/Button'
import { TextField, Grid } from '@material-ui/core'
import { postJob } from 'api'

const styles = theme => ({
  jsonfield: {
    paddingTop: theme.spacing.unit * 1.25
  },
  card: {
    paddingBottom: theme.spacing.unit * 2
  }
})

const FormLayout = ({ isSubmitting, classes, handleChange }) => (
  <Grid justify='center' container spacing={24}>
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
        <Grid container justify='center'>
          <Button color='primary' type='submit' disabled={isSubmitting}>
              Build Job
          </Button>
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
  handleSubmit (values) {
    postJob(values.json)
  }
})(FormLayout)

export default withStyles(styles)(JobForm)
