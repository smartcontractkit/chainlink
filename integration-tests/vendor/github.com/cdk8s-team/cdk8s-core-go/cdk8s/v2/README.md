# cdk8s

### Cloud Development Kit for Kubernetes

[![build](https://github.com/cdk8s-team/cdk8s-core/workflows/release/badge.svg)](https://github.com/cdk8s-team/cdk8s-core/actions/workflows/release.yml)
[![npm version](https://badge.fury.io/js/cdk8s.svg)](https://badge.fury.io/js/cdk8s)
[![PyPI version](https://badge.fury.io/py/cdk8s.svg)](https://badge.fury.io/py/cdk8s)
[![Maven Central](https://maven-badges.herokuapp.com/maven-central/org.cdk8s/cdk8s/badge.svg)](https://maven-badges.herokuapp.com/maven-central/org.cdk8s/cdk8s)

**cdk8s** is a software development framework for defining Kubernetes
applications using rich object-oriented APIs. It allows developers to leverage
the full power of software in order to define abstract components called
"constructs" which compose Kubernetes resources or other constructs into
higher-level abstractions.

> **Note:** This repository is the "core library" of cdk8s, with logic for synthesizing Kubernetes manifests using the [constructs framework](https://github.com/aws/constructs). It is published to NPM as [`cdk8s`](https://www.npmjs.com/package/cdk8s) and should not be confused with the cdk8s command-line tool [`cdk8s-cli`](https://www.npmjs.com/package/cdk8s-cli). For more general information about cdk8s, please see [cdk8s.io](https://cdk8s.io), or visit the umbrella repository located at [cdk8s-team/cdk8s](https://github.com/cdk8s-team/cdk8s).

## Documentation

See [cdk8s.io](https://cdk8s.io).

## License

This project is distributed under the [Apache License, Version 2.0](./LICENSE).

This module is part of the [cdk8s project](https://github.com/cdk8s-team/cdk8s).
