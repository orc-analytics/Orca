# Orca

Orca is a analytics orchestration framework that makes it easy for development and product teams to
extract insights from timeseries data. It provides a structured and scalable way to schedule, process,
and analyse data using a time window based triggering mechanism and a flexible DAG-based
architecture. All of this combined makes it seamless to tweak the Cost <-> Availability <-> Accuracy
tradeoff that is always present in timeseries processing.

Orca is built by [Predixus](https://www.predixus.com) for developers that are ready to get their timeseries
projects to production 🚀.

## ✨ Features

- **Pluggable processors:** Register cross-language processors dynamically and scale horizontally
- **Window-based execution:** Define regions of interest (windows) to trigger algorithms
- **Execution engine:** Automatically handles algorithm dependencies and execution order,
  without you having to worry about the [DAG](https://en.wikipedia.org/wiki/Directed_acyclic_graph).
- **Abstracted storage layer:** Growing list of storage solutions that work with Orca. Currently supported databases:

  - PostgresSQL

    With the following in the works:

  - MongoDB
  - BigQuery
  - RDS

## 🚀 Getting Started

### 1. Install the Orca cli

Prior to installing the CLI ensure that Docker is installed on your system.

#### Linux / MacOSX

```bash
curl -fsSL https://raw.githubusercontent.com/Predixus/orca/main/install-cli.sh | bash
```

#### Windows

Orca heavily leverages dockerised systems as part of the Orca stack. These generally work better on unix
based systems, so it's advised to install Orca on WSL. Once in WSL, use the above CLI install script.

## 📦 Architecture

1. Processors register with the Orca Core service.
2. Windows are emitted into the system.
3. Orca builds an execution DAG from dependencies.
4. Tasks are streamed to processors.
5. Results return to the core.
6. Dependent algorithms are triggered automatically.

---

## 🔌 Extending Orca

To implement a custom datalayer, your driver must implement:

- `CreateProcessor`
- `EmitWindow`

See `internal/datalayers/types.go` for the interface.

---

## 💻 Development

### Install proto tools

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
python -m pip install grpcio grpcio-tools
```

### Generate docs (optional)

```bash
docker pull pseudomuto/protoc-gen-doc
```

---

## 📜 Rules

1. Algorithm DAGs can only be triggered by a single WindowType.
2. Algorithms can’t depend on algorithms from a different WindowType.

---

## 💬 Community

- Issues: [GitHub Issues](https://github.com/predixus/orca/issues)
- Discussions: Coming soon!

---

## 📄 License

Orca is licensed under the [Business Source License (BSL) 1.1](./LICENSE.md).

- Free for companies under £5 million total value (including production use).
- Free for trial and evaluation by companies over £5 million.
- Free for registered charities and educational institutions.
- No one may build competing software on top of Orca.

**Change Date:**
Each Orca version will automatically become open source under the GPLv3 license four (4) years after its first public release.

For full license terms, see [LICENSE.md](./LICENSE.md).
