import React from 'react'

import TOML from '@iarna/toml'
import { Field, Form, Formik, FormikHelpers } from 'formik'
import { TextField } from 'formik-material-ui'
import * as Yup from 'yup'

import Button from '@material-ui/core/Button'
import Grid from '@material-ui/core/Grid'

export interface FormValues {
  toml: string
}

const ValidationSchema = Yup.object().shape({
  toml: Yup.string()
    .required('Required')
    .test('toml', 'Invalid TOML', function (value = '') {
      try {
        TOML.parse(value)

        return true
      } catch {
        return false
      }
    }),
})

export interface Props {
  initialValues: FormValues
  onSubmit: (
    values: FormValues,
    formikHelpers: FormikHelpers<FormValues>,
  ) => void | Promise<any>
  onTOMLChange?: (toml: string) => void
}

export const JobForm = ({ initialValues, onSubmit, onTOMLChange }: Props) => {
  return (
    <Formik
      initialValues={initialValues}
      validationSchema={ValidationSchema}
      onSubmit={onSubmit}
    >
      {({ isSubmitting, values }) => {
        if (onTOMLChange) {
          onTOMLChange(values.toml)
        }

        return (
          <Form data-testid="job-form" noValidate>
            <Grid container spacing={16}>
              <Grid item xs={12}>
                <Field
                  component={TextField}
                  id="toml"
                  name="toml"
                  label="Job Spec (TOML)"
                  required
                  fullWidth
                  multiline
                  rows={10}
                  rowsMax={25}
                  variant="outlined"
                  autoComplete="off"
                  FormHelperTextProps={{ 'data-testid': 'toml-helper-text' }}
                />
              </Grid>

              <Grid item xs={12} md={7}>
                <Button
                  variant="contained"
                  color="primary"
                  type="submit"
                  disabled={isSubmitting}
                  size="large"
                >
                  Create Job
                </Button>
              </Grid>
            </Grid>
          </Form>
        )
      }}
    </Formik>
  )
}
