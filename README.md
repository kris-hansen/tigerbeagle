# TigerBeagle

![TigerBeagle Logo](assets/logo.png)

TigerBeagle is a CLI tool designed to simplify account and transaction data management for [TigerBeetle](https://github.com/tigerbeetle/tigerbeetle), the high-performance distributed database engine.

Project Goals:

- Be a tool to help bridge between TigerBeetle and systems of transactional reference (i.e., transactional metadata)
- Faciliate mass account / transaction manipulation for conversion, migration and testing
- Generate test data for ledger testing

## Table of Contents

- [TigerBeagle](#tigerbeagle)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Installation](#installation)
  - [Usage](#usage)
  - [Commands](#commands)
  - [Migration Guide](#migration-guide)
  - [Contributing](#contributing)
  - [License](#license)
  - [Acknowledgments](#acknowledgments)

## Features

- Create and manage accounts in TigerBeetle
- Perform single and bulk transfers
- Migrate accounts and transactions from JSON files
- Validate connectivity to TigerBeetle
- Simplify testing and development workflows

## Installation

To install TigerBeagle, you need to have Go 1.16 or later installed on your system. Then, you can install it using the following command:

```bash
go install github.com/yourusername/tigerbeagle@latest
```

Alternatively, you can clone the repository and build it manually:

```bash
git clone https://github.com/yourusername/tigerbeagle.git
cd tigerbeagle
go build
```

## Usage

Before using TigerBeagle, make sure you have a TigerBeetle server running. You can set the TigerBeetle address using the `TB_ADDRESS` environment variable or the `--tb-address` flag.

```bash
export TB_ADDRESS=3000
tigerbeagle [command]
```

or

```bash
tigerbeagle --tb-address=3000 [command]
```

## Commands

- `create-account`: Create a new account
- `get-account`: Get account details
- `transfer`: Perform a transfer between accounts
- `bulk-transfer`: Perform multiple transfers in bulk
- `migrate-accounts`: Migrate accounts from a JSON file
- `migrate-transfers`: Migrate transfers from a JSON file
- `doctor`: Validate connectivity to TigerBeetle

For detailed information on each command, use the `--help` flag:

```bash
tigerbeagle [command] --help
```

## Migration Guide

For detailed instructions on migrating accounts and transfers, please refer to our [Migration Guide](MIGRATE.md).

## Contributing

We welcome contributions to TigerBeagle! Please feel free to submit issues, fork the repository and send pull requests!

1. Fork the project
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [TigerBeetle](https://github.com/tigerbeetle/tigerbeetle) - The high-performance accounting database that inspired this project
- All the contributors who have helped shape TigerBeagle

---

For more information, bug reports, or feature requests, please [open an issue](https://github.com/yourusername/tigerbeagle/issues) on GitHub.
