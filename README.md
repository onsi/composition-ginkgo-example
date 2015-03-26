# Composition Ginkgo Example

This example showcases one way to write reusable [Ginkgo](https://github.com/onsi/ginkgo) tests.  The intent is to present patterns that illustrate Ginkgo's flexibility - this is not the typical way to write Ginkgo tests and should only be considered if you want to be able to reuse tests in different contexts.  It's primary intent is to fuel the discussion around [#145](https://github.com/onsi/ginkgo/issues/146)

## Organization

The program under test is a hokey key-value store.  The code for the key-value store is in the `key_value_store` package.

There are *two* sets of tests that excercise the key-value store.  `tests/key_value_tests` cover the basic functionality of the store.  `tets/prefix_tests` cover the get-prefix and delete-prefix features (the code is pretty self-explanatory here).  These test packages are not traditional Ginkgo tests.  In fact they aren't strictly speaking Go `test` packages at all.  They are simply Go packages that construct a number of Ginkgo tests using the Ginkgo DSL.

These tests are written in such a way that they must be pointed at a running key-value store.  You do this by passing the test packages a `SharedContext` (defined in the `helpers` package).  This `SharedContext` includes a preconfigured client that points at the running key-value store.  It's important to understand that the tests under `/tests` **do not** set up this context or spin up the key-value store.  That responsibility is left to the Ginkgo suite that imports the test packages.

There are two such suites: 

1. `/integration` is a basic integration suite that:
    - compiles the key-value store
    - launches the key-value store binary
    - imports the two `tests` packages
    - configures the `SharedContext` and passes it into the `tests` packages
    - takes care of cleaning up the key-value store between test runs
    - tears down the key-value store at the end of the test.

2. `/stress` does everything the integration suite does but also:
    - runs a chaos-monkey style goroutine that messes with the key-value store (it's called EntropyOrangutan in the code)
    - presents a pattern for intercepting failures and acting upon them (this could be used, for example, to emit a metric or pause the EntropyOrangutan so that an operator can investigate the failure)

Because of how things are structured, `integration` and `stress` are two *different* suites that can actually share the *same tests*.

The other major thing this example illustrates is a pattern for writing parallelizable tests.  Ginkgo has strong first-class support for running test suites in parallel.  When using a shared resource (e.g. the key-value store), however, the onus is left on the developer to ensure the concurrently running test nodes do not interfere with each other.

There are two ways to do this.  A trivial approach would be to spin up a key-value store for each parallel test node.  Sometimes, however, this is not possible (for example, the object under test might be an external resource like a cluster).  In this case carefully sharding your tests is the preferred approach.

These tests opt for the latter approach to illustrate how it might be done.  Each parallel test node is given a key prefix which effectively shards the database across nodes.

We've simulated having slow integration tests by introducing a random sleep into the (otherwise fast) key-value store.  You can see the difference by running `ginkgo` vs `ginkgo -p`

## Trying it out

This assumes you already have the `ginkgo` cli installed.

```
go get github.com/onsi/composition-ginkgo-example
cd $GOPATH/src/github.com/onsi/composition-ginkgo-example
```

Now you can run the tests.  To run all the tests in series:

```
ginkgo -r
```

(note that the stress test may fail -- the EntropyOrangutan can be ruthless!)

To run all the tests in parallel:

```
ginkgo -r -p
```

To run just the stress tests until they fail:
```
ginkgo -untilItFails -p stress
```

etc...