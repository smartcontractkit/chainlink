import React from 'react'
import { Drawer, Button } from 'antd'
import moment from 'moment'

const NodeDetailsModal = ({
  onClose,
  visible,
  data,
  pendingAnswerId,
  jobId,
  options
}) => {
  return (
    <Drawer
      placement="right"
      closable={true}
      width={600}
      onClose={onClose}
      visible={visible}>
      <NodeDetailsContent
        data={data}
        pendingAnswerId={pendingAnswerId}
        jobId={jobId}
        options={options}
      />
    </Drawer>
  )
}

const NodeDetailsContent = ({ data = {}, pendingAnswerId, jobId, options }) => {
  const dateResonse = timestamp =>
    moment.unix(timestamp).format('DD/MM/YY h:mm:ss A')

  return (
    <div className="network-graph__node-details">
      <h2>{data.name}</h2>
      <hr className="hr" />
      <div className="network-graph__node-details__item">
        <div className="network-graph__node-details__item--label">
          Latest answer
        </div>
        <h3 className="network-graph__node-details__item--value">
          {options.valuePrefix || ''}{' '}
          {(data.state && data.state.responseFormatted) || '-'}
        </h3>
      </div>

      <div className="network-graph__node-details__item">
        <div className="network-graph__node-details__item--label">
          Response date
        </div>
        <h3 className="network-graph__node-details__item--value">
          {data.state && dateResonse(data.state && data.state.meta.timestamp)}
        </h3>
      </div>

      <hr className="hr" />

      <div>
        <h4>Find out more in:</h4>
        <Button
          style={{ marginRight: 10 }}
          ghost
          type="primary"
          disabled={!jobId}>
          <a
            target="_BLANK"
            rel="noopener noreferrer"
            href={`https://explorer.chain.link/job-runs?search=${jobId}`}>
            Chainlink Explorer
          </a>
        </Button>
        <Button style={{ marginRight: 10 }} ghost type="primary">
          <a
            target="_BLANK"
            rel="noopener noreferrer"
            href={`https://market.link/search/nodes?&name=${data.name}`}>
            Market.link
          </a>
        </Button>
        <Button ghost type="primary">
          <a
            target="_BLANK"
            rel="noopener noreferrer"
            href={`https://etherscan.io/address/${data.address}`}>
            Etherscan
          </a>
        </Button>
      </div>
    </div>
  )
}

export default NodeDetailsModal
