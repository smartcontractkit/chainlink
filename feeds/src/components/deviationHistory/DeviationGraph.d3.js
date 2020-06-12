import * as d3 from 'd3'
import { humanizeUnixTimestamp } from 'utils'

export default class DeviationGraph {
  margin = { top: 30, right: 30, bottom: 30, left: 50 }
  width = 1300
  height = 250
  svg
  path
  tooltip
  info
  x
  y
  line
  overlay
  tooltipPrice
  tooltipTimestamp
  tooltipPercentage
  config = {}

  constructor(config) {
    this.config = config
  }

  bisectDate = d3.bisector(d => d.timestamp).left

  build() {
    this.svg = d3
      .select('.deviation-history-graph')
      .append('svg')
      .attr(
        'viewBox',
        `0 0 ${this.width} ${this.height +
          this.margin.top +
          this.margin.bottom}`,
      )

    this.clip = this.svg
      .append('defs')
      .append('svg:clipPath')
      .attr('id', 'clip')
      .append('svg:rect')
      .attr('width', this.width - this.margin.left)
      .attr('height', this.height)
      .attr('x', 0)
      .attr('y', 0)

    this.path = this.svg
      .append('g')
      .attr(
        'transform',
        'translate(' + this.margin.left + ',' + this.margin.top + ')',
      )
      .append('path')
      .attr('class', 'line')
      .attr('stroke', '#a0a0a0')
      .attr('fill', 'none')
      .attr('clip-path', 'url(#clip)')

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

    this.info = this.svg
      .append('g')
      .attr('class', 'info')
      .attr(
        'transform',
        'translate(' + this.width + ',' + this.margin.top + ')',
      )

    this.tooltipPercentage = this.info
      .append('text')
      .attr('class', 'deviation-history-graph--percentage')
      .attr('x', '0')
      .attr('y', '0')

    this.tooltipPrice = this.info
      .append('text')
      .attr('class', 'deviation-history-graph--price')
      .attr('x', '0')
      .attr('y', '15')

    this.tooltipTimestamp = this.info
      .append('text')
      .attr('class', 'deviation-history-graph--timestamp')
      .attr('x', '0')
      .attr('y', '30')

    this.brushX = d3.brushX()

    this.brush = this.svg
      .append('g')
      .attr('class', 'brush')
      .attr(
        'transform',
        'translate(' + this.margin.left + ',' + this.margin.top + ')',
      )
      .on('mouseout', () => {
        this.tooltip.style('display', 'none')
        this.info.style('display', 'none')
      })

    this.zoomOutBtn = this.svg
      .append('g')
      .attr(
        'transform',
        'translate(' + (this.width - 70) + ',' + this.height + ')',
      )
      .style('opacity', 0)
      .style('cursor', 'pointer')

    this.zoomOutBtn
      .append('rect')
      .attr('class', 'y')
      .style('fill', '#375bd2')
      .style('stroke', '#375bd2')
      .attr('width', '70')
      .attr('height', '25')

    this.zoomOutBtn
      .append('text')
      .attr('fill', '#fff')
      .text('Zoom out')
      .attr('x', '10')
      .attr('y', '16')
      .style('font-size', 12)
  }

  updateBrushed() {
    if (!this.x) {
      return
    }

    const extent = d3.event.selection
    if (extent) {
      this.x.domain([this.x.invert(extent[0]), this.x.invert(extent[1])])
      this.brush.call(this.brushX.move, null)

      this.zoomOutBtn
        .transition()
        .duration(300)
        .style('opacity', 1)
    }

    this.xAxis
      .transition()
      .duration(300)
      .call(
        d3
          .axisBottom()
          .scale(this.x)
          .ticks(7)
          .tickFormat(f => humanizeUnixTimestamp(f)),
      )

    const line = d3
      .line()
      .x(d => this.x(d.timestamp))
      .y(d => this.y(Number(d.deviation)))

    this.path
      .transition()
      .duration(300)
      .attr('d', line)
  }

  zoomOut(data) {
    this.x.domain(d3.extent(data, d => d.timestamp))

    const xAxis = d3
      .axisBottom()
      .scale(this.x)
      .tickFormat(f => humanizeUnixTimestamp(f))

    this.xAxis
      .transition()
      .duration(300)
      .call(xAxis)

    const line = d3
      .line()
      .x(d => this.x(d.timestamp))
      .y(d => this.y(Number(d.deviation)))

    this.path
      .transition()
      .duration(300)
      .attr('d', line)

    this.zoomOutBtn
      .transition()
      .duration(300)
      .style('opacity', 0)
  }

  addDeviation(data) {
    const reducedData = data.reduce((memo, current, index) => {
      const slice = memo.slice(index - 1, index)
      const average = d3.mean(slice, d => d.answer) || current.answer

      const deviation =
        (100 * (current.answer - average)) / ((current.answer + average) / 2)

      return [
        ...memo,
        ...[
          {
            ...current,
            ...{ deviation: Math.abs(Number(deviation.toFixed(3))) },
          },
        ],
      ]
    }, [])
    return reducedData
  }

  update(updatedData) {
    if (!updatedData) {
      return
    }

    const parsedData = JSON.parse(JSON.stringify(updatedData))
    const data = this.addDeviation(parsedData)

    this.x = d3
      .scaleLinear()
      .domain(d3.extent(data, d => d.timestamp))
      .range([0, this.width - this.margin.left])

    const ymin = d3.min(data, d => d.deviation)
    const ymax = d3.max(data, d => d.deviation)

    this.y = d3
      .scaleLinear()
      .domain([Math.floor(ymin), Math.ceil(ymax)])
      .range([this.height, 0])

    const yAxis = d3
      .axisLeft()
      .scale(this.y)
      .ticks(3)
      .tickFormat(f => `${f}%`)

    this.svg
      .append('g')
      .attr('class', 'y-axis')
      .attr(
        'transform',
        `translate(${this.margin.left - 10}, ${this.margin.top})`,
      )
      .call(yAxis)

    const xAxis = d3
      .axisBottom()
      .scale(this.x)
      .ticks(7)
      .tickFormat(f => humanizeUnixTimestamp(f))

    this.xAxis = this.svg
      .append('g')
      .attr('class', 'x-axis')
      .attr(
        'transform',
        `translate(${this.margin.left}, ${this.height + this.margin.top + 10})`,
      )
      .call(xAxis)

    this.line = d3
      .line()
      .x(d => this.x(d.timestamp))
      .y(d => this.y(Number(d.deviation)))

    this.path
      .attr('stroke-width', 1)
      .attr('stroke', '#a0a0a0')
      .attr('fill', 'none')
      .datum(data)
      .attr('d', this.line)

    this.brush.on('mousemove', () => this.mousemove(data))

    this.svg.on('dblclick', () => this.zoomOut(data))
    this.zoomOutBtn.on('click', () => this.zoomOut(data))

    this.brush.call(
      this.brushX
        .extent([
          [0, 0],
          [this.width - this.margin.left, this.height],
        ])
        .on('end', () => this.updateBrushed()),
    )
  }

  mousemove(data) {
    const x0 = this.x.invert(d3.mouse(this.brush.node())[0])
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
          (this.y(d.deviation) + this.margin.top) +
          ')',
      )

    this.info.style('display', 'block')
    this.tooltipTimestamp.text(() =>
      humanizeUnixTimestamp(d.timestamp, 'LL LTS'),
    )
    this.tooltipPrice.text(
      () => `${this.config.valuePrefix} ${d.answerFormatted}`,
    )
    this.tooltipPercentage.text(() => `${d.deviation.toFixed(2)}%`)
  }
}
