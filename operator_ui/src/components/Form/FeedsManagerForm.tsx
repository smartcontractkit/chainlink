import React from 'react'
import { Field, Form, Formik, FormikHelpers } from 'formik'
import { TextField } from 'formik-material-ui'
import * as Yup from 'yup'

import Button from '@material-ui/core/Button'
import Grid from '@material-ui/core/Grid'

export type FormValues = {
  name: string
  uri: string
  publicKey: string
}

const ValidationSchema = Yup.object().shape({
  name: Yup.string().required('Required'),
  uri: Yup.string().required('Required'),
  publicKey: Yup.string().required('Required'),
})

export interface Props {
  initialValues: FormValues
  onSubmit: (
    values: FormValues,
    formikHelpers: FormikHelpers<FormValues>,
  ) => void | Promise<any>
}

export const FeedsManagerForm: React.FC<Props> = ({
  initialValues,
  onSubmit,
}) => {
  return (
    <Formik
      initialValues={initialValues}
      validationSchema={ValidationSchema}
      onSubmit={onSubmit}
    >
      {({ isSubmitting, submitForm }) => (
        <Form data-testid="feeds-manager-form">
          <Grid container spacing={16}>
            <Grid item xs={12} md={6}>
              <Field
                component={TextField}
                id="name"
                name="name"
                label="Name"
                required
                fullWidth
                FormHelperTextProps={{ 'data-testid': 'name-helper-text' }}
              />
            </Grid>

            <Grid item xs={false} md={6}></Grid>

            <Grid item xs={12} md={6}>
              <Field
                component={TextField}
                id="uri"
                name="uri"
                label="URI"
                required
                fullWidth
                helperText="Provided by the Feeds Manager operator"
                FormHelperTextProps={{ 'data-testid': 'uri-helper-text' }}
              />
            </Grid>

            <Grid item xs={12} md={6}>
              <Field
                component={TextField}
                id="publicKey"
                name="publicKey"
                label="Public Key"
                required
                fullWidth
                helperText="Provided by the Feeds Manager operator"
                FormHelperTextProps={{ 'data-testid': 'publicKey-helper-text' }}
              />
            </Grid>

            <Grid item xs={12}>
              <Button
                variant="contained"
                color="primary"
                disabled={isSubmitting}
                onClick={submitForm}
              >
                Submit
              </Button>
            </Grid>
          </Grid>
        </Form>
      )}
    </Formik>
  )
}
