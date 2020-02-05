# sensu-wavefront-handler

## Table of Contents
- [Overview](#overview)
- [Usage Examples](#usage-examples)
- [Configuration](#configuration)
  - [Asset registration](#asset-registration)
  - [Resource definition](#resource-definition)
- [Installation from source](#installation-from-source)
- [Contributing](#contributing)

## Overview
The Sensu Wavefront Handler is a [Sensu Event Handler][9] that sends metrics to the SaaS time series
database [Wavefront][10]. [Sensu][11] can collect metrics using check output metric extraction or
the StatsD listener. Those collected metrics pass through the event pipeline, allowing Sensu to
deliver normalized metrics to the configured metric event handlers. This Wavefront handler will
allow you to store, instrument, and visualize the metric data from Sensu.

## Usage Examples

Help:
```
a wavefront metrics handler built for use with sensu

Usage:
  sensu-wavefront-handler [flags]
  sensu-wavefront-handler [command]

Available Commands:
  help        Help about any command
  version     Print the version number of this plugin

Flags:
  -h, --help            help for sensu-wavefront-handler
  -s, --server string   the address of the wavefront server
  -t, --token string    the API token for the wavefront server
```

## Configuration

### Asset registration

Assets are the best way to make use of this plugin. If you're not using an asset, please consider
doing so! If you're using sensuctl 5.13 with Sensu Backend 5.13 or later, you can use the following
command to add the asset:

```
sensuctl asset add sensu/sensu-wavefront-handler
```

If you're using an earlier version of sensuctl, you can find the asset on the [Bonsai Asset Index]([https://bonsai.sensu.io/assets/sensu/sensu-wavefront-handler](https://bonsai.sensu.io/assets/sensu/sensu-wavefront-handler)).

### Resource definition

```yml
---
type: Handler
api_version: core/v2
metadata:
  name: sensu-wavefront-handler
  namespace: default
spec:
  command: sensu-wavefront-handler --server https://beachops.wavefront.com --token d3b33111-0000-4405-9999-fad497082222
  runtime_assets:
  - sensu-wavefront-handler
  type: pipe
```

## Installation from source

The preferred way of installing and deploying this plugin is to use it as an Asset. If you would
like to compile and install the plugin from source or contribute to it, download the latest version
or create an executable script from this source.

From the local path of the sensu-wavefront-handler repository:

```
go build -o /usr/local/bin/sensu-wavefront-handler main.go
```

## Contributing

For more information about contributing to this plugin, see [Contributing][1].

[1]: https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md
[2]: github.com/sensu-community/sensu-plugin-sdk
[3]: https://github.com/sensu-plugins/community/blob/master/PLUGIN_STYLEGUIDE.md
[4]: https://github.com/sensu/sensu-wavefront-handler/blob/master/.github/workflows/release.yml
[5]: https://github.com/sensu/sensu-wavefront-handler/actions
[6]: https://github.com/sensu/sensu-wavefront-handler/releases
[7]: https://github.com/sensu/sensu-wavefront-handler/blob/master/main.go
[8]: https://bonsai.sensu.io/
[9]: https://docs.sensu.io/sensu-go/latest/reference/handlers/#how-do-sensu-handlers-work
[10]: https://www.wavefront.com/
[11]: https://github.com/sensu/sensu-go