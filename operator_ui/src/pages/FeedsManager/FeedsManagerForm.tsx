import React from 'react'
import { Field, Form, Formik } from 'formik'
import { TextField, CheckboxWithLabel } from 'formik-material-ui'
import * as Yup from 'yup'

import Button from '@material-ui/core/Button'
import FormControl from '@material-ui/core/FormControl'
import FormGroup from '@material-ui/core/FormGroup'
import FormLabel from '@material-ui/core/FormLabel'
import Grid from '@material-ui/core/Grid'

const jobTypes = [
  {
    label: 'Flux Monitor',
    value: 'fluxmonitor',
  },
  {
    label: 'OCR',
    value: 'ocr',
  },
]

export type FormValues = {
  name: string
  uri: string
  jobTypes: string[]
  publicKey: string
  isBootstrapPeer: boolean
  bootstrapPeerMultiaddr?: string
}

const ValidationSchema = Yup.object().shape({
  name: Yup.string().required('Required'),
  uri: Yup.string().required('Required'),
  publicKey: Yup.string().required('Required'),
  bootstrapPeerMultiaddr: Yup.string()
    .when('isBootstrapPeer', {
      is: true,
      then: Yup.string().required('Required').nullable(),
    })
    .nullable(),
})

interface Props {
  initialValues: FormValues
  onSubmit: (values: FormValues) => void
}

export const FeedsManagerForm: React.FC<Props> = ({
  initialValues,
  onSubmit,
}) => {
  return (
    <Formik
      initialValues={initialValues}
      validationSchema={ValidationSchema}
      onSubmit={async (values) => {
        onSubmit(values)
      }}
    >
      {({ isSubmitting, submitForm, values }) => (
        <Form>
          <Grid container spacing={16}>
            <Grid item xs={12} md={6}>
              <Field
                component={TextField}
                id="name"
                name="name"
                label="Name"
                required
                fullWidth
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
              />
            </Grid>

            <Grid item xs={12} md={7}>
              <FormControl>
                <FormLabel>Which job types does this node run?</FormLabel>
                <FormGroup>
                  <div style={{ display: 'flex', flexDirection: 'row' }}>
                    {jobTypes.map((jobType) => (
                      <Field
                        type="checkbox"
                        component={CheckboxWithLabel}
                        name="jobTypes"
                        key={jobType.value}
                        value={jobType.value}
                        Label={{ label: jobType.label }}
                      />
                    ))}
                  </div>
                </FormGroup>
              </FormControl>
            </Grid>

            <Grid item xs={12} md={7}>
              <FormControl>
                <FormGroup>
                  <div style={{ display: 'flex', flexDirection: 'row' }}>
                    <Field
                      type="checkbox"
                      component={CheckboxWithLabel}
                      name="isBootstrapPeer"
                      Label={{
                        label: 'Is this node running as a bootstrap peer?',
                      }}
                    />
                  </div>
                </FormGroup>
              </FormControl>
            </Grid>

            {values.isBootstrapPeer && (
              <Grid item xs={12} md={7}>
                <Field
                  component={TextField}
                  id="bootstrapPeerMultiaddr"
                  name="bootstrapPeerMultiaddr"
                  label="Bootstrap Peer Multiaddress"
                  fullWidth
                  helperText=""
                />
              </Grid>
            )}

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
