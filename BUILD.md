# Build Instructions

This document explains how to build the project locally.

## Prerequisites

* Go >= 1.21.6
* Git

## Clone the Repository

```bash
git clone https://github.com/Rajdeep-Nemo/sugarglaze
cd sugarglaze
```

## Build

```bash
# Build the project
go build ./cmd/glaze
```

## Run

```bash
# Run the compiled binary
./glaze
```

## Tests

```bash
# Run all tests
go test ./...
```

