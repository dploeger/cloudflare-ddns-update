# cloudflare-ddns-update - Dyn v3 API that manages DNS entries in cloudflare

[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/cloudflare-ddns-update)](https://artifacthub.io/packages/search?repo=cloudflare-ddns-update)

## Introduction

The [Dyn v3 API](https://help.dyn.com/remote-access-api/perform-update/) to update dynamic DNS records is used in 
multiple routers on the market. This tool allows to provide this API while managing the dynamic DNS records in cloudflare.

## Configuration

The tool requires the following environment variables to be set:

* CLOUDFLARE_API_TOKEN: The API token used to authenticate against the cloudflare API
* CLOUDFLARE_ZONE: The name of the DNS zone to manage
* AUTH_USERNAME: A username that is required to be used for the Dyn API
* AUTH_PASSWORD: A password that is required to be used for the Dyn API

Optionally, DEBUG can be set to anything to enable logging of cloudflare requests and responses.

## Usage

After setting the environment variables, simply run the tool

    ./cloudflare-ddns-update

## Container image / Docker

This is also available as a container image, e.g. via Docker:

    docker run -e CLOUDFLARE_API_TOKEN=secret -e CLOUDFLARE_ZONE=example.com -e AUTH_USERNAME=testuser -e AUTH_PASSWORD=secret ghcr.io/dploeger/cloudflare-ddns-update:latest

## Helm chart

TODO