# go-opera Skills

Guidance for generating, migrating, and reviewing Go code using module
`github.com/colafanta/go-opera`, imported as `opera`.

## Categories

| Skill | Use For | Detailed References |
| --- | --- | --- |
| [Result flow](./opera-result-flow/SKILL.md) | Fallible operations, `Do`/`Try`/`Yield`, legacy error checks, error mapping and recovery | Migration patterns and Result API |
| [Optional values](./opera-option/SKILL.md) | Missing values, maps, slices, pointers, nullable fields, serialization | Presence patterns and integration rules |
| [Async coordination](./opera-async/SKILL.md) | Concurrent Result tasks, waits, joins, first success, cancellation | Async contracts and patterns |
| [Resource workflows](./opera-workflow/SKILL.md) | Resource scopes, reverse-order finalization, rollback | Cleanup equivalence and examples |

Start with the category that owns the behavior. Load another category only when
the change crosses its boundary. Each `SKILL.md` links directly to its references;
read only those relevant to the task.

## Shared Rules

- Use `Result[T]` for operational failure and `Option[T]` for expected absence.
- Prefer small, linear `Do(func() T)` blocks internally; preserve conventional
  `(T, error)` or `error` public APIs with `.Get()` or `.Err()`.
- Use `opera.Unit` and `opera.U` for success without a payload.
- `Yield` and `Must*` abort through a panic wrapper; use them inside `Do` unless
  the caller explicitly manages that wrapper. Ordinary panics still propagate.
- Preserve observable behavior before shortening code: returned values, error
  identity and text, call order, side effects, resource lifetimes, and cancellation.
- Do not force a migration when explicit Go branching is clearer or equivalent
  behavior cannot be maintained.
- Verify against the installed package revision. This checkout declares Go 1.27
  and uses generic methods; do not assume an older Go toolchain accepts its API.

## Distribution And Discovery

These are package-distributed skills under `skills/`, not automatically discovered
workspace skills. To enable discovery in a consuming project, copy the desired
category folders into `.agents/skills/` or `.github/skills/`, preserving their
`SKILL.md` and `references/` layout. Alternatively, explicitly load a category's
`SKILL.md` from this directory. No assets are needed for these documentation-only
skills.