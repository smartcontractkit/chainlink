import React, { useEffect, useState } from 'react'
import { Drawer, Button } from 'antd'
import { etherscanAddress, humanizeUnixTimestamp } from 'utils'
import { connect } from 'react-redux'
import { networkGraphOperations } from 'state/ducks/networkGraph'

const NodeDetailsModal = ({ options, fetchJobId, sideDrawer, setDrawer }) => {
  const [jobId, setJobId] = useState()

  useEffect(() => {
    async function fetchJobIdEffect() {
      if (sideDrawer && sideDrawer.type === 'oracle') {
        const jobIdResponse = await fetchJobId(sideDrawer.address)
        setJobId(jobIdResponse)
      }
    }
    fetchJobIdEffect()
  }, [sideDrawer, fetchJobId])

  return (
    <Drawer
      placement="right"
      closable={true}
      width={600}
      onClose={() => setDrawer(null)}
      visible={!!sideDrawer}
    >
      {sideDrawer && sideDrawer.type === 'oracle' ? (
        <NodeDetailsContent data={sideDrawer} jobId={jobId} options={options} />
      ) : (
        <ContractDetailsContent data={sideDrawer} options={options} />
      )}
    </Drawer>
  )
}

const NodeDetailsContent = ({ data = {}, jobId, options }) => {
  if (!data) return null

  return (
    <div className="network-graph__node-details">
      <h2>{data.name}</h2>
      <hr className="hr" />
      <div className="network-graph__node-details__item">
        <div className="network-graph__node-details__item--label">
          Latest answer
        </div>
        <h3 className="network-graph__node-details__item--value">
          {options.valuePrefix || ''} {data.responseFormatted || '...'}
        </h3>
      </div>

      <div className="network-graph__node-details__item">
        <div className="network-graph__node-details__item--label">
          Response date
        </div>
        <h3 className="network-graph__node-details__item--value">
          {data.meta && humanizeUnixTimestamp(data.meta.timestamp, 'LLL')}
        </h3>
      </div>

      <hr className="hr" />

      <div>
        <h4>Find out more in:</h4>
        {options.network === 'mainnet' && (
          <>
            <Button
              style={{ marginRight: 10 }}
              ghost
              type="primary"
              disabled={!jobId}
            >
              <a
                target="_BLANK"
                rel="noopener noreferrer"
                href={`https://explorer.chain.link/job-runs?search=${jobId}`}
              >
                Chainlink Explorer
              </a>
            </Button>
            <Button style={{ marginRight: 10 }} ghost type="primary">
              <a
                target="_BLANK"
                rel="noopener noreferrer"
                href={`https://market.link/search/nodes?&name=${data.name}`}
              >
                Market.link
              </a>
            </Button>{' '}
          </>
        )}
        <Button ghost type="primary">
          <a
            target="_BLANK"
            rel="noopener noreferrer"
            href={etherscanAddress(options.network, data.address)}
          >
            Etherscan
          </a>
        </Button>
      </div>
    </div>
  )
}

const ContractDetailsContent = ({ data = {}, options }) => {
  if (!data) return null

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
          {options.valuePrefix || ''} {data.currentAnswer || '...'}
        </h3>
      </div>

      <div className="network-graph__contract-details__item">
        <div className="network-graph__contract-details__item--label">
          Response date
        </div>
        <h3 className="network-graph__contract-details__item--value">
          {humanizeUnixTimestamp(data.updateHeight, 'LLL')}
        </h3>
      </div>

      <hr className="hr" />

      <div>
        <h4>Find out more in:</h4>
        {options.network === 'mainnet' && (
          <Button style={{ marginRight: 10 }} ghost type="primary">
            <a
              target="_BLANK"
              rel="noopener noreferrer"
              href={`https://explorer.chain.link/job-runs?search=${options.contractAddress}`}
            >
              Chainlink Explorer
            </a>
          </Button>
        )}
        <Button ghost type="primary">
          <a
            target="_BLANK"
            rel="noopener noreferrer"
            href={etherscanAddress(options.network, options.contractAddress)}
          >
            Etherscan
          </a>
        </Button>
      </div>
    </div>
  )
}

const mapStateToProps = state => ({
  tooltip: state.networkGraph.tooltip,
  sideDrawer: state.networkGraph.drawer,
})

const mapDispatchToProps = {
  setTooltip: networkGraphOperations.setTooltip,
  setDrawer: networkGraphOperations.setDrawer,
}

export default connect(mapStateToProps, mapDispatchToProps)(NodeDetailsModal)
