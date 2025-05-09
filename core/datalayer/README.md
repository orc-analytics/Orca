# Datalayers

## Logic Flow

### Registering a Processor

When a processor registers with the central orca server, the following
needs to occur:

Begin transaction:

1. Create any window types that do not exist
2. Create algorithms with their window types and associated processor
3. Create any algorithm dependencies
4. Register the processor and remove old algorithm associations to that processor
5. Associate the processor with it's supported algorithms

If there are any errors in the process (e.g. cyclic dependencies on algorithms)
then rollback the transaction.

## Rules

1. Only DAGs can be registered
2. DAG refers to the window type, processor type and algorithm type
