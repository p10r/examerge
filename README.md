## Examerge

Merges exams with ratings while preserving the directory structure to be sent back to students.

Expected directory structure:

```
|-- example
|   |-- student1
|   |   |-- example_exam1.pdf
|   |   `-- example_rating1.pdf
|   `-- student2
|       |-- example_exam2.pdf
|       `-- example_rating2.pdf
```

The script differentiates between an exam and its rating by checking a specific prefix that can
configured in `main.go`.

## Development

### Commit Hooks

This project uses [lefthook](https://github.com/evilmartians/lefthook)
and [golangci-lint](https://golangci-lint.run/). Install via:

```shell
brew install lefthook golangci-lint
```

### Linting

Make sure to run `golangci-lint run`

## Release

Examerge is currently only being run on Windows. It's build on every push via a hook.

See [lefthook.yml](./lefthook.yml) for details.