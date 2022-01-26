import React from 'react'

import { Field, Form, Formik, FormikHelpers } from 'formik'
import { TextField } from 'formik-material-ui'
import * as Yup from 'yup'

import Button from '@material-ui/core/Button'
import Grid from '@material-ui/core/Grid'

export interface FormValues {
  name: string
  minimumContractPayment: string
  confirmations: number
  url: string
}

const ValidationSchema = Yup.object().shape({
  name: Yup.string().required('Required'),
  url: Yup.string().required('Required'),
})

export interface Props {
  initialValues: FormValues
  submitButtonText: string
  nameDisabled?: boolean
  onSubmit: (
    values: FormValues,
    formikHelpers: FormikHelpers<FormValues>,
  ) => void | Promise<any>
}

export const BridgeForm = ({
  initialValues,
  onSubmit,
  submitButtonText,
  nameDisabled = false,
}: Props) => {
  return (
    <Formik
      initialValues={initialValues}
      validationSchema={ValidationSchema}
      onSubmit={onSubmit}
    >
      {({ isSubmitting }) => (
        <>
          <Form data-testid="bridge-form" noValidate>
            <Grid container spacing={16}>
              <Grid item xs={12} md={7}>
                <Field
                  component={TextField}
                  id="name"
                  name="name"
                  label="Name"
                  disabled={nameDisabled}
                  required
                  fullWidth
                  FormHelperTextProps={{ 'data-testid': 'name-helper-text' }}
                />
              </Grid>

              <Grid item xs={12} md={7}>
                <Field
                  component={TextField}
                  id="url"
                  name="url"
                  label="Bridge URL"
                  placeholder="https://"
                  required
                  fullWidth
                  FormHelperTextProps={{ 'data-testid': 'url-helper-text' }}
                />
              </Grid>

              <Grid item xs={12} md={7}>
                <Grid container spacing={16}>
                  <Grid item xs={7}>
                    <Field
                      component={TextField}
                      id="minimumContractPayment"
                      name="minimumContractPayment"
                      label="Minimum Contract Payment"
                      placeholder="0"
                      fullWidth
                      inputProps={{ min: 0 }}
                      FormHelperTextProps={{
                        'data-testid': 'minimumContractPayment-helper-text',
                      }}
                    />
                  </Grid>
                  <Grid item xs={7}>
                    <Field
                      component={TextField}
                      id="confirmations"
                      name="confirmations"
                      label="Confirmations"
                      placeholder="0"
                      type="number"
                      fullWidth
                      inputProps={{ min: 0 }}
                      FormHelperTextProps={{
                        'data-testid': 'confirmations-helper-text',
                      }}
                    />
                  </Grid>
                </Grid>
              </Grid>

              <Grid item xs={12} md={7}>
                <Button
                  variant="contained"
                  color="primary"
                  type="submit"
                  disabled={isSubmitting}
                  size="large"
                >
                  {submitButtonText}
                </Button>
              </Grid>
            </Grid>
          </Form>
        </>
      )}
    </Formik>
  )
}
