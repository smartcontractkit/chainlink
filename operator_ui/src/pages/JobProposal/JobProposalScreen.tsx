import React from 'react'
import { useParams } from 'react-router-dom'
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter'
import { prism } from 'react-syntax-highlighter/dist/esm/styles/prism'

import { v2 } from 'api'
import { Resource, JobProposal } from 'core/store/models'
import Content from 'components/Content'
import { useErrorHandler } from 'hooks/useErrorHandler'
import { useLoadingPlaceholder } from 'hooks/useLoadingPlaceholder'

import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import CardHeader from '@material-ui/core/CardHeader'
import Grid from '@material-ui/core/Grid'

interface RouteParams {
  id: string
}

export const JobProposalScreen = () => {
  const { id } = useParams<RouteParams>()
  const [proposal, setProposal] = React.useState<Resource<JobProposal>>()
  const { error, ErrorComponent, setError } = useErrorHandler()
  const { LoadingPlaceholder } = useLoadingPlaceholder(!error && !proposal)

  React.useEffect(() => {
    v2.jobProposals
      .getJobProposal(id)
      .then((res) => {
        setProposal(res.data)
      })
      .catch(setError)
  }, [])

  return (
    <Content>
      <Grid container>
        <Grid item xs={12}>
          <ErrorComponent />
          <LoadingPlaceholder />

          <Card>
            <CardHeader title={`Job Proposal #${proposal?.id}`} />

            <CardContent>
              {proposal && (
                <div>
                  <SyntaxHighlighter language="toml" style={prism}>
                    {proposal.attributes.spec}
                  </SyntaxHighlighter>
                </div>
              )}
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Content>
  )
}
