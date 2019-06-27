import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'react-redux'
import Typography from '@material-ui/core/Typography'
import Grid from '@material-ui/core/Grid'
import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import Divider from '@material-ui/core/Divider'
import BaseLink from 'components/BaseLink'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import { fetchBridgeSpec } from 'actions'
import bridgeSelector from 'selectors/bridge'
import Content from 'components/Content'
import Button from 'components/Button'
import { useHooks, useEffect } from 'use-react-hooks'

const renderLoading = () => <div>Loading...</div>

const renderLoaded = props => (
  <CardContent>
    <Typography variant="subtitle1" color="textSecondary">
      Name
    </Typography>
    <Typography variant="body1" color="inherit">
      {props.bridge.name}
    </Typography>

    <Typography variant="subtitle1" color="textSecondary">
      URL
    </Typography>
    <Typography variant="body1" color="inherit">
      {props.bridge.url}
    </Typography>

    <Typography variant="subtitle1" color="textSecondary">
      Confirmations
    </Typography>
    <Typography variant="body1" color="inherit">
      {props.bridge.confirmations}
    </Typography>

    <Typography variant="subtitle1" color="textSecondary">
      Minimum Contract Payment
    </Typography>
    <Typography variant="body1" color="inherit">
      {props.bridge.minimumContractPayment}
    </Typography>

    <Typography variant="subtitle1" color="textSecondary">
      Outgoing Token
    </Typography>
    <Typography variant="body1" color="inherit">
      {props.bridge.outgoingToken}
    </Typography>
  </CardContent>
)

const renderDetails = props =>
  props.bridge ? renderLoaded(props) : renderLoading(props)

export const Show = useHooks(props => {
  useEffect(() => {
    document.title = 'Show Bridge'
    props.fetchBridgeSpec(props.match.params.bridgeId)
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
                      {props.bridge && (
                        <Button
                          variant="secondary"
                          component={BaseLink}
                          to={`/bridges/${props.bridge.id}/edit`}
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

            {renderDetails(props)}
          </Card>
        </Grid>
      </Grid>
    </Content>
  )
})

Show.propTypes = {
  bridge: PropTypes.object
}

const mapStateToProps = (state, ownProps) => ({
  bridge: bridgeSelector(state, ownProps.match.params.bridgeId)
})

export const ConnectedShow = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({ fetchBridgeSpec })
)(Show)

export default ConnectedShow
