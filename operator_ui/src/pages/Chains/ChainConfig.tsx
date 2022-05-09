import React from 'react'
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter'
import { prism } from 'react-syntax-highlighter/dist/esm/styles/prism'

import Content from 'components/Content'

import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import CardHeader from '@material-ui/core/CardHeader'
import Typography from '@material-ui/core/Typography'
import { Theme, createStyles, withStyles, WithStyles } from '@material-ui/core'
import { ChainResource } from './Show'

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
  chain: ChainResource
}

export const ChainConfig = withStyles(definitionStyles)(({ chain }: Props) => {
  const configOverrides = Object.fromEntries(
    Object.entries(chain.attributes.config).filter(
      ([_key, value]) => value !== null,
    ),
  )

  return (
    <Content>
      {chain && (
        <Card>
          <CardHeader title="Config Overrides" />
          <CardContent>
            <Typography style={{ margin: 0 }} variant="body1" component="pre">
              <SyntaxHighlighter language="json" style={prism}>
                {JSON.stringify(configOverrides, null, 2)}
              </SyntaxHighlighter>
            </Typography>
          </CardContent>
        </Card>
      )}
    </Content>
  )
})
