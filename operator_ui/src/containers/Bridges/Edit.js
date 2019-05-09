import React from 'react'
import PropTypes from 'prop-types'
import { Link } from 'react-router-dom'
import { connect } from 'react-redux'
import Grid from '@material-ui/core/Grid'
import Typography from '@material-ui/core/Typography'
import Form from 'components/Bridges/Form'
import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import Divider from '@material-ui/core/Divider'
import ErrorMessage from 'components/Notifications/DefaultError'
import Button from 'components/Button'
import bridgeSelector from 'selectors/bridge'
import { fetchBridgeSpec, updateBridge } from 'actions'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import ReactStaticLinkComponent from 'components/ReactStaticLinkComponent'
import Content from 'components/Content'
import { useHooks, useEffect } from 'use-react-hooks'

const SuccessNotification = ({ id }) => {
  return (
    <React.Fragment>
      Successfully updated <Link to={`/bridges/${id}`}>{id}</Link>
    </React.Fragment>
  )
}

export const Edit = useHooks(props => {
  useEffect(() => {
    document.title = 'Edit Bridge'
    const { fetchBridgeSpec, match } = props
    fetchBridgeSpec(match.params.bridgeId)
  }, [])
  const { bridge, updateBridge } = props
  const checkLoaded = () => bridge
  const onLoad = buildLoadedComponent => {
    if (checkLoaded()) return buildLoadedComponent(props)
    return <div>Loading...</div>
  }

  return (
    <Content>
      <Grid container>
        <Grid item xs={12} md={11} lg={9}>
          <Card>
            <CardContent>
              <Grid container>
                <Grid item xs={9}>
                  <Typography variant="h5" color="secondary">
                    Edit Bridge
                  </Typography>
                </Grid>
                <Grid item xs={3}>
                  <Grid container justify="flex-end">
                    <Grid item>
                      {bridge && (
                        <Button
                          component={ReactStaticLinkComponent}
                          to={`/bridges/${bridge.id}`}
                        >
                          Cancel
                        </Button>
                      )}
                    </Grid>
                  </Grid>
                </Grid>
              </Grid>
            </CardContent>

            <Divider />

            <CardContent>
              {onLoad(({ bridge }) => (
                <Form
                  actionText="Save Bridge"
                  onSubmit={updateBridge}
                  name={bridge.name}
                  nameDisabled
                  url={bridge.url}
                  confirmations={bridge.confirmations}
                  minimumContractPayment={bridge.minimumContractPayment}
                  onSuccess={SuccessNotification}
                  onError={ErrorMessage}
                />
              ))}
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Content>
  )
})

Edit.propTypes = {
  bridge: PropTypes.object
}

const mapStateToProps = (state, ownProps) => {
  const bridge = bridgeSelector(state, ownProps.match.params.bridgeId)
  return { bridge }
}

export const ConnectedEdit = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({ fetchBridgeSpec, updateBridge })
)(Edit)

export default ConnectedEdit
