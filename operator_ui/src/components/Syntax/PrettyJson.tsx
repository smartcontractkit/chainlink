import React from 'react'

import jsonPrettyHtml from 'json-pretty-html'
import './PrettyJson.css'

interface Props {
  object: object
}

// TODO - Look to standardise syntax highlighting so we can reduce dependencies
export const PrettyJson: React.FC<Props> = ({ object }) => (
  <div
    dangerouslySetInnerHTML={{ __html: jsonPrettyHtml(object) }}
    data-testid="pretty-json"
  />
)
