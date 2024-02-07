# Test Suites

The `tests/` directory is reserved for non-unit test suites. It is
declared as a go module so that test execution is explicit and not
included when running `go test ./...` from the root of the repo.

Different test suites may have different setup required. It is therefore
invalid to run `go test ./...` from the `tests/` directory. Specific
test suites should be invoked as `go tests ./<test-suite>`, where
`<test-suite>` corresponds to a subpackage, e.g., `go test ./e2e`.
