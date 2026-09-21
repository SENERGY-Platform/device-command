# Aspects arrive as a list

A command names a set of aspects, not one aspect. The single-valued fields that
predate the list still exist, are marked deprecated, and behave as an alias for
a list with one element.

## Scope

Applies from `device-repository/v2 v2.2.1`, `marshaller v0.2.0`,
`external-task-worker v0.1.1` and
`models/go v0.0.0-20260911075423-f01521c01da2` onwards, to the requests this
service accepts on `/commands` and `/commands/batch` and to the protocol message
it writes.

Not about the aspects of a **content variable**. Those are `aspect_ids` on the
device type and say what a variable carries; the fields here say what a command
demands. The two are named similarly and point in opposite directions.

The rule itself — what several aspects mean, and how a variable satisfies them —
is not this repository's. It belongs to the device-repository and is implemented
in the marshaller's `lib/marshaller/model/aspects.go`, which this service calls
rather than reimplements. The same boundary is written down on the worker side
in the external-task-worker repository, `docs/aspects-arrive-as-a-list.md`.

## The fold happens at the two entry points

`Command()` in `pkg/command/command.go` and `Batch()` in `pkg/command/batch.go`
are the only ways into the command package, and both call `SetAspectIds()` on
the decoded request before anything else. From there inward nothing reads the
deprecated field: `DeviceCommand`, `GroupCommand`, `GetSubTasks` and
`deviceCommand` all take `aspectIds []string`.

`SetAspectIds` **clears** `AspectId` after folding it into `AspectIds`. That
matters for `Batch`: it deduplicates elements by `CommandMessage.Hash()`, which
hashes the encoded json, so two elements naming the same aspect — one in the
deprecated field, one in the list — only collapse into a single command if the
deprecated field is gone by then.

## The fields and their accessors

| Type | Deprecated | List | Accessors |
|---|---|---|---|
| `command.CommandMessage` | `aspect_id` | `aspect_ids` | `GetAspectIds()`, `SetAspectIds()` |
| `command.BatchRequest` | — | — | `SetAspectIds()` folds every element |
| `command.SubCommand` | — | `AspectIds` | — |
| `messages.Metadata` (external-task-worker) | `output_aspect_node` | `output_aspect_nodes` | `GetOutputAspectNodes()`, `SetOutputAspectNodes()` |

`SubCommand` is in-process only and never serialized, so it carries the list
alone.

The accessors in `pkg/command/aspects.go` are built on the marshaller's own
`AspectIdsAlias`, not on a local copy, so the fold behaves the same here as in
every other caller of that model.

## Writing keeps the deprecated node filled

`deviceCommand` calls `Metadata.SetOutputAspectNodes()`, which writes the list
**and** the deprecated single node, so a protocol handler that only knows the
old field still gets an answer. The node it picks is the one with the
alphabetically first id, the way the device-repository fills the deprecated
field of a path option. That is a silent narrowing for such a reader: it sees
one node out of several, chosen by how the urns sort rather than by which aspect
matters.

`HandleTaskResponse` reads the answer back with `GetOutputAspectNodes()`, so a
response to a command that was sent before the lists existed still unmarshals.

## Device group filtering

`getFilteredServices` in `pkg/command/groupcommand.go` asks, per content
variable, whether it satisfies the whole set:

```go
if variable.FunctionId == functionId &&
	marshallermodel.AspectMatchLevel(marshallermodel.ContentVariableAspectIds(variable), aspectNodes) >= 0 {
	return true
}
```

`AspectMatchLevel` returns `-1` as soon as one queried aspect is unmatched, so
several aspects are an AND on **one** variable, and each covers its own subtree.
It reads the subtree from `ChildIds` and `DescendentIds` of the queried node, so
the nodes have to come from the device-repository. `GetSubTasks` falls back to a
node carrying only the id when the lookup fails; such a node matches nothing but
itself, which is the behavior that was there before the lists.

## The input aspect now reaches the marshaller

Before the lists, the `MarshallingV2RequestData` of a controlling function was
built with an aspect node that was assigned one line *after* the struct literal,
so it was always nil and the input path was chosen by function id alone. The
list version passes the resolved nodes, which is what the external-task-worker
does on the same path. A controlling command that names an aspect whose subtree
no input variable serves therefore fails to marshal now, where it previously
picked a path by function id.
