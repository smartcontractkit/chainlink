## Architecture

```
                                                                  Schema
                                                                     |
                                                                     v
                                    /-> KafkaExporter = Mapper -> Encoder -> Producer
ChainReader -> Poller -> Exporter(s)
                                    \-> PrometheusExporter
\----------------+-----------------/
                 |
                 v
            FeedMonitor
                 |
                 v
          MultiFeedMonitor
                 |
                 v
   RddPoller-->Manager

```

## Abstractions

Don't create abstractions for the sake of it.
Always ask "What is the simplest thing that can solve the problem?"
Always ask "Would this code make sense if I didn't write it but had to change it?"
As a rule of thumb, abstraction around IO - think database or http connections - are almost always a good idea, because they make testing easier.
Another rule of thumb, is to compose abstractions by means of dependency injection.

## Concurrency

Concurrency is hard.
To make it manageable, always use established concurrent patterns specific to golang. Eg. https://go.dev/blog/pipelines or https://talks.golang.org/2012/concurrency.slide
Have all the concurrent code in one place and extract everything else in functions or interfaces.
This will make testing the concurrent code easier - but still not easy!

A tenant of good concurrent code is resource management.
Your primary resources are goroutines, channels and contexts (effectively wrappers ontop of channels).
Make sure, that upon exit, your program cleanly terminates all goroutine, releases all OS resources (file pointers, sockets, etc.), no channel is being used anymore and all contexts are cancelled.
This will force you to manage these things actively in your code and - hopefully - prevent leaks.

## Logging

I have yet to find an engineer who like GBs of logging.
Useless logs have a cognitive load on the person trying to solve an issue.
My approach is to log as little as possible, but when you do log, put all the data needed to reproduce the issue and fix it in the log!
Logging takes time to tune. Try to trigger or simulate errors in development and see if the log line is useful for debugging!

## Testing

This is controversial but I'm not a huge fan of testing as much as possible.
Most tests I've seen - and written - are brittle, are non-deterministic - ofc they break only in CI - and are not very valuable.
The most valuable test is an end-to-end test that checks a use-case.
The lest valuable test is a unit test that tests the implementation of a function.

Another thing that makes writing valuable tests easier is good "interfaces".
If a piece of code has clearly defined inputs and outputs, it's easier to test.

## Errors

An often overlooked output of a piece of code are it's errors. It's easy to `return nil, err`!
Well defined errors can be either public values or error types - when more context is needed to trace the error.
Make sure you consider whether a specific error can be handled by the caller or needs to be pushed up the stack!

## Benchmarks

Execute the existing benchmark whenever a significant change to the system is introduced.
While these benchmarks run in an ideal situation - eg. 0 network latency, correctly formatted message, etc. -
they give a reference point for potential performance degradation introduced by new features.

Benchmarks are - arguably - the easiest way to profile your code!
