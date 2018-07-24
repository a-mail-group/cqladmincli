# cqladmincli
Tiny Cassandra Admin tool.

## Getting `cqladm`
```
go get github.com/a-mail-group/cqladmincli/cqladm
```

## Using `cqladm`

These commands are entered within the `cqladm` prompt.

### List up keyspaces
```
list keyspaces
```

### List up tables
```
list tables <<keyspace>>
```

### List up columns
```
list columns <<keyspace>>.<<table>>
```

### Run CQL commands.
```
do <<CQL-Command>>
```
### Run Multiline CQL commands.
```
do:
<<Multi
 -Line
 -CQL
 -Command>>;
```
