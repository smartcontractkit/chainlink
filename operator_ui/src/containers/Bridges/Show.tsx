import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import Divider from '@material-ui/core/Divider'
import Grid from '@material-ui/core/Grid'
import Typography from '@material-ui/core/Typography'
import { fetchBridgeSpec } from 'actions'
import BaseLink from 'components/BaseLink'
import Button from 'components/Button'
import Content from 'components/Content'
import { AppState } from 'connectors/redux/reducers'
import { BridgeType } from 'operator_ui'
import React from 'react'
import { connect } from 'react-redux'
import bridgeSelector from 'selectors/bridge'
import { useEffect, useHooks } from 'use-react-hooks'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'

const Loading = () => <div>Loading...</div>

interface LoadedProps {
  bridge: BridgeType
}

const fields: [string, string][] = [
  ['name', 'Name'],
  ['url', 'URL'],
  ['confirmations', 'Confirmations'],
  ['minimumContractPayment', 'Minimum Contract Payment'],
  ['outgoingToken', 'Outgoing Token'],
]

const Loaded = ({ bridge }: LoadedProps) => (
  <CardContent>
    {fields.map(([k, t]) => {
      return (
        <React.Fragment key={k}>
          <Typography variant="subtitle1" color="textSecondary">
            {t}
          </Typography>
          <Typography variant="body1" color="inherit">
            {bridge[k as keyof typeof bridge]}
          </Typography>
        </React.Fragment>
      )
    })}
  </CardContent>
)

interface Props {
  match: {
    params: {
      bridgeId: string
    }
  }
  fetchBridgeSpec: (name: string) => Promise<any>
  bridge?: BridgeType
}

export const Show = useHooks(({ bridge, fetchBridgeSpec, match }: Props) => {
  useEffect(() => {
    document.title = 'Show Bridge'
    fetchBridgeSpec(match.params.bridgeId)
  }, [])

  return (
    <Content>
      <Grid container>
        <Grid item xs={12} md={11} lg={9}>
          <Card>
            <CardContent>
              <Grid container>
                <Grid item xs={9}>
                  <Typography variant="h5" color="secondary">
                    Bridge Info
                  </Typography>
                </Grid>
                <Grid item xs={3}>
                  <Grid container justify="flex-end">
                    <Grid item>
                      {bridge && (
                        <Button
                          variant="secondary"
                          component={BaseLink}
                          href={`/bridges/${bridge.id}/edit`}
                        >
                          Edit
                        </Button>
                      )}
                    </Grid>
                  </Grid>
                </Grid>
              </Grid>
            </CardContent>

            <Divider />

            {bridge ? <Loaded bridge={bridge} /> : <Loading />}
          </Card>
        </Grid>
      </Grid>
    </Content>
  )
})

const mapStateToProps = (state: AppState, ownProps: Props) => ({
  bridge: bridgeSelector(state, ownProps.match.params.bridgeId),
})

export const ConnectedShow = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({ fetchBridgeSpec }),
)(Show)

export default ConnectedShow
