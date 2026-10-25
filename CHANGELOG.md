# Changelog

All notable changes to this project are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and the project uses
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Client for Hevo's API v2: Pipelines, Pipeline objects, jobs, Destinations and connectors, one
  builder per operation.
- Region or base URL, Basic auth, User-Agent and a cap on concurrent calls (default 4).
- Retries on `429` (honouring `Retry-After`, then `X-RateLimit-Reset`) for every method, and on
  `502`/`503`/`504` and network errors for `GET` and `DELETE`.
- `ResponseError` with the API's error code, and the rate-limit quota on a `429`.
