import React from 'react'

import { useQuery } from '@apollo/client'

import { connect } from 'react-redux'
import Grid from '@material-ui/core/Grid'
import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import Divider from '@material-ui/core/Divider'
import Typography from '@material-ui/core/Typography'
import Form from 'components/Bridges/Form'
import ErrorMessage from 'components/Notifications/DefaultError'
import Content from 'components/Content'
import { createBridge } from 'actionCreators'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import BaseLink from '../../components/BaseLink'
import { BridgeTypeAuthentication } from 'core/store/models'
import { BRIDGES_QUERY } from 'src/screens/Bridges/BridgesScreen'

interface SuccessResponse {
  type: string
  id: string
  attributes: BridgeTypeAuthentication
}

const SuccessNotification = (response: SuccessResponse) => {
  return (
    <React.Fragment>
      Successfully created bridge{' '}
      <BaseLink href={`/bridges/${response.id}`}>{response.id}</BaseLink> with
      incoming access token: {response.attributes.incomingToken}
    </React.Fragment>
  )
}

interface Props {
  createBridge: () => Promise<any>
}

const New: React.FC<Props> = (props) => {
  document.title = 'New Bridge'

  // This is a hacky fix to refetch the page data after an edit so changes
  // appear on the table. Once this gets moved to GQL, we can use the
  // refetchQueries on the mutation.
  const { refetch } = useQuery<FetchBridges, FetchBridgesVariables>(
    BRIDGES_QUERY,
    {
      variables: { offset: 0, limit: 10 },
    },
  )

  return (
    <Content>
      <Grid container spacing={40}>
        <Grid item xs={12} md={11} lg={9}>
          <Card>
            <CardContent>
              <Typography variant="h5" color="secondary">
                New Bridge
              </Typography>
            </CardContent>

            <Divider />

            <CardContent>
              <Form
                actionText="Create Bridge"
                onSubmit={props.createBridge}
                onSuccess={SuccessNotification}
                onError={ErrorMessage}
                refetchGQL={refetch}
              />
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Content>
  )
}

export const ConnectedNew = connect(
  null,
  matchRouteAndMapDispatchToProps({ createBridge }),
)(New)

export default ConnectedNew
