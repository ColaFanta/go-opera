---
name: opera-workflow
description: "Use when generating or migrating go-opera Scoped resource handling or Workflow/Step/Finalize cleanup and rollback. Preserve acquire/release timing, reverse-order finalizers, failed-step cleanup, error joining, and panic behavior; do not mechanically replace defer or transactions."
---

# Resource Workflows

Use `Scoped` for a bounded acquire/use/release operation and `Workflow` for
explicit, reverse-order step finalization. Neither is a general transaction or
automatic compensation for arbitrary code.

## Procedure

1. Record when each resource becomes owned, when cleanup is registered, whether
   cleanup runs after acquisition failure, and how cleanup errors are handled.
2. Read [cleanup patterns](./references/cleanup-patterns.md) before migrating a
   `defer`, rollback sequence, or resource scope.
3. Choose a helper only if its registration, lifetime, and failure semantics
   match. Retain plain Go `defer` when that is the simpler equivalent.
4. Put `Step` and fallible `Finalize` calls inside `Do`. Call `Finalize` on the
   success path; do not assume a later unrelated failure finalizes a workflow.
5. Test each failed acquisition/step, success, cleanup order, finalizer failures,
   and any required panic cleanup. Verify error joining matches the old contract.