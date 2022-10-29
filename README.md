# Semver

[![License](https://img.shields.io/github/license/FollowTheProcess/semver)](https://github.com/FollowTheProcess/semver)
[![Go Report Card](https://goreportcard.com/badge/github.com/FollowTheProcess/semver)](https://goreportcard.com/report/github.com/FollowTheProcess/semver)
[![GitHub](https://img.shields.io/github/v/release/FollowTheProcess/semver?logo=github&sort=semver)](https://github.com/FollowTheProcess/semver)
[![CI](https://github.com/FollowTheProcess/semver/workflows/CI/badge.svg)](https://github.com/FollowTheProcess/semver/actions?query=workflow%3ACI)
[![codecov](https://codecov.io/gh/FollowTheProcess/semver/branch/main/graph/badge.svg?token=Q8Y5KFA9ZK)](https://codecov.io/gh/FollowTheProcess/semver)

Semver parsing and validation library for Go

* Free software: MIT License

## Project Description

Semver is a small, simple [semver] parsing and validation library for Go

## Installation

```shell
go get github.com/FollowTheProcess/semver@latest
```

## Quickstart

### Create a New Version

```go
version := semver.New(1, 2, 3, "rc.1", "build.123")
```

### Parse a Version from text

```go
version := semver.Parse("v1.6.12")
```

### Credits

This package was created with [copier] and the [FollowTheProcess/go_copier] project template.

[copier]: https://copier.readthedocs.io/en/stable/
[FollowTheProcess/go_copier]: https://github.com/FollowTheProcess/go_copier
[semver]: https://semver.org
