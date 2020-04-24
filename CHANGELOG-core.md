# Changelog Chainlink Core
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fluxmonitor

- The rules by which a fluxmonitor initiator decides whether to report a result
  onchain are now expressed in a `valueTriggers` entry on the initiator params
  object. For instance,
  
  ```json
  "valueTriggers": {"absoluteThreshold": 0.0059, "relativeThreshold": 0.59}
  ```
  
  means that a new value will be reported onchain if it differs by at least
  0.0059 units AND it's a .59% change relative to the current onchain value.
  Either threshold may be left off the `valueTriggers` object, but at least one
  must be present. New trigger functions should implement the `TriggerFn`
  interface from the `core/store/models/triggerfns` package, and should be
  registered with `triggerfns.RegisterTriggerFunctionFactory`. See
  `core/services/fluxmonitor/trigger_fns.go` for the examples of
  `absoluteThreshold` and `relativeThreshold`.
  
  All trigger functions in a fluxmonitor initiator must trigger, for an onchain
  report to be triggered.

### Added
- Added Changelog

## [0.8.2] - 2020-04-20

## [0.8.1] - 2020-04-08

## [0.8.0] - 2020-04-06
