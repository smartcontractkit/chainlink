import React from 'react'
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter'
import { prism } from 'react-syntax-highlighter/dist/esm/styles/prism'

import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import CardHeader from '@material-ui/core/CardHeader'
import Typography from '@material-ui/core/Typography'

import { CopyButton } from 'src/components/Copy/CopyButton'
import { generateJobDefinition } from './generateJobDefinition'

interface Props {
  job: JobPayload_Fields
}

export const TabDefinition = ({ job }: Props) => {
  const { definition, envDefinition } = generateJobDefinition(job)

  return (
    <Card>
      <CardHeader
        title="Definition"
        action={<CopyButton title="Copy" data={definition} />}
      />
      <CardContent>
        <Typography style={{ margin: 0 }} variant="body1" component="pre">
          <SyntaxHighlighter
            language="toml"
            style={prism}
            data-testid="definition"
          >
            {definition}
          </SyntaxHighlighter>
        </Typography>

        {envDefinition.trim() && (
          <>
            <Typography variant="h5" color="secondary">
              Job attributes set by environment variables
            </Typography>

            <Typography style={{ margin: 0 }} variant="body1" component="pre">
              <SyntaxHighlighter language="toml" style={prism}>
                {envDefinition}
              </SyntaxHighlighter>
            </Typography>
          </>
        )}
      </CardContent>
    </Card>
  )
}
