# Gitlab branch tracker

[![Build Status](https://travis-ci.org/ns3777k/gitlab-branch-tracker.svg?branch=master)](https://travis-ci.org/ns3777k/gitlab-branch-tracker)

Simple utility to find and report left branches in Gitlab.

## Requirements

- Go 1.11+ (go modules required)

## Docker image

```bash
$ docker run --rm ns3777k/gitlab-branch-tracker
```

## Building

After cloning this repository simply run:

```bash
$ make build
```

## Mailhog

For testing purposes you can use [mailhog](https://hub.docker.com/r/mailhog/mailhog/):

```bash
$ docker run -p 1025:1025 -p 8025:8025 mailhog/mailhog
```

That'll run a smtp server on host's port `1025` and a web viewer to it on `8025`.
