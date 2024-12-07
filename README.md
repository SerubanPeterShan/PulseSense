# PulseSense

PulseSense is a Go package project designed to retrieve CPU, disk, memory, and network data to monitor a system's performance and health. It provides a simple and efficient way to process and interpret system metrics for monitoring and analysis applications.

**Note: This is a pet project and I'm still figuring out Go, so there might be bugs when you implement it. Please proceed with caution.**

## Data Sources

PulseSense supports the following data sources:

- CPU usage
- Disk usage
- Memory usage

## Features

- Real-time data monitoring
- Ligthweight package

## Installation

To install PulseSense, use `go get`:

```sh
go get github.com/SerubanPeterShan/PulseSense
```

## Usage

Here is a basic example of how to use PulseSense:

```go
package main

import (
    "fmt"
    "github.com/SerubanPeterShan/PulseSense/Baremetal/cpusense/cpuusagesense"
)

func main() {
    cpuUsage := cpuusagesense.GetCPUUsage()
    fmt.Println("Current CPU usage:", cpuUsage)
}
```

## Documentation
