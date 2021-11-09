import React from 'react'
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter'
import { prism } from 'react-syntax-highlighter/dist/esm/styles/prism'

import Content from 'components/Content'
import { JobData } from './sharedTypes'

import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import CardHeader from '@material-ui/core/CardHeader'
import Typography from '@material-ui/core/Typography'
import { Theme, createStyles, withStyles, WithStyles } from '@material-ui/core'

const definitionStyles = (theme: Theme) =>
  createStyles({
    definitionTitle: {
      marginTop: theme.spacing.unit * 2,
      marginBottom: theme.spacing.unit * 2,
    },
    divider: {
      marginTop: theme.spacing.unit,
      marginBottom: theme.spacing.unit * 3,
    },
  })

interface Props extends WithStyles<typeof definitionStyles> {
  error: unknown
  ErrorComponent: React.FC
  LoadingPlaceholder: React.FC
  job?: JobData['job']
  envAttributesDefinition?: string
}

export const JobDefinition = withStyles(definitionStyles)(
  ({
    error,
    ErrorComponent,
    LoadingPlaceholder,
    job,
    envAttributesDefinition,
  }: Props) => {
    React.useEffect(() => {
      document.title = job?.name
        ? `${job.name} | Job definition`
        : 'Job definition'
    }, [job])

    return (
      <Content>
        <ErrorComponent />
        <LoadingPlaceholder />

        {!error && job && (
          <Card>
            <CardHeader title="Definition" />
            <CardContent>
              <Typography style={{ margin: 0 }} variant="body1" component="pre">
                <SyntaxHighlighter
                  language="toml"
                  style={prism}
                  data-testid="definition"
                >
                  {job.definition}
                </SyntaxHighlighter>
              </Typography>

              {envAttributesDefinition?.trim() && (
                <>
                  <Typography variant="h5" color="secondary">
                    Job attributes set by environment variables
                  </Typography>

                  <Typography
                    style={{ margin: 0 }}
                    variant="body1"
                    component="pre"
                  >
                    <SyntaxHighlighter language="toml" style={prism}>
                      {envAttributesDefinition}
                    </SyntaxHighlighter>
                  </Typography>
                </>
              )}
            </CardContent>
          </Card>
        )}
      </Content>
    )
  },
)
