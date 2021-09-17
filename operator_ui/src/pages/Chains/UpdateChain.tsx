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
import { ChainSpecV2 } from './RegionalNav'
import { Field, Form, Formik } from 'formik'
import { TextField } from 'formik-material-ui'
import * as Yup from 'yup'
import { values } from 'lodash'

const SuccessNotification = ({ id }: { id: string }) => (
  <>
    Successfully created node{' '}
    <BaseLink id="created-node" href={`/nodes`}>
      {id}
    </BaseLink>
  </>
)

function apiCall({
  chain,
  config,
}: {
  chain: ChainSpecV2
  config: Record<string, JSONPrimitive>
}): Promise<ApiResponse<Chain>> {
  const definition: UpdateChainRequest = { config }
  return api.v2.chains.updateChain(chain.id, definition)
}

const UpdateChain = ({ chain }: { chain: ChainSpecV2 }) => {
  const dispatch = useDispatch()

  async function handleSubmit({ config }: { config: string }) {
    apiCall({
      chain,
      config: JSON.parse(config),
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
                {({ isSubmitting, submitForm }) => (
                  <Form>
                    <Grid container spacing={16}>
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
