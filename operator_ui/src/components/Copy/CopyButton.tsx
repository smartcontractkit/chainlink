import React from 'react'

import { CopyToClipboard } from 'react-copy-to-clipboard'

import Button from '../Button'
import Tooltip from '@material-ui/core/Tooltip'

interface Props {
  data: string
  title: string
}

export const CopyButton: React.FC<Props> = ({ data, title }) => {
  const [copied, setCopied] = React.useState(false)

  return (
    <CopyToClipboard text={data} onCopy={() => setCopied(true)}>
      <Tooltip title="Copied!" open={copied} onClose={() => setCopied(false)}>
        <Button variant="secondary" onMouseLeave={() => setCopied(false)}>
          {title}
        </Button>
      </Tooltip>
    </CopyToClipboard>
  )
}
