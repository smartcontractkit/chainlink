# Changelog Chainlink Core
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.8.4] - 2020-05-18

### New features

Fluxmonitor initiators may now optionally include an `absoluteThreshold`
parameter. To trigger a new on-chain report, the absolute difference in the feed
value must change by at least the `absoluteThreshold` value. If it is
unspecified or zero, fluxmonitor behavior is unchanged.

## [0.8.3] - 2020-05-04

### Breaking changes

New cron strings MUST now include time zone. If you want your jobs to run in UTC for example: `CRON_TZ=UTC * * * * *`. Previously, jobs specified without a time zone would run in the server's native time zone, which in most cases is UTC but this was never guaranteed.

### New features

Added Changelog

### Bugfixes

Fix crash in experimental gas updater when run on Kovan network

### Database migrations

There a number of database migrations included in this release as part of our ongoing effort to make the node even more reliable and stable, and build a firm foundation for future development.

## [0.8.2] - 2020-04-20

## [0.8.1] - 2020-04-08

## [0.8.0] - 2020-04-06
