# Phase 2 Fixture Performance

Measured locally on 2026-05-09 with `make test` and the 10 committed real-data fixtures.

| Operation | Median | Worst | Notes |
| --- | ---: | ---: | --- |
| Parser fixture suite | <1s package runtime | <1s package runtime | `go test ./internal/mailbox` completes under one second on the fixture set. |
| Full local unit suite | ~7s | ~7s | `make test` includes Vitest and Go packages. |
| Browser smoke with upload/export path | ~15s e2e-only, ~45-70s with build | Build time dominates; the real EML upload step is sub-second after backend startup. |

Performance cliff: `.eml` uploads are capped at 25MB (`mailbox.MaxEMLBytes`). Indexed body text is capped at 2MB per part to keep imports responsive and deterministic.
