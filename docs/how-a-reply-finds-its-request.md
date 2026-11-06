# How a reply finds its request

A command is answered synchronously: the http request that sent it waits until
the connector replies or the timeout ends. Everything between sending and
answering hangs on one value, the task id of the protocol message, and on two
naming rules that keep foreign replies out.

## Scope

Applies from `v1.2.1` to the command path of both variants, kafka in the cloud
(`com_impl` `cloud`) and the mqtt broker of a gateway (`com_impl` `mgw`), for
commands that go to a connector.

Not about measuring functions of event services. Those are answered from the
last value (`pkg/lastvalue`) and never wait for a reply.

## The waiting request owns the reply

`deviceCommand` in `pkg/command/devicecommand.go` registers the task id with
`Replies` **before** it sends, so a reply that arrives faster than the send
returns is not lost. The transport (`pkg/transport`) passes every reply it reads
to `Replies.Receive` as `transport.Reply{TaskId, Output, Failed}` and does
nothing else with it.

`Receive` hands the reply to the request waiting for that task id and forgets
the id. A reply for an id nobody waits for is dropped: a second reply for the
same command, a reply after the timeout, a reply to another client. It never
blocks the transport, so one stuck request cannot hold back the replies of
others.

The request unmarshals the output itself, with the service, protocol, output
characteristic, function and aspect nodes it sent. The metadata a connector
echoes in its reply is not read, so a connector that changes or shortens it —
for example one that knows only the deprecated single aspect node — does not
change the result.

The result of one device command; `/commands` for a single device wraps a 200
result in a list, a group command collects them.

| Outcome | Status | Result |
|---|---|---|
| reply, no output characteristic | 200 | `null` |
| reply, unmarshalled | 200 | the value |
| reply fails to unmarshal | 500 | the error |
| connector reported an error (`Failed`) | 500 | the output map of the error |
| send failed | 500 | `unable to produce message` |
| no reply within the timeout | 408 | `timeout` |

## Kafka: every instance reads its own reply topics

The message names where to reply in `metadata.response_to` and
`metadata.error_to`. Both are the configured topics with
`topic_suffix_for_scaling` appended, and the consumer reads exactly those
(`Config.WithScalingSuffix`, applied once in `pkg.Start`). The deployment sets
the suffix to the pod name.

The suffix is what makes more than one instance work. All instances share the
consumer group `kafka_consumer_group`; on a shared topic the group splits the
partitions between them, and a reply then reaches an instance that does not
know its task id, drops it, and the sender runs into its timeout. With a single
instance a missing suffix goes unnoticed, which is how it was lost once in a
refactoring. `TestCommandScaled` therefore checks the reply topics of the sent
message, not only that the command succeeds.

## Mqtt: the command id carries a prefix

On a gateway the connector echoes the command id, and every client of the broker
receives every response on `response/#` and every error on `error/command/#`.
The automation-engine sends commands to the same connectors. The command id is
`mgw_correlation_id_prefix` followed by the task id; `pkg/transport/mqtt.go`
reads only replies with that prefix and strips it to get the task id back. An
empty prefix is refused at start, because then every foreign reply would be read
as a candidate.

The prefix needs no correlation table: the task id travels inside the command id
and comes back with the reply.

## Where to look

- `pkg/command/replies.go` and `replies_test.go` — registration, delivery, dropping
- `pkg/transport/kafka.go`, `pkg/transport/mqtt.go`, `mqtt_test.go` — reading replies, the prefix check
- `pkg/configuration/config.go`, `WithScalingSuffix`
- `pkg/tests/command_test.go`, `testCommand` — the reply topic check of the scaled test
