import React from 'react'
import { withFormik, Form } from 'formik'
import { withStyles } from '@material-ui/core/styles'
import Button from '@material-ui/core/Button'
import { TextField, Grid } from '@material-ui/core'
import postJob from 'utils/postJob'

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
      <Form>
        <div>
          <TextField
            fullWidth
            id='textarea'
            label='Paste JSON'
            onChange={handleChange}
            placeholder='Paste JSON'
            multiline
            className={classes.jsonfield}
            margin='normal'
          />
        </div>
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
    postJob(values).then(res => console.log(res))
  }
})(FormLayout)

export default withStyles(styles)(JobForm)
