import React from 'react'
import PropTypes from 'prop-types'
import jsonPrettyHtml from 'json-pretty-html'
import './PrettyJson.css'

const PrettyJson = ({ object }) => (
  <div dangerouslySetInnerHTML={{ __html: jsonPrettyHtml(object) }} />
)

PrettyJson.propTypes = {
  object: PropTypes.object.isRequired,
}

export default PrettyJson
