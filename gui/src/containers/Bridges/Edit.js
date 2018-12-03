import React from 'react'
import PropTypes from 'prop-types'
import { Link } from 'react-static'
import { connect } from 'react-redux'
import Grid from '@material-ui/core/Grid'
import { Button } from '@material-ui/core'
import Title from 'components/Title'
import PaddedCard from 'components/PaddedCard'
import BridgesForm from 'components/Bridges/Form'
import Breadcrumb from 'components/Breadcrumb'
import BreadcrumbItem from 'components/BreadcrumbItem'
import ErrorMessage from 'components/Notifications/DefaultError'
import bridgeSelector from 'selectors/bridge'
import {
  fetchBridgeSpec,
  updateBridge
} from 'actions'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import ReactStaticLinkComponent from 'components/ReactStaticLinkComponent'
import Content from 'components/Content'
import { useHooks, useEffect } from 'use-react-hooks'

const SuccessNotification = ({name}) => (<React.Fragment>
  Successfully updated <Link to={`/bridges/${name}`}>{name}</Link>
</React.Fragment>)

export const Edit = useHooks(props => {
  useEffect(() => {
    const {fetchBridgeSpec, match} = props
    fetchBridgeSpec(match.params.bridgeId)
  }, [])
  const {bridge, updateBridge} = props
  const checkLoaded = () => bridge
  const onLoad = (buildLoadedComponent) => {
    if (checkLoaded()) return buildLoadedComponent(props)
    return <div>Loading...</div>
  }

  return (
    <Content>
      <Grid container>
        <Grid item xs={12}>
          <Breadcrumb>
            <BreadcrumbItem href='/'>Dashboard</BreadcrumbItem>
            <BreadcrumbItem>></BreadcrumbItem>
            <BreadcrumbItem href='/bridges'>Bridges</BreadcrumbItem>
            <BreadcrumbItem>></BreadcrumbItem>
            <BreadcrumbItem>{bridge && bridge.id}</BreadcrumbItem>
          </Breadcrumb>
        </Grid>
        <Grid item xs={12} md={12} xl={6}>
          <Grid container>
            <Grid item xs={9}>
              <Title>Edit Bridge</Title>
            </Grid>
            <Grid item xs={3}>
              <Grid container justify='flex-end'>
                <Grid item>
                  {bridge &&
                    <Button
                      variant='outlined'
                      color='primary'
                      component={ReactStaticLinkComponent}
                      to={`/bridges/${bridge.id}`}
                    >
                      Cancel
                    </Button>
                  }
                </Grid>
              </Grid>
            </Grid>
          </Grid>

          {onLoad(({bridge}) => (
            <PaddedCard>
              <BridgesForm
                actionText='Save Bridge'
                onSubmit={updateBridge}
                name={bridge.name}
                nameDisabled
                url={bridge.url}
                confirmations={bridge.confirmations}
                minimumContractPayment={bridge.minimumContractPayment}
                onSuccess={SuccessNotification}
                onError={ErrorMessage}
              />
            </PaddedCard>
          ))}
        </Grid>
      </Grid>
    </Content>
  )
}
)

Edit.propTypes = {
  bridge: PropTypes.object
}

const mapStateToProps = (state, ownProps) => {
  const bridge = bridgeSelector(state, ownProps.match.params.bridgeId)
  return {bridge}
}

export const ConnectedEdit = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({fetchBridgeSpec, updateBridge})
)(Edit)

export default ConnectedEdit
