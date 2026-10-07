# Eunuch
A scheme interpreter written in Zig.

> [!WARNING]
> This branch' API will *try* to follow the master's API, however, due to language constraints, some function signatures may change.

## USAGE
```
eunuch [file.scm]
```
### Versioning
#### Or: Why is this still in v0.X.X?
This project is a long-term commitment and won't receive a full release until the entire stdlib as-per the scheme lisp RFC-v7 has been implemented.
### Linters
#### Or: why are there so many linter errors?
The linters shouldn't be used. Its that simple.
The linters are only in the project to help catch dangling switches and forgotten field intializations (and they are still hit-or-miss about that too).