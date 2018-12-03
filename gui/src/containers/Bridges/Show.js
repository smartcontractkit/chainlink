import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'react-redux'
import { Button } from '@material-ui/core'
import Typography from '@material-ui/core/Typography'
import Grid from '@material-ui/core/Grid'
import Title from 'components/Title'
import PaddedCard from 'components/PaddedCard'
import Breadcrumb from 'components/Breadcrumb'
import BreadcrumbItem from 'components/BreadcrumbItem'
import ReactStaticLinkComponent from 'components/ReactStaticLinkComponent'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import { fetchBridgeSpec } from 'actions'
import bridgeSelector from 'selectors/bridge'
import Content from 'components/Content'
import { useHooks, useEffect } from 'use-react-hooks'

const renderLoading = () => (
  <div>Loading...</div>
)

const renderLoaded = props => (
  <PaddedCard>
    <Typography variant='subtitle1' color='textSecondary'>Name</Typography>
    <Typography variant='body1' color='inherit'>{props.bridge.name}</Typography>

    <Typography variant='subtitle1' color='textSecondary'>URL</Typography>
    <Typography variant='body1' color='inherit'>{props.bridge.url}</Typography>

    <Typography variant='subtitle1' color='textSecondary'>Confirmations</Typography>
    <Typography variant='body1' color='inherit'>{props.bridge.confirmations}</Typography>

    <Typography variant='subtitle1' color='textSecondary'>Minimum Contract Payment</Typography>
    <Typography variant='body1' color='inherit'>{props.bridge.minimumContractPayment}</Typography>

    <Typography variant='subtitle1' color='textSecondary'>Incoming Token</Typography>
    <Typography variant='body1' color='inherit'>{props.bridge.incomingToken}</Typography>

    <Typography variant='subtitle1' color='textSecondary'>Outgoing Token</Typography>
    <Typography variant='body1' color='inherit'>{props.bridge.outgoingToken}</Typography>
  </PaddedCard>
)

const renderDetails = props => props.bridge ? renderLoaded(props) : renderLoading(props)

export const Show = useHooks(props => {
  useEffect(() => { props.fetchBridgeSpec(props.match.params.bridgeId) }, [])
  return (
    <Content>
      <Grid container>
        <Grid item xs={12}>
          <Breadcrumb>
            <BreadcrumbItem href='/'>Dashboard</BreadcrumbItem>
            <BreadcrumbItem>></BreadcrumbItem>
            <BreadcrumbItem href='/bridges'>Bridges</BreadcrumbItem>
            <BreadcrumbItem>></BreadcrumbItem>
            <BreadcrumbItem>{props.bridge && props.bridge.id}</BreadcrumbItem>
          </Breadcrumb>
        </Grid>
        <Grid item xs={12} md={12} xl={6}>
          <Grid container>
            <Grid item xs={9}>
              <Title>Bridge Info</Title>
            </Grid>
            <Grid item xs={3}>
              <Grid container justify='flex-end'>
                <Grid item>
                  {props.bridge &&
                    <Button
                      variant='outlined'
                      color='primary'
                      component={ReactStaticLinkComponent}
                      to={`/bridges/${props.bridge.id}/edit`}
                    >
                      Edit
                    </Button>
                  }
                </Grid>
              </Grid>
            </Grid>
          </Grid>

          {renderDetails(props)}
        </Grid>
      </Grid>
    </Content>
  )
}
)

Show.propTypes = {
  bridge: PropTypes.object
}

const mapStateToProps = (state, ownProps) => ({
  bridge: bridgeSelector(state, ownProps.match.params.bridgeId)
})

export const ConnectedShow = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({fetchBridgeSpec})
)(Show)

export default ConnectedShow
