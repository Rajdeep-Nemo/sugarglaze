# Build Instructions

This document explains how to build and run the project locally.

---

## Prerequisites

| Requirement | Version |
|-------------|---------|
| Go | >= 1.21.6 |
| Git | any recent version |

---

## Clone the Repository

```bash
git clone https://github.com/Rajdeep-Nemo/sugarglaze
cd sugarglaze
```

---

## Build

```bash
go build ./cmd/glaze
```

---

## Run

```bash
./glaze
```

---

## Tests

```bash
go test ./...
```