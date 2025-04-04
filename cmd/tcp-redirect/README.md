# TCP Port Forwarding Server

## Overview

Asis is a simple TCP port forwarding server written in Go. It listens on a specified local address and forwards incoming connections to a remote address.

## Usage

```bash
go run main.go -listen :10001 -remotremote remote-machine:22
```

### Command Line Flags

- `-listen` (string)
   - Local address to listen on
   - Default: `:10001`
   - Example: `:5112`

- `-remote` (string)
   - Remote address to forward connections to
   - Required parameter
   - Example: `remote-machine:22`

## Implementation Details

### Main Components

1. **Main Function**
   - Parses command line flags
   - Sets up TCP listener
   - Accepts incoming connections and handles them in goroutines

2. **dealClientConn Function**
   - Handles individual client connections
   - Establishes connection to remote server
   - Sets up bidirectional data transfer
   - Manages connection cleanup

### Key Features

- JTCP Connection Handling**
   - Resolves both local and remote TCP addresses
   - Supports IPv4 connections
   - Implements keep-alive functionality

- **Data Transfer**
   - Uses `io
 CopyBuffer` for efficient data transfer
   - Bidirectional forwarding between local and remote connections
   - Buffered data transfer

- **Error Handling**
   - Comprehensive error checking
   - Logging of connection issues
   - Graceful connection cleanup

- **Concurrency**
   - Uses goroutines for
     Handling multiple connections
   - Implements WaitGroup for proper connection cleanup

### Configuration

```go
// Default listen address
flag.StringVar(&listen, "listen", ":10001", "listen address")

// Remote address (required)
flag.StringVar(&remoteAddr, "remote", "", "remote address")
````


### Connection Settings

- Keep-alive enabled on both local and remote connections
- Keep-alive period set to 10 seconds
- Automatic connection cleanup on completion or error

## Error Messages

-&no listen address" - When listen address is empty
- "no remote address" - When remote address is not specified
- Various connection-specific error messages with connection details

## Logging

The program logs:
- Connection establishment failures
- Data transfer errors
- Connection cleanup events
- Keep-alive configuration issues

All log messages include connection information showing local and remote endpoints.

## Dependencies

- Standard Go libraries:
   - `flag`
   - `fmt`
   - `io`
   - `log`
   - net`
   - sync`
   - time`
