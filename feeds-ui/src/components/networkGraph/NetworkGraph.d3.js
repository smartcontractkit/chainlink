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
      .attr(
        'class',
        d => `network-graph__node-group network-graph__node-group--${d.type}`,
      )
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

    const oracles = this.nodes.filter(d => d.type !== 'contract')

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

    nodesEnter
      .filter(d => d.type === 'contract')
      .append('g')
      .attr('transform', 'translate(-15,-35)')
      .append('path')
      .attr(
        'd',
        'M866.9 169.9L527.1 54.1C523 52.7 517.5 52 512 52s-11 .7-15.1 2.1L157.1 169.9c-8.3 2.8-15.1 12.4-15.1 21.2v482.4c0 8.8 5.7 20.4 12.6 25.9L499.3 968c3.5 2.7 8 4.1 12.6 4.1s9.2-1.4 12.6-4.1l344.7-268.6c6.9-5.4 12.6-17 12.6-25.9V191.1c.2-8.8-6.6-18.3-14.9-21.2zM810 654.3L512 886.5 214 654.3V226.7l298-101.6 298 101.6v427.6zm-405.8-201c-3-4.1-7.8-6.6-13-6.6H336c-6.5 0-10.3 7.4-6.5 12.7l126.4 174a16.1 16.1 0 0 0 26 0l212.6-292.7c3.8-5.3 0-12.7-6.5-12.7h-55.2c-5.1 0-10 2.5-13 6.6L468.9 542.4l-64.7-89.1z',
      )
      .attr('transform', 'scale(0.03)')

    nodesEnter
      .selectAll('.network-graph__node__oracle')
      .on('mouseover', this.setOracleTooltip.bind(this))
      .on('mouseout', this.setOpacity.bind(this.oracleTooltip, 0))

    const label = nodesEnter
      .filter(d => d.type !== 'contract')
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
      .text(d => d.name)

    label
      .append('text')
      .attr('class', 'network-graph__node__oracle-label--price')
      .style('opacity', 0.6)
      .attr('x', '20')
      .attr('y', '10')

    const contract = nodesEnter.filter(d => d.type === 'contract')

    contract
      .on('mouseover', this.setContractTooltip.bind(this))
      .on('mouseout', this.setOpacity.bind(this.contractTooltip, 0))

    const contractLabel = contract
      .append('g')
      .attr('class', 'network-graph__node__contract-label')

    contractLabel
      .append('text')
      .style('opacity', 0)
      .transition()
      .style('opacity', 1)
      .attr('class', 'network-graph__node__contract-label--price')
      .attr('y', '15')
      .attr('text-anchor', 'middle')
      .text('Loading...')

    this.nodes.select('.network-graph__node__contract').attr('class', () => {
      if (this.nodes.data().length > 1) {
        return 'network-graph__node__contract'
      }
      return 'network-graph__node__contract wait'
    })

    this.links = this.links.data(nodes, d => d)

    this.links.exit().remove()

    this.links = this.links
      .enter()
      .append('line')
      .attr('class', 'network-graph__line--wait')
      .attr('x1', this.width / 2)
      .attr('y1', this.height / 2)
      .attr('x2', l => l.x)
      .attr('y2', l => l.y)
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

    this.oracleTooltip.select('.name').text(() => d.name)

    if (d.state) {
      this.oracleTooltip
        .select('.price')
        .text(
          () =>
            `${this.options.valuePrefix || ''} ${d.state.responseFormatted}`,
        )

      this.oracleTooltip
        .select('.date')
        .text(
          () =>
            `Date: ${moment
              .unix(d.state.meta.timestamp)
              .format('DD/MM/YY h:mm:ss A')}`,
        )
      this.oracleTooltip
        .select('.block')
        .text(() => `Block: ${d.state.meta.blockNumber}`)
    }
  }

  setOpacity(opacity) {
    return this.style('opacity', opacity)
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
      this.contractTooltip
        .select('.price')
        .text(
          () => `${this.options.valuePrefix || ''} ${d.state.currentAnswer}`,
        )
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
      .style('opacity', d => (this.isPendingAnswered(d) ? 1 : 0.6))
      .text(d => d.name)

    nodeGroup
      .select('.network-graph__node__oracle-label--price')
      .text(d =>
        hasPrice(d)
          ? `${this.options.valuePrefix || ''} ${d.state.responseFormatted}`
          : '',
      )
      .transition()
      .duration(600)
      .style('opacity', d => (this.isPendingAnswered(d) ? 1 : 0.6))

    nodeGroup
      .select('.network-graph__node__oracle')
      .attr('class', d =>
        this.isPendingAnswered(d)
          ? 'network-graph__node__oracle fulfilled'
          : 'network-graph__node__oracle fetching',
      )

    nodeGroup
      .select('.network-graph__node__contract-label--price')
      .text(d =>
        d.state && d.state.currentAnswer
          ? `${this.options.valuePrefix || ''} ${d.state.currentAnswer}`
          : '',
      )
  }

  updateLinksState() {
    this.links.attr('class', d =>
      this.isPendingAnswered(d)
        ? 'network-graph__line--success'
        : 'network-graph__line--wait',
    )
  }

  updateState(state, pendingAnswerId) {
    this.pendingAnswerId = pendingAnswerId

    if (!state || _.isEmpty(state)) return

    const nodesData = this.nodes.data()

    // Put state on nodes object

    nodesData.forEach(d => {
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
