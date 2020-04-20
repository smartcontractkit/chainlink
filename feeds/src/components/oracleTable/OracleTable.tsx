import React, { useEffect, useState } from 'react'
import { Table, Icon } from 'antd'
import { ColumnProps } from 'antd/lib/table'
import { humanizeUnixTimestamp } from 'utils'
import { connect } from 'react-redux'
import {
  aggregatorSelectors,
  aggregatorOperations,
} from 'state/ducks/aggregator'
import { AppState } from 'state'

interface StateProps {
  fetchEthGasPrice: any
  ethGasPrice: any
  latestOraclesState: any
}

interface DispatchProps {
  fetchEthGasPrice: any
}

export interface Props extends StateProps, DispatchProps {}

const OracleTable: React.FC<Props> = ({
  fetchEthGasPrice,
  ethGasPrice,
  latestOraclesState,
}) => {
  const [data, setData] = useState<any | undefined>()
  const [gasPrice, setGasPrice] = useState<any | undefined>()

  useEffect(() => {
    fetchEthGasPrice()
  }, [fetchEthGasPrice])

  useEffect(() => {
    if (ethGasPrice) {
      setGasPrice(ethGasPrice.fast / 10)
    }
  }, [ethGasPrice])

  useEffect(() => {
    setData(latestOraclesState)
  }, [latestOraclesState])

  const columns: ColumnProps<any>[] = [
    {
      title: 'Oracle',
      dataIndex: 'name',
      key: 'name',
      sorter: (a: any, b: any): number => a.name.localeCompare(b && b.name),
    },

    {
      title: 'Answer',
      dataIndex: 'answerFormatted',
      key: 'answerFormatted',
      sorter: (a: any, b: any): number => {
        if (!a.answerFormatted || !b.answerFormatted) return 0
        return a.answerFormatted - b.answerFormatted
      },
    },
    {
      title: 'Gas Price (Gwei)',
      dataIndex: 'meta.gasPrice',
      key: 'gas',
      sorter: (a: any, b: any): number => {
        if (!a.meta || !b.meta) return 0
        return a.meta.gasPrice - b.meta.gasPrice
      },
      defaultSortOrder: 'descend' as 'descend',
    },
    {
      title: 'Date',
      dataIndex: 'meta.timestamp',
      key: 'timestamp',
      sorter: (a: any, b: any): number => {
        if (!a.meta || !b.meta) return 0
        return a.meta.timestamp - b.meta.timestamp
      },
      render: (timestamp: number) =>
        timestamp && humanizeUnixTimestamp(timestamp, 'LLL'),
    },
  ]

  return (
    <div className="oracle-table">
      <h2 className="oracle-table-header">Oracles data</h2>
      <div className="gas-price-info">
        <h4>
          <Icon type="check-circle" /> Recommended gas price:{' '}
          <b>{gasPrice} Gwei</b>
        </h4>
      </div>
      <Table
        dataSource={data}
        columns={columns}
        pagination={false}
        size={'middle'}
        locale={{ emptyText: <Icon type="loading" /> }}
        rowKey={record => record.id}
      />
    </div>
  )
}

const mapStateToProps = (state: AppState) => ({
  ethGasPrice: state.aggregator.ethGasPrice,
  latestOraclesState: aggregatorSelectors.latestOraclesState(state),
})

const mapDispatchToProps = {
  fetchEthGasPrice: aggregatorOperations.fetchEthGasPrice,
}

export default connect(mapStateToProps, mapDispatchToProps)(OracleTable)
