import * as d3 from 'd3'
import _ from 'lodash'
import moment from 'moment'

const width = 1200
const height = 600
const theme = {
  brand: '#2a5ada',
  backgroud: '#f1f1f1',
  strokeActive: '#a0a0a0',
  nodeActive: '#2d2a2b'
}

let svg
let simulation
let node
let link
let tooltip
let dragging

function appendSvg() {
  d3.select('.network-graph')
    .append('svg')
    .attr('viewBox', `0 0 ${width} ${height}`)

  svg = d3.select('.network-graph').select('svg')

  svg.append('g').attr('class', 'links')
  svg.append('g').attr('class', 'nodes')

  node = svg.select('.nodes').selectAll('g')
  link = svg.select('.links').selectAll('line')

  tooltip = d3
    .select('.network-graph')
    .select('.network-graph__tooltip')
    .style('opacity', 0)
}

function createSimulation() {
  simulation = d3
    .forceSimulation()
    .force(
      'link',
      d3.forceLink().id(d => {
        return d.id
      })
    )
    .force('charge', d3.forceManyBody().strength(-1000))
    .force(
      'collision',
      d3.forceCollide().radius(d => {
        return d.radius
      })
    )
    .force('center', d3.forceCenter(width / 2, height / 2))
    .on('tick', ticked)
}

function ticked() {
  for (let i = 0; i < 5; i++) {
    simulation.tick()
  }

  svg
    .select('.links')
    .selectAll('line')
    .attr('x1', function(d) {
      return d.source.x
    })
    .attr('y1', function(d) {
      return d.source.y
    })
    .attr('x2', function(d) {
      return d.target.x
    })
    .attr('y2', function(d) {
      return d.target.y
    })

  svg
    .select('.nodes')
    .selectAll('.network-graph__node-group')
    .attr('transform', function(d) {
      return 'translate(' + d.x + ',' + d.y + ')'
    })
}

var createNodes = function(nodes, numNodes, radius) {
  let width = radius * 2 + 50,
    height = radius * 2 + 50,
    angle,
    x,
    y,
    i
  for (i = 0; i < numNodes; i++) {
    angle = (i / (numNodes / 2)) * Math.PI
    x = radius * Math.cos(angle) + width / 2
    y = radius * Math.sin(angle) + height / 2
    nodes[i].x = x
    nodes[i].y = y
  }
  return nodes
}

function update(data, linksData) {
  linksData = JSON.parse(JSON.stringify(linksData))

  createNodes(data, data.length, 100)

  node = node.data(data, d => {
    return d
  })

  node.exit().remove()

  node = node
    .enter()
    .append('g')
    .attr('class', 'network-graph__node-group')
    .on('update-state', updateStateEvent)
    .merge(node)
    .call(
      d3
        .drag()
        .on('start', dragstarted)
        .on('drag', dragged)
        .on('end', dragended)
    )

  node
    .append('circle')
    .attr('class', d => d.type)
    .attr('r', d => {
      return d.type === 'contract' ? 30 : 10
    })
    .on('mouseover', d => {
      if (dragging) {
        return tooltip.style('opacity', 0)
      }

      tooltip
        .style('opacity', 1)
        .style('left', d.x + 15 + 'px')
        .style('top', d.y + 'px')

      tooltip.select('.network-graph__tooltip--price').text('')
      tooltip.select('.network-graph__tooltip--date').text('')
      tooltip.select('.network-graph__tooltip--block').text('')

      tooltip.select('.network-graph__tooltip--type').text(() => {
        return d.type
      })

      tooltip.select('.network-graph__tooltip--name').text(() => {
        return d.name
      })

      if (d.state && d.type === 'oracle') {
        tooltip.select('.network-graph__tooltip--price').text(() => {
          return `$ ${d.state.responseFormatted}`
        })
        tooltip.select('.network-graph__tooltip--date').text(() => {
          return `Date: ${moment
            .unix(d.state.meta.timestamp)
            .format('DD/MM/YY h:mm:ss A')}`
        })
        tooltip.select('.network-graph__tooltip--block').text(() => {
          return `Block: ${d.state.meta.blockNumber}`
        })
      }

      if (d.state && d.type === 'contract') {
        tooltip.select('.network-graph__tooltip--price').text(() => {
          return `$ ${d.state.currentAnswer}`
        })
      }
    })
    .on('mouseout', () => {
      tooltip.style('opacity', 0)
    })

  node
    .select('.contract')
    .attr('stroke', theme.strokeActive)
    .attr('fill', theme.backgroud)
    .attr('r', 50)

  const label = node
    .filter(d => {
      return d.type === 'contract' ? false : true
    })
    .append('g')
    .attr('class', 'network-graph__oracle-label')

  label
    .append('text')
    .attr('class', 'network-graph__oracle-name')
    .style('opacity', 0)
    .attr('x', '20')
    .attr('y', '-5')

  label
    .append('text')
    .attr('class', 'network-graph__oracle-price')
    .style('opacity', 0)
    .attr('x', '20')
    .attr('y', '10')

  // contract label

  const contractLabel = node
    .filter(d => {
      return d.type === 'contract'
    })
    .append('g')
    .attr('class', 'network-contract-label')

  contractLabel
    .append('text')
    .attr('class', 'network-graph__contract-price')
    .attr('y', '5')
    .attr('text-anchor', 'middle')

  // links

  link = link.data(linksData, function(d) {
    return d
  })

  link.exit().remove()

  link = link
    .enter()
    .append('line')
    .merge(link)

  simulation.nodes(data)
  simulation
    .force('link')
    .distance(240)
    .strength(1)
    .links(linksData)

  simulation.alpha(1).restart()
}

function dragstarted(e) {
  if (!d3.event.active) simulation.alphaTarget(0.3).restart()
  dragging = true
  tooltip.style('opacity', 0)

  d3.event.subject.fx = d3.event.subject.x
  d3.event.subject.fy = d3.event.subject.y
}
function dragged() {
  d3.event.subject.fx = d3.event.x
  d3.event.subject.fy = d3.event.y
}
function dragended() {
  dragging = false

  if (!d3.event.active) simulation.alphaTarget(0)
  d3.event.subject.fx = null
  d3.event.subject.fy = null
}

export function createChart() {
  appendSvg()
  createSimulation()
}

export function updateData(nodesData, linksData) {
  if (!nodesData.length || !linksData.length) {
    return
  }
  update(nodesData, linksData)
}

function hasPrice(d) {
  return d.state && d.state.responseFormatted
}

function updateStateEvent(d, i) {}

function updateOracles() {
  const nodeGroup = d3.selectAll('.network-graph__node-group')

  nodeGroup
    .selectAll('.network-graph__oracle-label')
    .transition()
    .style('fill', d => {
      return hasPrice(d) ? '#505050' : '#8e8e8e'
    })

  nodeGroup
    .selectAll('.network-graph__oracle-name')
    .transition()
    .duration(1000)
    .style('opacity', 1)
    .text(d => {
      return d.name
    })

  nodeGroup
    .selectAll('.network-graph__oracle-price')
    .transition()
    .duration(1000)
    .style('opacity', d => {
      return hasPrice(d) ? 1 : 0
    })
    .text(function(d) {
      return hasPrice(d) ? `$ ${d.state.responseFormatted}` : ''
    })
    .transition()
    .duration(1000)
    .style('opacity', d => {
      return hasPrice(d) ? 1 : 0
    })

  node
    .selectAll('circle.oracle')
    .transition()
    .duration(1000)
    .attr('fill', d => {
      return hasPrice(d) ? theme.nodeActive : theme.backgroud
    })
    .attr('stroke', d => {
      return hasPrice(d) ? theme.nodeActive : '#8e8e8e'
    })

  nodeGroup.selectAll('.network-graph__contract-price').text(d => {
    return d.state && d.state.currentAnswer ? `$ ${d.state.currentAnswer}` : ''
  })
}

function updateLinks() {
  link.attr('class', d => {
    return d.target.state && d.target.state.responseFormatted
      ? 'network-graph__line-success'
      : 'network-graph__line-wait'
  })
}

export function updateState(data) {
  if (!data) return

  const nodeData = node.data()

  nodeData.forEach((d, i) => {
    const oracleData = _.find(data, { sender: d.address })
    d.state = oracleData
    if (d.type === 'contract') {
      d.state = _.find(data, { type: 'contract' })
    }
    return d
  })

  updateOracles()
  updateLinks()
}

export const updateForce = () => {
  simulation.force('link').distance(100)
  simulation.alpha(1).restart()
}
