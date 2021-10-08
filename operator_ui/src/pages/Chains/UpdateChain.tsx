import React from 'react'
import { ApiResponse } from 'utils/json-api-client'
import Button from '@material-ui/core/Button'
import * as api from 'api'
import { useDispatch } from 'react-redux'
import { Chain, UpdateChainRequest } from 'core/store/models'
import BaseLink from 'components/BaseLink'
import ErrorMessage from 'components/Notifications/DefaultError'
import { notifySuccess, notifyError } from 'actionCreators'
import Content from 'components/Content'
import { Grid, Card, CardContent, CardHeader } from '@material-ui/core'
import { ChainResource } from './RegionalNav'
import { Field, Form, Formik } from 'formik'
import { TextField, CheckboxWithLabel } from 'formik-material-ui'
import * as Yup from 'yup'

const SuccessNotification = ({ id }: { id: string }) => (
  <>
    Successfully updated chain{' '}
    <BaseLink id="updated-chain" href={`/chains`}>
      {id}
    </BaseLink>
  </>
)

function apiCall({
  chain,
  config,
  enabled,
}: {
  chain: ChainResource
  config: Record<string, JSONPrimitive>
  enabled: boolean
}): Promise<ApiResponse<Chain>> {
  const definition: UpdateChainRequest = { config, enabled }
  return api.v2.chains.updateChain(chain.id, definition)
}

const UpdateChain = ({ chain }: { chain: ChainResource }) => {
  const dispatch = useDispatch()

  async function handleSubmit({
    config,
    enabled,
  }: {
    config: string
    enabled: boolean
  }) {
    apiCall({
      chain,
      config: JSON.parse(config),
      enabled,
    })
      .then(({ data }) => {
        dispatch(notifySuccess(SuccessNotification, data))
      })
      .catch((error) => {
        dispatch(notifyError(ErrorMessage, error))
      })
  }

  const configOverrides = Object.fromEntries(
    Object.entries(chain.attributes.config).filter(
      ([_key, value]) => value !== null,
    ),
  )

  const initialValues = {
    config: JSON.stringify(configOverrides, null, 2),
    enabled: chain.attributes.enabled,
  }

  const ValidationSchema = Yup.object().shape({
    config: Yup.string().required('Required'),
  })

  return (
    <Content>
      <Grid container spacing={40}>
        <Grid item xs={12}>
          <Card>
            <CardHeader title={`Edit Chain ${chain.id}`} />
            <CardContent>
              <Formik
                initialValues={initialValues}
                validationSchema={ValidationSchema}
                onSubmit={async (values) => {
                  handleSubmit(values)
                }}
              >
                {({ isSubmitting, submitForm, values }) => (
                  <Form>
                    <Grid container spacing={16}>
                      <Grid item xs={12} md={4}>
                        <Field
                          type="checkbox"
                          component={CheckboxWithLabel}
                          name="enabled"
                          id="enabled"
                          checked={values.enabled}
                          Label={{ label: 'Enabled' }}
                        />
                      </Grid>
                      <Grid item xs={false} md={8}></Grid>
                      <Grid item xs={12} md={4}>
                        <Field
                          component={TextField}
                          autoComplete="off"
                          label="Config Overrides"
                          rows={10}
                          rowsMax={25}
                          multiline
                          margin="normal"
                          name="config"
                          id="config"
                          variant="outlined"
                          fullWidth
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
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Content>
  )
}

export default UpdateChain
