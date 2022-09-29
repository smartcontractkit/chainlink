import React, { useState } from 'react'
import { ApiResponse } from 'utils/json-api-client'
import Button from '@material-ui/core/Button'
import * as api from 'api'
import { useDispatch } from 'react-redux'
import { Chain, UpdateChainRequest } from 'core/store/models'
import BaseLink from 'components/BaseLink'
import ErrorMessage from 'components/Notifications/DefaultError'
import { notifyError, notifySuccess } from 'actionCreators'
import Content from 'components/Content'
import {
  Card,
  CardContent,
  CardHeader,
  Checkbox,
  CircularProgress,
  FormControlLabel,
  Grid,
} from '@material-ui/core'
import { ChainResource } from './RegionalNav'
import ChainConfigFields, {
  ConfigOverrides,
} from 'pages/Chains/ChainConfigFields'

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

  const [enabled, setEnabled] = useState<boolean>(chain.attributes.enabled)
  const [overrides, setOverrides] = useState<ConfigOverrides>({})
  const [keySpecificOverridesErrorMsg, setKeySpecificOverridesErrorMsg] =
    useState<string>('')
  const [loading, setLoading] = useState<boolean>(false)

  function onConfigChange(config: ConfigOverrides, error: string) {
    if (error) {
      setKeySpecificOverridesErrorMsg(error)
      return
    }

    setOverrides({ ...overrides, ...config })
  }

  function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()

    if (keySpecificOverridesErrorMsg) {
      return
    }

    setLoading(true)

    apiCall({
      chain,
      enabled,
      config: overrides,
    })
      .then(({ data }) => {
        dispatch(notifySuccess(SuccessNotification, data))
      })
      .catch((error) => {
        dispatch(notifyError(ErrorMessage, error))
      })
      .finally(() => {
        setLoading(false)
      })
  }

  const configOverrides = Object.fromEntries(
    Object.entries(chain.attributes.config).filter(
      ([_key, value]) => value !== null,
    ),
  )

  const initialValues = {
    config: configOverrides,
    enabled: chain.attributes.enabled,
  }

  return (
    <Content>
      <Grid container spacing={40}>
        <Grid item xs={12}>
          <Card>
            <CardHeader title={`Edit Chain ${chain.id}`} />
            <CardContent>
              <form noValidate onSubmit={handleSubmit}>
                <Grid container spacing={16}>
                  <Grid item xs={12} md={4}>
                    <FormControlLabel
                      control={
                        <Checkbox
                          name="enabled"
                          checked={enabled}
                          value={enabled}
                          onChange={(event) => setEnabled(event.target.checked)}
                        />
                      }
                      label="Enabled"
                    />
                  </Grid>

                  <Grid item xs={false} md={8}></Grid>

                  <ChainConfigFields
                    onChange={onConfigChange}
                    initialValues={initialValues.config}
                  />

                  <Grid item xs={12}>
                    <Button
                      variant="contained"
                      color="primary"
                      type="submit"
                      disabled={
                        loading || Boolean(keySpecificOverridesErrorMsg)
                      }
                    >
                      Submit
                      {loading && (
                        <CircularProgress
                          style={{ position: 'absolute' }}
                          size={30}
                          color="primary"
                        />
                      )}
                    </Button>
                  </Grid>
                </Grid>
              </form>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Content>
  )
}

export default UpdateChain
