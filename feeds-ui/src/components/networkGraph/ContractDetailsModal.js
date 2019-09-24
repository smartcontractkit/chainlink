import React from 'react'
import { Drawer, Button } from 'antd'
import moment from 'moment'

const ContractDetailsModal = ({
  onClose,
  visible,
  data,
  options,
  updateHeight
}) => {
  return (
    <Drawer
      placement="right"
      closable={true}
      width={600}
      onClose={onClose}
      visible={visible}>
      <ContractDetailsContent
        data={data}
        options={options}
        updateHeight={updateHeight}
      />
    </Drawer>
  )
}

const ContractDetailsContent = ({
  data = {},
  pendingAnswerId,
  jobId,
  options,
  updateHeight
}) => {
  const dateResonse = timestamp =>
    moment.unix(timestamp).format('DD/MM/YY h:mm:ss A')
  return (
    <div className="network-graph__contract-details">
      <h2>Aggregation Contract</h2>

      <hr className="hr" />

      <div className="network-graph__contract-details__item">
        <div className="network-graph__contract-details__item--label">Type</div>
        <h3 className="network-graph__contract-details__item--value">
          {options.name}
        </h3>
      </div>

      <div className="network-graph__contract-details__item">
        <div className="network-graph__contract-details__item--label">
          Latest answer
        </div>
        <h3 className="network-graph__contract-details__item--value">
          {options.valuePrefix || ''}{' '}
          {(data.state && data.state.currentAnswer) || '-'}
        </h3>
      </div>

      <div className="network-graph__contract-details__item">
        <div className="network-graph__contract-details__item--label">
          Response date
        </div>
        <h3 className="network-graph__contract-details__item--value">
          {updateHeight && dateResonse(updateHeight.timestamp)}
        </h3>
      </div>

      <hr className="hr" />

      <div>
        <h4>Find out more in:</h4>
        <Button style={{ marginRight: 10 }} ghost type="primary">
          <a
            target="_BLANK"
            rel="noopener noreferrer"
            href={`https://explorer.chain.link/job-runs?search=${options.contractAddress}`}>
            Chainlink Explorer
          </a>
        </Button>
        <Button ghost type="primary">
          <a
            target="_BLANK"
            rel="noopener noreferrer"
            href={`https://etherscan.io/address/${options.contractAddress}`}>
            Etherscan
          </a>
        </Button>
      </div>
    </div>
  )
}

export default ContractDetailsModal
