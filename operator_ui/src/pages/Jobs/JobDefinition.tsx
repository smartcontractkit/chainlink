import React from 'react'
import {
  createStyles,
  CardContent,
  Card,
  Theme,
  Typography,
  withStyles,
  WithStyles,
} from '@material-ui/core'
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter'
import { prism } from 'react-syntax-highlighter/dist/esm/styles/prism'

import { CardTitle } from 'components/CardTitle'
import Content from 'components/Content'
import PrettyJson from 'components/PrettyJson'
import { JobData } from './sharedTypes'

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

type Props = {
  error: unknown
  ErrorComponent: React.FC
  LoadingPlaceholder: React.FC
  job?: JobData['job']
} & WithStyles<typeof definitionStyles>

export const JobDefinition = withStyles(definitionStyles)(
  ({ error, ErrorComponent, LoadingPlaceholder, job }: Props) => {
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
            <CardTitle divider>Definition</CardTitle>
            <CardContent>
              {job.type === 'Direct request' && (
                <PrettyJson object={JSON.parse(job.definition)} />
              )}
              {job.type === 'v2' && (
                <Typography
                  style={{ margin: 0 }}
                  variant="body1"
                  component="pre"
                >
                  <SyntaxHighlighter language="toml" style={prism}>
                    {job.definition}
                  </SyntaxHighlighter>
                </Typography>
              )}
            </CardContent>
          </Card>
        )}
      </Content>
    )
  },
)
