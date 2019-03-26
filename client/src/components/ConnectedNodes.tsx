import React, { Component } from 'react'
import io from 'socket.io-client'

interface IProps {
  className: string
}

interface IState {
  count?: number
}

const showCount = (count?: number): string => {
  if (count === undefined) {
    return '...'
  }
  return count.toString()
}

class ConnectedNodes extends Component<IProps, IState> {
  constructor(props: IProps) {
    super(props)
    this.state = {}
  }

  componentDidMount() {
    const socket = io('/', { path: '/client' })

    socket.on('clnodeCount', (count: number) => {
      this.setState({ count: count })
    })
  }

  render() {
    return (
      <div className={this.props.className}>
        Connected Nodes: {showCount(this.state.count)}
      </div>
    )
  }
}

export default ConnectedNodes
