import * as d3 from 'd3'

let svg
let path
let tooltip
let x
let y
let line
let overlay
let tooltipPrice
let tooltipAnswerId
// let heighestPricePoint

let margin = { top: 30, right: 30, bottom: 30, left: 30 }
let width = 1200
let height = 300

const bisectDate = d3.bisector(d => {
  return d.answerId
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

  tooltipAnswerId = tooltip
    .append('text')
    .attr('class', 'answer-history-graph--timestamp')
    .attr('x', '10')
    .attr('y', '10')

  // heighestPricePoint = svg
  //   .append('circle')
  //   .attr('class', 'answer-history-graph--top')

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

  x = d3
    .scaleLinear()
    .domain(d3.extent(data, d => d.answerId))
    .range([0, width])

  y = d3
    .scaleLinear()
    .domain(d3.extent(data, d => Number(d.response)))
    .range([height, 0])

  line = d3
    .line()
    .x(d => {
      return x(d.answerId)
    })
    .y(d => {
      return y(Number(d.response))
    })
    .curve(d3.curveMonotoneX)

  path.datum(data).attr('d', line)

  const totalLength = path.node().getTotalLength()

  path
    .attr('stroke-dasharray', totalLength + ' ' + totalLength)
    .attr('stroke-dashoffset', totalLength)
    .transition()
    .duration(2000)
    .attr('stroke-dashoffset', 0)

  // const maxY = d3.max(data, d => {
  //   return Number(d.response)
  // })

  // const sortedByValue = [...data].sort(function(a, b) {
  //   return Number(b.response) - Number(a.response)
  // })[0]

  // heighestPricePoint
  //   .style('fill', 'black')
  //   .attr('r', 4)
  //   .attr(
  //     'transform',
  //     'translate(' +
  //       (x(sortedByValue.answerId) + margin.left) +
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
    const d = x0 - d0.answerId > d1.answerId - x0 ? d1 : d0

    tooltip
      .style('display', 'block')
      .attr(
        'transform',
        'translate(' +
          (x(d.answerId) + margin.left) +
          ',' +
          (y(d.response) + margin.top) +
          ')'
      )

    tooltipAnswerId.text(function() {
      return `ID ${d.answerId}`
    })

    tooltipPrice.text(function() {
      return `$ ${d.responseFormatted}`
    })
  }
}
