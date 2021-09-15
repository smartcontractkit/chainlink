import React, { useState } from 'react'
import { ApiResponse, BadRequestError } from 'utils/json-api-client'
import Button from 'components/Button'
import * as api from 'api'
import { useDispatch } from 'react-redux'
import { CreateNodeRequest, Node } from 'core/store/models'
import BaseLink from 'components/BaseLink'
import ErrorMessage from 'components/Notifications/DefaultError'
import { notifySuccess, notifyError } from 'actionCreators'
import Content from 'components/Content'
import {
  TextField,
  Grid,
  Card,
  CardContent,
  CardHeader,
  CircularProgress,
} from '@material-ui/core'
import { createStyles, withStyles, WithStyles } from '@material-ui/core/styles'
import { ChainSpecV2 } from './RegionalNav'

const styles = () =>
  createStyles({
    loader: {
      position: 'absolute',
    },
  })

const SuccessNotification = ({ id }: { id: string }) => (
  <>
    Successfully created node{' '}
    <BaseLink id="created-node" href={`/nodes`}>
      {id}
    </BaseLink>
  </>
)

function apiCall({
  name,
  wsURL,
  httpURL,
  evmChainID,
}: {
  name: string
  httpURL: string
  wsURL: string
  evmChainID: string
}): Promise<ApiResponse<Node>> {
  const definition: CreateNodeRequest = { name, wsURL, httpURL, evmChainID }
  return api.v2.nodes.createNode(definition)
}

const NewChainNode = ({
  classes,
  chain,
}: {
  classes: WithStyles<typeof styles>['classes']
  chain: ChainSpecV2
}) => {
  const dispatch = useDispatch()
  const [name, setName] = useState<string>('')
  const [nameErrorMsg, setNameErrorMsg] = useState<string>('')
  const [httpURL, setHttpURL] = useState<string>('')
  const [httpURLErrorMsg, setHttpURLErrorMsg] = useState<string>('')
  const [wsURL, setWsURL] = useState<string>('')
  const [wsURLErrorMsg, setWsURLErrorMsg] = useState<string>('')
  const [loading, setLoading] = useState<boolean>(false)

  function validate({
    name,
    httpURL,
    wsURL,
  }: {
    name: string
    httpURL: string
    wsURL: string
  }) {
    let valid = true
    if (!name) {
      setNameErrorMsg('Invalid name')
      valid = false
    }
    try {
      const url = new URL(httpURL)
      if (!(url.protocol === 'http:' || url.protocol === 'https:')) {
        setHttpURLErrorMsg('Invalid HTTP URL')
        valid = false
      }
    } catch (_) {
      setHttpURLErrorMsg('Invalid HTTP URL')
      valid = false
    }
    try {
      const url = new URL(wsURL)
      if (!(url.protocol === 'ws:' || url.protocol === 'wss:')) {
        setWsURLErrorMsg('Invalid Websocket URL')
        valid = false
      }
    } catch (_) {
      setWsURLErrorMsg('Invalid Websocket URL')
      valid = false
    }
    return valid
  }

  function handleNameChange(event: React.ChangeEvent<HTMLTextAreaElement>) {
    setName(event.target.value)
    setNameErrorMsg('')
  }
  function handlehttpURLChange(event: React.ChangeEvent<HTMLTextAreaElement>) {
    setHttpURL(event.target.value)
    setHttpURLErrorMsg('')
  }
  function handlewsURLChange(event: React.ChangeEvent<HTMLTextAreaElement>) {
    setWsURL(event.target.value)
    setWsURLErrorMsg('')
  }

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const isValid = validate({ name, httpURL, wsURL })

    if (isValid) {
      setLoading(true)
      apiCall({
        name,
        wsURL,
        httpURL,
        evmChainID: chain.id,
      })
        .then(({ data }) => {
          dispatch(notifySuccess(SuccessNotification, data))
        })
        .catch((error) => {
          dispatch(notifyError(ErrorMessage, error))
          if (error instanceof BadRequestError) {
            setNameErrorMsg('Invalid Name')
          }
        })
        .finally(() => {
          setLoading(false)
        })
    }
  }

  return (
    <Content>
      <Grid container spacing={40}>
        <Grid item xs={12}>
          <Card>
            <CardHeader title="New Node" />
            <CardContent>
              <form noValidate onSubmit={handleSubmit}>
                <Grid container>
                  <Grid item xs={12}>
                    <TextField
                      error={Boolean(nameErrorMsg)}
                      helperText={Boolean(nameErrorMsg) && nameErrorMsg}
                      label="Name"
                      name="Name"
                      placeholder="Name"
                      value={name}
                      onChange={handleNameChange}
                    />
                  </Grid>
                  <Grid item xs={12}>
                    <TextField
                      error={Boolean(httpURLErrorMsg)}
                      helperText={Boolean(httpURLErrorMsg) && httpURLErrorMsg}
                      label="HTTP URL"
                      name="httpURL"
                      placeholder="httpURL"
                      value={httpURL}
                      onChange={handlehttpURLChange}
                    />
                  </Grid>
                  <Grid item xs={12}>
                    <TextField
                      error={Boolean(wsURLErrorMsg)}
                      helperText={Boolean(wsURLErrorMsg) && wsURLErrorMsg}
                      label="Websocket URL"
                      name="wsURL"
                      placeholder="wsURL"
                      value={wsURL}
                      onChange={handlewsURLChange}
                    />
                  </Grid>
                  <Grid item xs={12}>
                    <Button
                      data-testid="new-node-config-submit"
                      variant="primary"
                      type="submit"
                      size="large"
                      disabled={loading || Boolean(nameErrorMsg)}
                    >
                      Create Node
                      {loading && (
                        <CircularProgress
                          className={classes.loader}
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

export default withStyles(styles)(NewChainNode)
