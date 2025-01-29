# Kerosin

**Kerosin** is a versatile converter designed to transform OpenBullet configuration components into C# and Golang code. It streamlines the process of migrating or integrating OpenBullet configs into your Go or C# projects, supporting various modules like HTTP Requests, Key Checks, and Parsers.

## Features

- **Module Conversion**: Convert OpenBullet modules to Go/C# code:
  - **Request**: Supports `GET` and `POST` HTTP methods.
  - **KeyCheck**: Validate responses against predefined conditions.
  - **Parser**: Convert `LR` (Linear Regression) and `JSON` parsers.
- **Dynamic Function Naming**: Generates context-aware function names based on URL domains (e.g., `createRequestForExample` for `example.com`).
- **Cross-Language Support**: Outputs idiomatic code for both Go and C#.

## Installation

1. Ensure **Go 1.20+** is installed on your system.
2. Clone this repository:
   ```bash
   git clone https://github.com/Axion-Security/Kerosin.git
   cd Kerosin
   ```
3. Build the project:
   ```bash
   go build -o kerosin
   ```

## Contributing

Contributions are welcome! Fork the repository, create a branch, and submit a pull request. Report bugs via GitHub Issues.

## License

Distributed under the MIT License. See `LICENSE` for details.
