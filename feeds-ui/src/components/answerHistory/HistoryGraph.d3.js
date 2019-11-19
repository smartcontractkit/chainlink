import * as d3 from 'd3'
import { formatEthPrice } from 'contracts/utils'
import { ethers } from 'ethers'
import { humanizeUnixTimestamp } from 'utils'

export default class HistoryGraph {
  margin = { top: 30, right: 30, bottom: 30, left: 50 }
  width = 1200
  height = 300
  svg
  path
  tooltip
  x
  y
  line
  overlay
  tooltipPrice
  tooltipTimestamp
  options = {}

  constructor(options) {
    this.options = options
  }

  bisectDate = d3.bisector(d => {
    return d.timestamp
  }).left

  build() {
    this.svg = d3
      .select('.answer-history-graph')
      .append('svg')
      .attr('viewBox', `0 0 ${1300} ${400}`)

    this.path = this.svg
      .append('g')
      .attr(
        'transform',
        'translate(' + this.margin.left + ',' + this.margin.top + ')'
      )
      .append('path')
      .attr('class', 'line')
      .attr('class', 'line')
      .attr('stroke', '#a0a0a0')
      .attr('fill', 'none')

    this.tooltip = this.svg
      .append('g')
      .attr('class', 'tooltip')
      .style('display', 'none')

    this.tooltip
      .append('circle')
      .attr('class', 'y')
      .style('fill', '#2a59da')
      .style('stroke', '#2a59da')
      .attr('r', 4)

    this.tooltipPrice = this.tooltip
      .append('text')
      .attr('class', 'answer-history-graph--price')
      .attr('x', '10')
      .attr('y', '-5')

    this.tooltipTimestamp = this.tooltip
      .append('text')
      .attr('class', 'answer-history-graph--timestamp')
      .attr('x', '10')
      .attr('y', '10')

    this.overlay = this.svg
      .append('rect')
      .attr('width', this.width - this.margin.left)
      .attr('height', this.height)
      .style('fill', 'none')
      .style('pointer-events', 'all')
      .attr(
        'transform',
        'translate(' + this.margin.left + ',' + this.margin.top + ')'
      )
      .on('mouseout', () => {
        this.tooltip.style('display', 'none')
      })
  }

  update(updatedData) {
    if (!updatedData) {
      return
    }
    const data = JSON.parse(JSON.stringify(updatedData))

    this.x = d3
      .scaleLinear()
      .domain(d3.extent(data, d => d.timestamp))
      .range([0, this.width - this.margin.left])

    this.y = d3
      .scaleLinear()
      .domain(d3.extent(data, d => d.response))
      .range([this.height, 0])

    const y_axis = d3
      .axisLeft()
      .scale(this.y)
      .ticks(4)
      .tickFormat(f => formatEthPrice(ethers.utils.bigNumberify(f)))

    this.svg
      .append('g')
      .attr('class', 'y-axis')
      .attr(
        'transform',
        `translate(${this.margin.left - 10}, ${this.margin.top})`
      )
      .call(y_axis)

    const x_axis = d3
      .axisBottom()
      .scale(this.x)
      .ticks(7)
      .tickFormat(f => humanizeUnixTimestamp(f))

    this.svg
      .append('g')
      .attr('class', 'x-axis')
      .attr(
        'transform',
        `translate(${this.margin.left}, ${this.height + this.margin.top + 10})`
      )
      .call(x_axis)

    this.line = d3
      .line()
      .x(d => {
        return this.x(d.timestamp)
      })
      .y(d => {
        return this.y(Number(d.response))
      })
      .curve(d3.curveMonotoneX)

    this.path.datum(data).attr('d', this.line)

    const totalLength = this.path.node().getTotalLength()

    this.path
      .attr('stroke-dasharray', totalLength + ' ' + totalLength)
      .attr('stroke-dashoffset', totalLength)
      .transition()
      .duration(2000)
      .attr('stroke-dashoffset', 0)

    this.overlay.on('mousemove', null)
    this.overlay.on('mousemove', () => mousemove())

    const mousemove = () => {
      const x0 = this.x.invert(d3.mouse(this.overlay.node())[0])
      const i = this.bisectDate(data, x0, 1)
      const d0 = data[i - 1]
      const d1 = data[i]
      if (!d1) {
        return
      }
      const d = x0 - d0.timestamp > d1.timestamp - x0 ? d1 : d0
      this.tooltip
        .style('display', 'block')
        .attr(
          'transform',
          'translate(' +
            (this.x(d.timestamp) + this.margin.left) +
            ',' +
            (this.y(d.response) + this.margin.top) +
            ')'
        )
      this.tooltipTimestamp.text(() => humanizeUnixTimestamp(d.timestamp))
      this.tooltipPrice.text(
        () => `${this.options.valuePrefix} ${d.responseFormatted}`
      )
    }
  }
}
