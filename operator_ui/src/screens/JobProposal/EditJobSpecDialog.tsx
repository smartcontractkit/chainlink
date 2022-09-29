import React from 'react'
import { Field, Form, Formik, FormikHelpers } from 'formik'
import { TextField } from 'formik-material-ui'
import * as Yup from 'yup'

import Button from '@material-ui/core/Button'
import Dialog from '@material-ui/core/Dialog'
import DialogActions from '@material-ui/core/DialogActions'
import DialogContent from '@material-ui/core/DialogContent'
import DialogTitle from '@material-ui/core/DialogTitle'
import { createStyles, WithStyles, withStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'

const styles = () => {
  return createStyles({
    paperRoot: {
      width: 700,
    },
  })
}

export type FormValues = {
  id: string
  definition: string
}

export interface Props extends WithStyles<typeof styles> {
  onClose: () => void
  open: boolean
  initialValues: FormValues
  onSubmit: (
    values: FormValues,
    formikHelpers: FormikHelpers<FormValues>,
  ) => void | Promise<any>
}

const ValidationSchema = Yup.object().shape({
  definition: Yup.string().required('Required'),
})

export const EditJobSpecDialog = withStyles(styles)(
  ({ classes, initialValues, onClose, onSubmit, open }: Props) => {
    return (
      <Formik
        initialValues={initialValues}
        validationSchema={ValidationSchema}
        onSubmit={async (values, formikHelper) => {
          await onSubmit(values, formikHelper)

          onClose()
        }}
      >
        {({ isSubmitting, submitForm }) => (
          <Form>
            <Dialog
              open={open}
              onClose={onClose}
              classes={{ paper: classes.paperRoot }}
            >
              <DialogTitle disableTypography>
                <Typography variant="h5">Edit Job Spec</Typography>
              </DialogTitle>
              <DialogContent>
                <Field
                  component={TextField}
                  id="definition"
                  name="definition"
                  label="Job Spec"
                  variant="outlined"
                  multiline
                  rows={10}
                  rowsMax={25}
                  required
                  autoComplete="off"
                  margin="normal"
                  fullWidth
                  spellCheck="false"
                />
              </DialogContent>
              <DialogActions>
                <Button onClick={onClose} variant="text">
                  Cancel
                </Button>
                <Button
                  variant="contained"
                  color="primary"
                  disabled={isSubmitting}
                  onClick={submitForm}
                >
                  Submit
                </Button>
              </DialogActions>
            </Dialog>
          </Form>
        )}
      </Formik>
    )
  },
)
