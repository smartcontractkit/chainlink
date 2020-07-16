import { Button, Drawer } from 'antd'
import { FeedConfig } from 'config'
import React, { useEffect, useState } from 'react'
import { connect } from 'react-redux'
import { AppState } from 'state'
import { networkGraphOperations } from 'state/ducks/networkGraph'
import { etherscanAddress, humanizeUnixTimestamp, Networks } from 'utils'

interface OwnProps {
  config: FeedConfig
  fetchJobId?: any
}

interface StateProps {
  sideDrawer: any
}

interface DispatchProps {
  setDrawer: any
}

interface Props extends OwnProps, StateProps, DispatchProps {}

const NodeDetailsModal: React.FC<Props> = ({
  config,
  fetchJobId,
  sideDrawer,
  setDrawer,
}) => {
  const [jobId, setJobId] = useState()

  useEffect(() => {
    async function fetchJobIdEffect() {
      if (sideDrawer && sideDrawer.type === 'oracle') {
        const jobIdResponse = await fetchJobId(sideDrawer.address)
        setJobId(jobIdResponse)
      }
    }
    if (fetchJobId) {
      fetchJobIdEffect()
    }
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
        <NodeDetailsContent data={sideDrawer} jobId={jobId} config={config} />
      ) : (
        <ContractDetailsContent data={sideDrawer} config={config} />
      )}
    </Drawer>
  )
}

interface NodeDetailsContentProps {
  data: any
  jobId?: string
  config: FeedConfig
}

const NodeDetailsContent: React.FC<NodeDetailsContentProps> = ({
  data = {},
  jobId,
  config,
}) => {
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
          {config.valuePrefix || ''} {data.answerFormatted || '...'}
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
        {config.networkId === Networks.MAINNET && (
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
                href={`https://market.link/search/nodes?search=${data.address}`}
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
            href={etherscanAddress(config.networkId, data.address)}
          >
            Etherscan
          </a>
        </Button>
      </div>
    </div>
  )
}

interface ContractDetailsContentProps {
  data: any
  config: FeedConfig
}

const ContractDetailsContent: React.FC<ContractDetailsContentProps> = ({
  data = {},
  config,
}) => {
  if (!data) return null

  return (
    <div className="network-graph__contract-details">
      <h2>Aggregation Contract</h2>

      <hr className="hr" />

      <div className="network-graph__contract-details__item">
        <div className="network-graph__contract-details__item--label">Type</div>
        <h3 className="network-graph__contract-details__item--value">
          {config.name}
        </h3>
      </div>

      <div className="network-graph__contract-details__item">
        <div className="network-graph__contract-details__item--label">
          Latest answer
        </div>
        <h3 className="network-graph__contract-details__item--value">
          {config.valuePrefix || ''} {data.latestAnswer || '...'}
        </h3>
      </div>

      <div className="network-graph__contract-details__item">
        <div className="network-graph__contract-details__item--label">
          Answer date
        </div>
        <h3 className="network-graph__contract-details__item--value">
          {data.latestAnswerTimestamp
            ? humanizeUnixTimestamp(data.latestAnswerTimestamp, 'LLL')
            : '...'}
        </h3>
      </div>

      <hr className="hr" />

      <div>
        <h4>Find out more in:</h4>
        {config.networkId === Networks.MAINNET && (
          <Button style={{ marginRight: 10 }} ghost type="primary">
            <a
              target="_BLANK"
              rel="noopener noreferrer"
              href={`https://explorer.chain.link/job-runs?search=${config.contractAddress}`}
            >
              Chainlink Explorer
            </a>
          </Button>
        )}
        <Button ghost type="primary">
          <a
            target="_BLANK"
            rel="noopener noreferrer"
            href={etherscanAddress(config.networkId, config.contractAddress)}
          >
            Etherscan
          </a>
        </Button>
      </div>
    </div>
  )
}

const mapStateToProps = (state: AppState) => ({
  sideDrawer: state.networkGraph.drawer,
})

const mapDispatchToProps = {
  setDrawer: networkGraphOperations.setDrawer,
}

export default connect(mapStateToProps, mapDispatchToProps)(NodeDetailsModal)
