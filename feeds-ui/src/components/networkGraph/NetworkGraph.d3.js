import * as d3 from 'd3'
import _ from 'lodash'
import moment from 'moment'

function hasPrice(d) {
  return d.state && d.state.responseFormatted
}

export default class NetworkGraph {
  width = 1200
  height = 600
  svg
  nodes
  links
  options = {}
  dragging = false
  pendingAnswerId

  constructor(options) {
    this.options = options
  }

  build() {
    this.svg = d3
      .select('.network-graph__svg')
      .attr('viewBox', `0 0 ${this.width} ${this.height}`)
      .attr('preserveAspectRatio', 'xMidYMid meet')

    this.nodes = this.svg.select('.network-graph__nodes').selectAll('g')
    this.links = this.svg.select('.network-graph__links').selectAll('line')

    this.oracleTooltip = d3
      .select('.network-graph')
      .select('.network-graph__tooltip--oracle')
      .style('opacity', 0)

    this.contractTooltip = d3
      .select('.network-graph')
      .select('.network-graph__tooltip--contract')
      .style('opacity', 0)
  }

  updateNodes(nodes) {
    const updateData = this.svg
      .select('.network-graph__nodes')
      .selectAll('g.network-graph__node-group')
      .data(nodes, d => {
        d.radius = d.type === 'contract' ? 60 : 10
        return d.id
      })

    updateData.exit().remove()

    const nodesEnter = updateData
      .enter()
      .append('g')
      .attr('class', 'network-graph__node-group')
      .on('click', d => {
        if (this.onNodeClick && d.type === 'oracle') {
          this.onNodeClick(d)
        }

        if (this.onContractClick && d.type === 'contract') {
          this.onContractClick(d)
        }
      })

    this.nodes = updateData.merge(nodesEnter)

    this.nodes.attr('transform', d => {
      d.x = this.width / 2
      d.y = this.height / 2
      return 'translate(' + d.x + ',' + d.y + ')'
    })

    const oracles = this.nodes.filter(d => {
      return d.type !== 'contract'
    })

    oracles.attr('transform', (d, i) => {
      const size = oracles.size()
      const width = this.width
      const height = this.height

      const angle = (i / (size / 2)) * Math.PI
      const x = 280 * Math.cos(angle) + width / 2
      const y = 280 * Math.sin(angle) + height / 2
      d.x = x
      d.y = y
      return 'translate(' + d.x + ',' + d.y + ')'
    })

    nodesEnter
      .append('circle')
      .attr('class', d => `network-graph__node__${d.type}`)
      .attr('r', d => d.radius)
      .attr('fill', '#f1f1f1')
      .attr('stroke', '#8e8e8e')

    nodesEnter
      .selectAll('.network-graph__node__oracle')
      .on('mouseover', d => this.setOracleTooltip(d))
      .on('mouseout', () => {
        this.oracleTooltip.style('opacity', 0)
      })

    nodesEnter
      .selectAll('.network-graph__node__contract')
      .on('mouseover', d => this.setContractTooltip(d))
      .on('mouseout', () => {
        this.contractTooltip.style('opacity', 0)
      })

    const label = nodesEnter
      .filter(d => {
        return d.type !== 'contract'
      })
      .append('g')
      .attr('class', 'network-graph__node__oracle-label')

    label
      .append('text')
      .attr('class', 'network-graph__node__oracle-label--name')
      .attr('x', '20')
      .attr('y', '-5')
      .attr('opacity', 0)
      .transition()
      .duration(600)
      .attr('opacity', 1)
      .text(d => {
        return d.name
      })

    label
      .append('text')
      .attr('class', 'network-graph__node__oracle-label--price')
      .style('opacity', 0.6)
      .attr('x', '20')
      .attr('y', '10')

    const contract = nodesEnter.filter(d => {
      return d.type === 'contract'
    })

    const contractLabel = contract
      .append('g')
      .attr('class', 'network-graph__node__contract-label')

    contractLabel
      .append('text')
      .style('opacity', 0)
      .transition()
      .style('opacity', 1)
      .attr('class', 'network-graph__node__contract-label--price')
      .attr('y', '5')
      .attr('text-anchor', 'middle')
      .text(c => {
        return 'Loading...'
      })

    this.nodes.select('.network-graph__node__contract').attr('class', () => {
      if (this.nodes.data().length > 1) {
        return 'network-graph__node__contract'
      }
      return 'network-graph__node__contract wait'
    })

    this.links = this.links.data(nodes, d => {
      return d
    })

    this.links.exit().remove()

    this.links = this.links
      .enter()
      .append('line')
      .attr('class', 'network-graph__line--wait')
      .attr('x1', l => {
        return this.width / 2
      })
      .attr('y1', l => {
        return this.height / 2
      })
      .attr('x2', l => {
        return l.x
      })
      .attr('y2', l => {
        return l.y
      })
  }

  setOracleTooltip(d) {
    if (this.dragging) {
      return this.oracleTooltip.style('opacity', 0)
    }

    this.oracleTooltip
      .style('opacity', 1)
      .style('left', d.x + 15 + 'px')
      .style('top', d.y + 'px')

    this.oracleTooltip.select('.price').text('')
    this.oracleTooltip.select('.date').text('')
    this.oracleTooltip.select('.block').text('')

    this.oracleTooltip.select('.name').text(() => {
      return d.name
    })

    if (d.state) {
      this.oracleTooltip.select('.price').text(() => {
        return `${this.options.valuePrefix || ''} ${d.state.responseFormatted}`
      })

      this.oracleTooltip.select('.date').text(() => {
        return `Date: ${moment
          .unix(d.state.meta.timestamp)
          .format('DD/MM/YY h:mm:ss A')}`
      })
      this.oracleTooltip.select('.block').text(() => {
        return `Block: ${d.state.meta.blockNumber}`
      })
    }
  }

  setContractTooltip(d) {
    if (this.dragging) {
      return this.contractTooltip.style('opacity', 0)
    }

    this.contractTooltip
      .style('opacity', 1)
      .style('left', d.x + 15 + 'px')
      .style('top', d.y + 15 + 'px')

    this.contractTooltip.select('.price').text('')
    this.contractTooltip.select('.date').text('')
    this.contractTooltip.select('.block').text('')

    if (d.state && d.type === 'contract') {
      this.contractTooltip.select('.price').text(() => {
        return `${this.options.valuePrefix || ''} ${d.state.currentAnswer}`
      })
    }
  }

  updateOracleState() {
    const nodeGroup = this.svg
      .select('.network-graph__nodes')
      .selectAll('g.network-graph__node-group')

    nodeGroup
      .select('.network-graph__node__oracle-label--name')
      .transition()
      .duration(600)
      .style('opacity', d => {
        return this.isPendingAnswered(d) ? 1 : 0.6
      })
      .text(d => {
        return d.name
      })

    nodeGroup
      .select('.network-graph__node__oracle-label--price')
      .text(d => {
        return hasPrice(d)
          ? `${this.options.valuePrefix || ''} ${d.state.responseFormatted}`
          : ''
      })
      .transition()
      .duration(600)
      .style('opacity', d => {
        return this.isPendingAnswered(d) ? 1 : 0.6
      })

    nodeGroup
      .select('circle')
      .transition()
      .duration(1000)
      .attr('fill', d => {
        return this.isPendingAnswered(d) ? '#2d2a2b' : '#f1f1f1'
      })
      .attr('stroke', d => {
        return this.isPendingAnswered(d) ? '#2d2a2b' : '#8e8e8e'
      })

    nodeGroup.select('.network-graph__node__contract-label--price').text(d => {
      return d.state && d.state.currentAnswer
        ? `${this.options.valuePrefix || ''} ${d.state.currentAnswer}`
        : ''
    })
  }

  updateLinksState() {
    this.links.attr('class', d => {
      return this.isPendingAnswered(d)
        ? 'network-graph__line--success'
        : 'network-graph__line--wait'
    })
  }

  updateState(state, pendingAnswerId) {
    this.pendingAnswerId = pendingAnswerId

    if (!state || _.isEmpty(state)) return

    const nodesData = this.nodes.data()

    // Put state on nodes object

    nodesData.forEach((d, i) => {
      const oracleData = _.find(state, { sender: d.address })
      d.state = oracleData

      if (d.type === 'contract') {
        d.state = _.find(state, { type: 'contract' })
      }

      return d
    })

    this.updateOracleState()
    this.updateLinksState()
  }

  isPendingAnswered(d) {
    if (!d.state) {
      return false
    }
    if (!this.pendingAnswerId) {
      return false
    }
    return d.state.answerId >= this.pendingAnswerId
  }
}
