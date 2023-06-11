# Status Checker

## status-checker is a Go command line tool to acknowledge the health of the Slack product and get information when an incident, outage or maintenance occurs

## How to build the tool

To build the tool simply run the `make` command, this will generate binaries for various platforms in the `bin` folder.

```sh
├── bin
│   ├── status-checker-darwin-amd64
│   ├── status-checker-darwin-arm64
│   ├── status-checker-linux-amd64
│   └── status-checker-windows-amd64
```

## How to use status-checker

Currently the available commands are:

```
Usage of status:
  -current
        get Slack's current health status
  -debug-mode
        enable debug mode
  -history
        get Slack's history health status
```

### Examples

To query Slack's current health status, run the following command:

`./bin/status-checker-darwin-arm64 -current`

If there are no ongoing incidents this message should be returned:

```
🟢 Slack Current Health Status is Ok! 😄
```

To display the response as JSON add `-debug-mode`

**e.g**

`./bin/status-checker-darwin-arm64 -current -debug-mode`
