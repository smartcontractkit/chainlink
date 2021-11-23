import React, { useEffect } from 'react'

import { useQuery } from '@apollo/client'
import PropTypes from 'prop-types'
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
import { fetchBridgeSpec, updateBridge } from 'actionCreators'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import BaseLink from 'components/BaseLink'
import Content from 'components/Content'
import { BRIDGES_QUERY } from 'src/screens/Bridges/BridgesScreen'

const SuccessNotification = ({ id }) => {
  return (
    <React.Fragment>
      Successfully updated <BaseLink href={`/bridges/${id}`}>{id}</BaseLink>
    </React.Fragment>
  )
}

export const Edit = (props) => {
  // This is a hacky fix to refetch the page data after an edit so changes
  // appear on the table. Once this gets moved to GQL, we can use the
  // refetchQueries on the mutation.
  const { refetch } = useQuery(BRIDGES_QUERY, {
    variables: { offset: 0, limit: 10 },
  })

  const { fetchBridgeSpec, match, bridge, updateBridge } = props
  useEffect(() => {
    document.title = 'Edit Bridge'
    fetchBridgeSpec(match.params.bridgeId)
  }, [fetchBridgeSpec, match])

  const checkLoaded = () => bridge
  const onLoad = (buildLoadedComponent) => {
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
                          component={BaseLink}
                          href={`/bridges/${bridge.id}`}
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
                  refetchGQL={refetch}
                />
              ))}
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Content>
  )
}

Edit.propTypes = {
  bridge: PropTypes.object,
}

const mapStateToProps = (state, ownProps) => {
  const bridge = bridgeSelector(state, ownProps.match.params.bridgeId)
  return { bridge }
}

export const ConnectedEdit = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({ fetchBridgeSpec, updateBridge }),
)(Edit)

export default ConnectedEdit
