# Neo-Async

<img src="https://raw.githubusercontent.com/wiki/suguru03/neo-async/images/neo_async_v2.png" width="230px" />

[![npm](https://img.shields.io/npm/v/neo-async.svg)](https://www.npmjs.com/package/neo-async)
[![Travis](https://img.shields.io/travis/suguru03/neo-async.svg)](https://travis-ci.org/suguru03/neo-async)
[![Codecov](https://img.shields.io/codecov/c/github/suguru03/neo-async.svg)](https://codecov.io/github/suguru03/neo-async?branch=master)
[![Dependency Status](https://gemnasium.com/suguru03/neo-async.svg)](https://gemnasium.com/suguru03/neo-async)
[![npm](https://img.shields.io/npm/dm/neo-async.svg)](https://www.npmjs.com/package/neo-async)

Neo-Async is thought to be used as a drop-in replacement for [Async](https://github.com/caolan/async), it almost fully covers its functionality and runs [faster](#benchmark).

Benchmark is [here](#benchmark)!

Bluebird's benchmark is [here](https://github.com/suguru03/bluebird/tree/aigle/benchmark)!

## Code Coverage
![coverage](https://raw.githubusercontent.com/wiki/suguru03/neo-async/images/coverage.png)

## Installation

### In a browser
```html
<script src="async.min.js"></script>
```

### In an AMD loader
```js
require(['async'], function(async) {});
```

### Promise and async/await

I recommend to use [`Aigle`](https://github.com/suguru03/aigle).

It is optimized for Promise handling and has almost the same functionality as `neo-async`.

### Node.js

#### standard

```bash
$ npm install neo-async
```
```js
var async = require('neo-async');
```

#### replacement
```bash
$ npm install neo-async
$ ln -s ./node_modules/neo-async ./node_modules/async
```
```js
var async = require('async');
```

### Bower

```bash
bower install neo-async
```

## Feature

[JSDoc](http://suguru03.github.io/neo-async/doc/async.html)

\* not in Async

### Collections

- [`each`](http://suguru03.github.io/neo-async/doc/async.each.html)
- [`eachSeries`](http://suguru03.github.io/neo-async/doc/async.eachSeries.html)
- [`eachLimit`](http://suguru03.github.io/neo-async/doc/async.eachLimit.html)
- [`forEach`](http://suguru03.github.io/neo-async/doc/async.each.html) -> [`each`](http://suguru03.github.io/neo-async/doc/async.each.html)
- [`forEachSeries`](http://suguru03.github.io/neo-async/doc/async.eachSeries.html) -> [`eachSeries`](http://suguru03.github.io/neo-async/doc/async.eachSeries.html)
- [`forEachLimit`](http://suguru03.github.io/neo-async/doc/async.eachLimit.html) -> [`eachLimit`](http://suguru03.github.io/neo-async/doc/async.eachLimit.html)
- [`eachOf`](http://suguru03.github.io/neo-async/doc/async.each.html) -> [`each`](http://suguru03.github.io/neo-async/doc/async.each.html)
- [`eachOfSeries`](http://suguru03.github.io/neo-async/doc/async.eachSeries.html) -> [`eachSeries`](http://suguru03.github.io/neo-async/doc/async.eachSeries.html)
- [`eachOfLimit`](http://suguru03.github.io/neo-async/doc/async.eachLimit.html) -> [`eachLimit`](http://suguru03.github.io/neo-async/doc/async.eachLimit.html)
- [`forEachOf`](http://suguru03.github.io/neo-async/doc/async.each.html) -> [`each`](http://suguru03.github.io/neo-async/doc/async.each.html)
- [`forEachOfSeries`](http://suguru03.github.io/neo-async/doc/async.eachSeries.html) -> [`eachSeries`](http://suguru03.github.io/neo-async/doc/async.eachSeries.html)
- [`eachOfLimit`](http://suguru03.github.io/neo-async/doc/async.eachLimit.html) -> [`forEachLimit`](http://suguru03.github.io/neo-async/doc/async.eachLimit.html)
- [`map`](http://suguru03.github.io/neo-async/doc/async.map.html)
- [`mapSeries`](http://suguru03.github.io/neo-async/doc/async.mapSeries.html)
- [`mapLimit`](http://suguru03.github.io/neo-async/doc/async.mapLimit.html)
- [`mapValues`](http://suguru03.github.io/neo-async/doc/async.mapValues.html)
- [`mapValuesSeries`](http://suguru03.github.io/neo-async/doc/async.mapValuesSeries.html)
- [`mapValuesLimit`](http://suguru03.github.io/neo-async/doc/async.mapValuesLimit.html)
- [`filter`](http://suguru03.github.io/neo-async/doc/async.filter.html)
- [`filterSeries`](http://suguru03.github.io/neo-async/doc/async.filterSeries.html)
- [`filterLimit`](http://suguru03.github.io/neo-async/doc/async.filterLimit.html)
- [`select`](http://suguru03.github.io/neo-async/doc/async.filter.html) -> [`filter`](http://suguru03.github.io/neo-async/doc/async.filter.html)
- [`selectSeries`](http://suguru03.github.io/neo-async/doc/async.filterSeries.html) -> [`filterSeries`](http://suguru03.github.io/neo-async/doc/async.filterSeries.html)
- [`selectLimit`](http://suguru03.github.io/neo-async/doc/async.filterLimit.html) -> [`filterLimit`](http://suguru03.github.io/neo-async/doc/async.filterLimit.html)
- [`reject`](http://suguru03.github.io/neo-async/doc/async.reject.html)
- [`rejectSeries`](http://suguru03.github.io/neo-async/doc/async.rejectSeries.html)
- [`rejectLimit`](http://suguru03.github.io/neo-async/doc/async.rejectLimit.html)
- [`detect`](http://suguru03.github.io/neo-async/doc/async.detect.html)
- [`detectSeries`](http://suguru03.github.io/neo-async/doc/async.detectSeries.html)
- [`detectLimit`](http://suguru03.github.io/neo-async/doc/async.detectLimit.html)
- [`find`](http://suguru03.github.io/neo-async/doc/async.detect.html) -> [`detect`](http://suguru03.github.io/neo-async/doc/async.detect.html)
- [`findSeries`](http://suguru03.github.io/neo-async/doc/async.detectSeries.html) -> [`detectSeries`](http://suguru03.github.io/neo-async/doc/async.detectSeries.html)
- [`findLimit`](http://suguru03.github.io/neo-async/doc/async.detectLimit.html) -> [`detectLimit`](http://suguru03.github.io/neo-async/doc/async.detectLimit.html)
- [`pick`](http://suguru03.github.io/neo-async/doc/async.pick.html) *
- [`pickSeries`](http://suguru03.github.io/neo-async/doc/async.pickSeries.html) *
- [`pickLimit`](http://suguru03.github.io/neo-async/doc/async.pickLimit.html) *
- [`omit`](http://suguru03.github.io/neo-async/doc/async.omit.html) *
- [`omitSeries`](http://suguru03.github.io/neo-async/doc/async.omitSeries.html) *
- [`omitLimit`](http://suguru03.github.io/neo-async/doc/async.omitLimit.html) *
- [`reduce`](http://suguru03.github.io/neo-async/doc/async.reduce.html)
- [`inject`](http://suguru03.github.io/neo-async/doc/async.reduce.html) -> [`reduce`](http://suguru03.github.io/neo-async/doc/async.reduce.html)
- [`foldl`](http://suguru03.github.io/neo-async/doc/async.reduce.html) -> [`reduce`](http://suguru03.github.io/neo-async/doc/async.reduce.html)
- [`reduceRight`](http://suguru03.github.io/neo-async/doc/async.reduceRight.html)
- [`foldr`](http://suguru03.github.io/neo-async/doc/async.reduceRight.html) -> [`reduceRight`](http://suguru03.github.io/neo-async/doc/async.reduceRight.html)
- [`transform`](http://suguru03.github.io/neo-async/doc/async.transform.html)
- [`transformSeries`](http://suguru03.github.io/neo-async/doc/async.transformSeries.html) *
- [`transformLimit`](http://suguru03.github.io/neo-async/doc/async.transformLimit.html) *
- [`sortBy`](http://suguru03.github.io/neo-async/doc/async.sortBy.html)
- [`sortBySeries`](http://suguru03.github.io/neo-async/doc/async.sortBySeries.html) *
- [`sortByLimit`](http://suguru03.github.io/neo-async/doc/async.sortByLimit.html) *
- [`some`](http://suguru03.github.io/neo-async/doc/async.some.html)
- [`someSeries`](http://suguru03.github.io/neo-async/doc/async.someSeries.html)
- [`someLimit`](http://suguru03.github.io/neo-async/doc/async.someLimit.html)
- [`any`](http://suguru03.github.io/neo-async/doc/async.some.html) -> [`some`](http://suguru03.github.io/neo-async/doc/async.some.html)
- [`anySeries`](http://suguru03.github.io/neo-async/doc/async.someSeries.html) -> [`someSeries`](http://suguru03.github.io/neo-async/doc/async.someSeries.html)
- [`anyLimit`](http://suguru03.github.io/neo-async/doc/async.someLimit.html) -> [`someLimit`](http://suguru03.github.io/neo-async/doc/async.someLimit.html)
- [`every`](http://suguru03.github.io/neo-async/doc/async.every.html)
- [`everySeries`](http://suguru03.github.io/neo-async/doc/async.everySeries.html)
- [`everyLimit`](http://suguru03.github.io/neo-async/doc/async.everyLimit.html)
- [`all`](http://suguru03.github.io/neo-async/doc/async.every.html) -> [`every`](http://suguru03.github.io/neo-async/doc/async.every.html)
- [`allSeries`](http://suguru03.github.io/neo-async/doc/async.everySeries.html) -> [`every`](http://suguru03.github.io/neo-async/doc/async.everySeries.html)
- [`allLimit`](http://suguru03.github.io/neo-async/doc/async.everyLimit.html) -> [`every`](http://suguru03.github.io/neo-async/doc/async.everyLimit.html)
- [`concat`](http://suguru03.github.io/neo-async/doc/async.concat.html)
- [`concatSeries`](http://suguru03.github.io/neo-async/doc/async.concatSeries.html)
- [`concatLimit`](http://suguru03.github.io/neo-async/doc/async.concatLimit.html) *

### Control Flow

- [`parallel`](http://suguru03.github.io/neo-async/doc/async.parallel.html)
- [`series`](http://suguru03.github.io/neo-async/doc/async.series.html)
- [`parallelLimit`](http://suguru03.github.io/neo-async/doc/async.series.html)
- [`tryEach`](http://suguru03.github.io/neo-async/doc/async.tryEach.html)
- [`waterfall`](http://suguru03.github.io/neo-async/doc/async.waterfall.html)
- [`angelFall`](http://suguru03.github.io/neo-async/doc/async.angelFall.html) *
- [`angelfall`](http://suguru03.github.io/neo-async/doc/async.angelFall.html) -> [`angelFall`](http://suguru03.github.io/neo-async/doc/async.angelFall.html) *
- [`whilst`](#whilst)
- [`doWhilst`](#doWhilst)
- [`until`](#until)
- [`doUntil`](#doUntil)
- [`during`](#during)
- [`doDuring`](#doDuring)
- [`forever`](#forever)
- [`compose`](#compose)
- [`seq`](#seq)
- [`applyEach`](#applyEach)
- [`applyEachSeries`](#applyEachSeries)
- [`queue`](#queue)
- [`priorityQueue`](#priorityQueue)
- [`cargo`](#cargo)
- [`auto`](#auto)
- [`autoInject`](#autoInject)
- [`retry`](#retry)
- [`retryable`](#retryable)
- [`iterator`](#iterator)
- [`times`](http://suguru03.github.io/neo-async/doc/async.times.html)
- [`timesSeries`](http://suguru03.github.io/neo-async/doc/async.timesSeries.html)
- [`timesLimit`](http://suguru03.github.io/neo-async/doc/async.timesLimit.html)
- [`race`](#race)

### Utils
- [`apply`](#apply)
- [`setImmediate`](#setImmediate)
- [`nextTick`](#nextTick)
- [`memoize`](#memoize)
- [`unmemoize`](#unmemoize)
- [`ensureAsync`](#ensureAsync)
- [`constant`](#constant)
- [`asyncify`](#asyncify)
- [`wrapSync`](#asyncify) -> [`asyncify`](#asyncify)
- [`log`](#log)
- [`dir`](#dir)
- [`timeout`](http://suguru03.github.io/neo-async/doc/async.timeout.html)
- [`reflect`](#reflect)
- [`reflectAll`](#reflectAll)
- [`createLogger`](#createLogger)

## Mode
- [`safe`](#safe) *
- [`fast`](#fast) *

## Benchmark

[Benchmark: Async vs Neo-Async](http://suguru03.hatenablog.com/entry/2016/06/10/135559)

### How to check

```bash
$ git clone git@github.com:suguru03/async-benchmark.git
$ cd async-benchmark
$ npm install
$ node . // It might take more than one hour...
```

### Environment

* Ubuntu v12.04
* Node.js v6.2.1
* async v2.0.0-rc.6
* neo-async v2.0.0-rc.1
* benchmark v2.1.0
* func-comparator v0.7.1

### Result

Neo-Async is 1.27 ~ 10.7 times faster than Async.

The value is the ratio (Neo-Async/Async) of the average speed.

#### Collections
|function|benchmark|func-comparator|
|---|--:|--:|
|each|3.71|2.54|
|eachSeries|2.14|1.90|
|eachLimit|2.14|1.88|
|eachOf|3.30|2.50|
|eachOfSeries|1.97|1.83|
|eachOfLimit|2.02|1.80|
|map|4.20|4.11|
|mapSeries|2.40|3.65|
|mapLimit|2.64|2.66|
|mapValues|5.71|5.32|
|mapValuesSeries|3.82|3.23|
|mapValuesLimit|3.10|2.38|
|filter|8.11|8.76|
|filterSeries|5.79|4.86|
|filterLimit|4.00|3.32|
|reject|9.47|9.52|
|rejectSeries|7.39|4.64|
|rejectLimit|4.54|3.49|
|detect|6.67|6.37|
|detectSeries|3.54|3.73|
|detectLimit|2.38|2.62|
|reduce|4.13|3.23|
|reduceRight|4.23|3.24|
|transform|5.30|5.17|
|sortBy|2.24|2.37|
|some|6.39|6.10|
|someSeries|5.37|4.66|
|someLimit|3.39|2.84|
|every|6.85|6.27|
|everySeries|4.53|3.90|
|everyLimit|3.36|2.75|
|concat|9.18|9.35|
|concatSeries|7.49|6.09|

#### Control Flow
|funciton|benchmark|func-comparator|
|---|--:|--:|
|parallel|7.54|5.45|
|series|3.29|2.41|
|waterfall|5.12|4.27|
|whilst|1.96|1.95|
|doWhilst|2.07|1.96|
|until|2.10|1.99|
|doUntil|1.98|2.04|
|during|10.7|7.09|
|doDuring|5.98|6.03|
|queue|1.83|1.75|
|priorityQueue|1.79|1.75|
|times|3.84|3.65|
|race|1.45|1.27|
|auto|3.23|3.50|
|retry|9.43|6.78|

