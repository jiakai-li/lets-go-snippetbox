# Let's go

*This repo is a note from reading [Let's go](https://lets-go.alexedwards.net/)*

## Chapter 2.7

- `cmd` directory: *application-specific* code for the executable applications in the project.
- `internal` directory: *non-application-specific* code used in the project.
  - packages live under this directory can only be imported by code inside the parent of the `internal` directory
  - meaning:
    - any packages which live in `internal` can only be imported by code inside `lets-go-snippetbox` project directory
    - any packages under `internal` cannot be imported by code outside of this project
- `ui` directory: *user-interface assets* used by the web application.
