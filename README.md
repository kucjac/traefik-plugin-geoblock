# traefik-plugin-geoblock

[![Build Status](https://github.com/kucjac/traefik-plugin-geoblock/actions/workflows/ci.yml/badge.svg)](https://github.com/kucjac/traefik-plugin-geoblock/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/kucjac/traefik-plugin-geoblock)](https://goreportcard.com/report/github.com/kucjac/traefik-plugin-geoblock)
[![Latest GitHub release](https://img.shields.io/github/v/release/kucjac/traefik-plugin-geoblock?sort=semver)](https://github.com/kucjac/traefik-plugin-geoblock/releases/latest)
[![License](https://img.shields.io/badge/license-Apache%202.0-brightgreen.svg)](LICENSE)

*traefik-plugin-geoblock is a traefik plugin to whitelist requests based on geolocation*

> This projects includes IP2Location LITE data available from [`lite.ip2location.com`](https://lite.ip2location.com/database/ip-country).

## Configuration

### Static

#### Local

```yaml
experimental:
  localPlugins:
    geoblock:
      moduleName: github.com/kucjac/traefik-plugin-geoblock
```

#### Pilot

```yaml
pilot:
  token: "xxxxxxxxx"

experimental:
  plugins:
    geoblock:
      moduleName: github.com/kucjac/traefik-plugin-geoblock
      version: v0.2.0
```

### Dynamic

```yaml
http:
  middlewares:
    geoblock:
      plugin:
        geoblock:
          # Whether or not to enable geoblocking.
          enabled: true
          # Path to the ip2location database.
          databaseFilePath: /plugins-local/src/github.com/kucjac/traefik-plugin-geoblock/IP2LOCATION-LITE-DB1.IPV6.BIN
          # Countries to allow requests from, using ISO 3166-1 alpha-2 codes.
          # See https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2#Officially_assigned_code_elements
          # Either allowedCountries or disallowedCountries could be set at once.
          allowedCountries: [ "AT", "CH", "DE" ]
          # Countries from which requests are forbidden, using ISO 3166-1 alpha-2 codes.
          # Either allowedCountries or disallowedCountries could be set at once.
          disallowedCountries: [ "AT", "CH", "PL" ]

          # Whether or not requests from private networks should be allowed.
          allowPrivate: true
```