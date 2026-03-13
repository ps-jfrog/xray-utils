# Release Notes

## [0.1.0] - Initial release

### Added

- **set-license-priority** command (alias: **slp**) to set license name priorities via the Xray API.
- POST request to `/api/v1/licensesNames/priorities` with parameters:
  - `license-name`: name of the license
  - `license-priority`: priority value (integer)
- Support for environment variables: `JFROG_URL` / `JF_URL`, `JFROG_ACCESS_TOKEN` / `JF_ACCESS_TOKEN`.
