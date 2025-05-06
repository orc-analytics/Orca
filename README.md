# Orca

Orca is a analytics orchestration framework that makes it easy for development and product teams to
extract insights from timeseries data. It provides a structured and scalable way to schedule, process,
and analyse data using a time window based triggering mechanism and a flexible DAG-based
architecture. All of this combined makes it seamless to tweak the Cost <-> Availability <-> Accuracy
tradeoff that is always present in timeseries processing.

Orca is built by (Predixus)[https://www.predixus.com] for developers, ready to deploy their timeseries
projects 🚀.

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

### 1. Install Orca

Clone the repo:

```bash
git clone https://github.com/predixus/orca.git
cd orca
```

Build the binary:

```bash
make build
```

### 2. Setup Database

Start a local PostgreSQL instance with:

```bash
make build_store
```

Other DB commands:

```bash
make start_store    # start DB
make stop_store     # stop DB
make remove_store   # delete DB and data
make redo_store     # reset DB
```

### 3. Run Orca Core

```bash
./orca --connStr "postgresql://orca:orca_password@localhost:5432/orca?sslmode=disable"  --platform postgresql --port 3335 --migrate
```

The `--migrate` flag will instruct orca to provision the schemas within the store.

### 4. Register a Processor

Use gRPC or a client library to register processors and algorithms. Processors should implement the `OrcaProcessor` gRPC interface.

Current processor SDKs:

- (Python)[https://www.github.com/Predixus/orca-python.git]

---

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
