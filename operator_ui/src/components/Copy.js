import React from 'react'
import Button from 'components/Button'
import { CopyToClipboard } from 'react-copy-to-clipboard'

const Copy = ({ data, buttonText, ...props }) => (
  <CopyToClipboard text={data}>
    <Button {...props}>{buttonText}</Button>
  </CopyToClipboard>
)

export default Copy
