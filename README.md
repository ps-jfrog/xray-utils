# xray-utils

## About this plugin

This plugin extends JFrog CLI with utilities for JFrog Xray, including setting license name priorities via the Xray API.

## Installation with JFrog CLI

Installing the latest version:

```bash
$ jf plugin install xray-utils
```

Installing a specific version:

```bash
$ jf plugin install xray-utils@version
```

Uninstalling the plugin:

```bash
$ jf plugin uninstall xray-utils
```

## Usage

### Commands

#### set-license-priority (alias: slp)

Sets the priority for a license name by sending a POST request to the Xray API endpoint `/api/v1/licensesNames/priorities`.

- **Arguments:**
  - **license-name** – The name of the license to set priority for.
  - **license-priority** – The priority value (e.g. 1, 2, 3).

- **Flags:**
  - **--server-id** – (Optional) Server ID from JFrog CLI config. If omitted, the default server from `jf c add` is used.

- **Examples:**

  Using JFrog CLI config (recommended; run `jf c add` first to configure URL and token):

  ```bash
  $ jf xray-utils set-license-priority "MIT" 1
  $ jf xray-utils set-license-priority "MIT" 1 --server-id my-server
  ```

  Or using the alias:

  ```bash
  $ jf xray-utils slp "Apache-2.0" 2
  ```

### Configuration and environment variables

The plugin resolves the JFrog URL and access token in this order:

1. **JFrog CLI config** – If you have run `jf c add`, the plugin uses the default server’s platform URL and access token. Use the optional `--server-id` flag to select a non-default server.
2. **Environment variables** – If config is missing or incomplete, the plugin falls back to:
   - **JFROG_URL** or **JF_URL** – JFrog platform URL (e.g. `https://your-instance.jfrog.io`).
   - **JFROG_ACCESS_TOKEN** or **JF_ACCESS_TOKEN** – Access token for authentication.

So you can either configure once with `jf c add` and run without env vars, or set the env vars (e.g. in CI) and omit config.

## Building from source

```bash
go build -o xray-utils .
```

The resulting binary can be used as a JFrog CLI plugin when placed in the JFrog CLI plugins directory.

## Building an executable for testing

To build and test the plugin locally with JFrog CLI before publishing:

1. **Build the plugin executable** (from the project root):

   ```bash
   go build -o xray-utils .
   ```

2. **Install it into the JFrog CLI plugins directory** so `jf` can discover it. By default, plugins live under the JFrog CLI home (e.g. `~/.jfrog/plugins/`). Create the plugin directory and copy the binary:

   ```bash
   mkdir -p ~/.jfrog/plugins/xray-utils/bin
   cp xray-utils ~/.jfrog/plugins/xray-utils/bin/
   ```

3. **Verify** that the plugin is available:

   ```bash
   jf
   ```

   You should see `xray-utils` listed. Then run the plugin:

   ```bash
   jf xray-utils set-license-priority "MIT" 1
   ```

**Tip: avoid copying after every build** — Use a symlink so `jf` always runs the binary in your project directory. After the first build, run:

```bash
ln -sf "$(pwd)/xray-utils" ~/.jfrog/plugins/xray-utils/bin/xray-utils
```

Then you only need to run `go build -o xray-utils .` when you change code; no need to copy the binary again. You can also run the plugin without installing via `jf` by using:

```bash
go run . -- set-license-priority "MIT" 1
```

## Publishing to a private registry

You can publish this plugin to a **private JFrog CLI Plugins Registry** hosted on Artifactory. See the [JFrog CLI Plugins Developer Guide](https://github.com/jfrog/documentation/blob/main/jfrog-applications/jfrog-cli/cli-plugins/developer-guide.md) for full details.

### 1. Set up the private registry on Artifactory

- Create a **local generic** repository on your Artifactory server named `jfrog-cli-plugins` (or another name of your choice).
- Ensure your Artifactory server is configured in JFrog CLI:

  ```bash
  jf c show
  ```

  If needed, add it:

  ```bash
  jf c add
  ```

- Set the ID of that configured server as the **JFROG_CLI_PLUGINS_SERVER** environment variable:

  ```bash
  export JFROG_CLI_PLUGINS_SERVER=your-server-id
  ```

- If your repository name is not `jfrog-cli-plugins`, set it with **JFROG_CLI_PLUGINS_REPO**:

  ```bash
  export JFROG_CLI_PLUGINS_REPO=your-plugins-repo-name
  ```

### 2. Publish the plugin

From the **root of the plugin source directory**, run:

```bash
jf plugin publish xray-utils <version>
```

Example:

```bash
jf plugin publish xray-utils v0.1.0
```

This builds the plugin for all supported operating systems and uploads the binaries to your private registry. Users can then install it with `jf plugin install xray-utils` (and optionally `@version`) as long as their JFrog CLI is configured to use that private plugins server (via **JFROG_CLI_PLUGINS_SERVER**).

## License

This project is licensed under the Apache License 2.0. See [LICENSE](LICENSE) for the full text.

## Additional info

None.

## Release Notes

The release notes are available [here](RELEASE.md).
