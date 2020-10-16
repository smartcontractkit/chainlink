import React from 'react'
import { connect } from 'react-redux'
import PropTypes from 'prop-types'
import Grid from '@material-ui/core/Grid'
import Title from 'components/Title'
import BridgeList from 'components/Bridges/List'
import BaseLink from 'components/BaseLink'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import bridgesSelector from 'selectors/bridges'
import { fetchBridges } from 'actionCreators'
import Content from 'components/Content'
import Button from 'components/Button'

export const Index = (props) => {
  document.title = 'Bridges'
  const { bridges, count, pageSize, fetchBridges, history, match } = props
  return (
    <Content>
      <Grid container spacing={8}>
        <Grid item xs={9}>
          <Title>Bridges</Title>
        </Grid>
        <Grid item xs={3}>
          <Grid container justify="flex-end">
            <Grid item>
              <Button
                variant="secondary"
                component={BaseLink}
                href={'/bridges/new'}
              >
                New Bridge
              </Button>
            </Grid>
          </Grid>
        </Grid>
      </Grid>
      <Grid container spacing={40}>
        <Grid item xs={12}>
          <BridgeList
            bridges={bridges}
            bridgeCount={count}
            pageSize={pageSize}
            fetchBridges={fetchBridges}
            history={history}
            match={match}
          />
        </Grid>
      </Grid>
    </Content>
  )
}

Index.propTypes = {
  count: PropTypes.number.isRequired,
  bridges: PropTypes.array.isRequired,
  pageSize: PropTypes.number,
}

Index.defaultProps = {
  pageSize: 10,
}

const mapStateToProps = (state) => {
  const bridges = bridgesSelector(state)
  const count = state.bridges.count

  return { bridges, count }
}

export const ConnectedIndex = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({ fetchBridges }),
)(Index)

export default ConnectedIndex
