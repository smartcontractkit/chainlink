import * as d3 from 'd3'
import _ from 'lodash'
import moment from 'moment'

const theme = {
  brand: '#2a5ada',
  backgroud: '#f1f1f1',
  strokeActive: '#a0a0a0',
  nodeActive: '#2d2a2b'
}

let svg
let path
let tooltip
let x
let y
let line
let overlay
let tooltipPrice
let tooltipTimestamp
let topValue

let margin = { top: 30, right: 30, bottom: 30, left: 30 }
let width = 1200
let height = 300

const bisectDate = d3.bisector(d => {
  return d.timestamp
}).left

export function createChart() {
  svg = d3
    .select('.answer-history-graph')
    .append('svg')
    .attr('viewBox', `0 0 ${1300} ${400}`)

  path = svg
    .append('g')
    .attr('transform', 'translate(' + margin.left + ',' + margin.top + ')')
    .append('path')
    .attr('class', 'line')
    .attr('class', 'line')
    .attr('stroke', '#a0a0a0')
    .attr('fill', 'none')

  tooltip = svg
    .append('g')
    .attr('class', 'tooltip')
    .style('display', 'none')

  tooltip
    .append('circle')
    .attr('class', 'y')
    .style('fill', '#2a59da')
    .style('stroke', '#2a59da')
    .attr('r', 4)

  tooltipPrice = tooltip
    .append('text')
    .attr('class', 'answer-history-graph--price')
    .attr('x', '10')
    .attr('y', '-5')

  tooltipTimestamp = tooltip
    .append('text')
    .attr('class', 'answer-history-graph--timestamp')
    .attr('x', '10')
    .attr('y', '10')

  topValue = svg.append('circle').attr('class', 'answer-history-graph--top')

  overlay = svg
    .append('rect')
    .attr('width', width)
    .attr('height', height)
    .style('fill', 'none')
    .style('pointer-events', 'all')
    .attr('transform', 'translate(' + margin.left + ',' + margin.top + ')')
    .on('mouseout', function() {
      tooltip.style('display', 'none')
    })
}

export function update(updatedData) {
  if (!updatedData) {
    return
  }
  const data = JSON.parse(JSON.stringify(updatedData))

  line = d3
    .line()
    .x(function(d) {
      return x(d.timestamp)
    })
    .y(function(d) {
      return y(d.response)
    })
    .curve(d3.curveMonotoneX)

  x = d3
    .scaleLinear()
    .domain(d3.extent(data, d => d.timestamp))
    .range([0, width])

  y = d3
    .scaleLinear()
    .domain(d3.extent(data, d => Number(d.response)))
    .range([height, 0])

  path.datum(data).attr('d', line)

  // const maxY = d3.max(data, d => {
  //   return Number(d.response)
  // })

  // const sortedByValue = [...data].sort(function(a, b) {
  //   return Number(b.response) - Number(a.response)
  // })[0]

  // topValue
  //   .style('fill', '#2a59da')
  //   .style('stroke', '#2a59da')
  //   .attr('r', 4)
  //   .attr(
  //     'transform',
  //     'translate(' +
  //       (x(sortedByValue.timestamp) + margin.left) +
  //       ',' +
  //       (y(maxY) + margin.top) +
  //       ')'
  //   )

  overlay.on('mousemove', null)
  overlay.on('mousemove', mousemove)

  function mousemove() {
    const x0 = x.invert(d3.mouse(this)[0])
    const i = bisectDate(data, x0, 1)
    const d0 = data[i - 1]
    const d1 = data[i]
    if (!d1) {
      return
    }
    const d = x0 - d0.timestamp > d1.timestamp - x0 ? d1 : d0

    tooltip
      .style('display', 'block')
      .attr(
        'transform',
        'translate(' +
          (x(d.timestamp) + margin.left) +
          ',' +
          (y(d.response) + margin.top) +
          ')'
      )

    tooltipTimestamp.text(function() {
      return moment.unix(d.timestamp).format('hh:mm:ss A')
    })

    tooltipPrice.text(function() {
      return `$ ${d.responseFormatted}`
    })
  }
}
