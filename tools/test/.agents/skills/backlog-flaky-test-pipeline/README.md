# backlog-flaky-test-pipeline

This **Agent skill** automates an end-to-end flaky-test workflow: discover work (JIRA backlog or explicit tickets, or local test names), gather Trunk/diagnose evidence, run a structured investigation with classification and multi-role debate, optionally apply and verify fixes, then commit and open a PR with JIRA follow-up. **Local mode** skips JIRA/Trunk/PR and stops after verification with uncommitted diffs.

**Model tiers** appear in each phase’s front matter as `model_tier`.
They align with the skill’s `<model-assignment>` shorthand:
- **lightweight** (mechanical branching, API calls — *haiku*, *gemini flash*, *etc*),
- **standard** (code understanding and judgment — *sonnet*, *gemini pro*, *etc*),
- **reasoning** (adversarial debate roles — *opus*, *gemini ultra*, *etc*).

---

## High-level flow

```mermaid
flowchart TB
  subgraph P0["Phase 0 · lightweight"]
    A[Verify Trunk, Atlassian, golangci-lint, diagnose, LSP/CRG]
  end

  subgraph P1["Phase 1 · lightweight"]
    B[Parse mode: project / direct-ticket / local]
  end

  subgraph P2["Phase 2 · lightweight"]
    C{Mode?}
    C -->|project| D[2a: fetch N tickets]
    C -->|direct| E[2b: validate tickets]
    C -->|local| F[2d: build local records]
    D --> G[2c: prior gate + claim]
    E --> G
    F --> H[Skip JIRA claim]
  end

  subgraph P3["Phase 3 · standard parent · mixed substeps"]
    I[Evidence + analysis + classifier + debate]
    J[Checkpoint · lightweight]
  end

  subgraph P4["Phase 4 · standard"]
    K[Apply fix + lint + diagnose 10x]
  end

  subgraph P5["Phase 5 · lightweight"]
    L[Branch, commit, push, gh pr create]
  end

  subgraph P6["Phase 6 · lightweight"]
    M[JIRA comments + transitions]
  end

  P0 --> P1 --> P2
  G --> I
  H --> I
  I --> J
  J -->|PROCEED JIRA| K
  J -->|PROCEED local| K
  J -->|no PROCEED JIRA| M
  J -->|no PROCEED local| N[Session summary]
  K -->|JIRA| L --> M
  K -->|local| N
```

---

## Phase 0 — Prerequisites

`model_tier: **lightweight**`

Mermaid labels below avoid nested backticks and leading `--` inside `"{...}"` diamonds (some renderers treat that as a comment or break lexing).

```mermaid
flowchart TD
  Start([Start]) --> Local{Invocation uses local benchmark mode?}
  Local -->|yes| SkipTrunk[Skip Trunk + Atlassian checks]
  Local -->|no| Trunk[Trunk MCP probe]
  Trunk --> Atlassian[Atlassian user info + cache accountId]
  SkipTrunk --> GL[golangci-lint on PATH]
  Atlassian --> GL
  GL --> Diagnose[Run diagnose help from tools or test]
  Diagnose --> Nav[ToolSearch LSP then LSP hover or CRG fallback]
  Nav -->|neither works| Stop([Hard stop: need navigation])
  Nav -->|ok| Out["Write phase0: nav_tool, lsp_available, golangci_lint_available"]
  Out --> P1[Phase 1]
```

---

## Phase 1 — Invocation & routing

`model_tier: **lightweight**`

```mermaid
flowchart TD
  A([Parse args order]) --> L{Local benchmark mode?}
  L -->|yes| LM[Local: test_specs + optional log path flag]
  L -->|no| D{"Ticket KEY or KEY@URL?"}
  D -->|yes| DT[Direct-ticket list + ci_run_url]
  D -->|no| P{"KEY + N?"}
  P -->|yes| PM[Project mode]
  P -->|no| Ask[AskUserQuestion: project vs direct]
  Ask --> PM
  Ask --> DT
  LM --> Inv["Write invocation: mode, args, auto_mode"]
  PM --> V{N greater than 5?}
  V --> Inv
  DT --> Inv
  Inv --> R{Route}
  R -->|local| P2D[phase2d]
  R -->|project| P2A[phase2a]
  R -->|direct| P2B[phase2b]
```

---

## Phase 2a — Project mode (fetch)

`model_tier: **lightweight**` · subagent `fetch-filter`: **lightweight**

```mermaid
flowchart TD
  A[Read investigation-comment rules] --> C["Subagent: fetch-flaky-tickets loop"]
  C --> D["JSON: slim_records + skipped counts"]
  D --> E[Merge into ticket_records]
  E --> F{len records less than N?}
  F -->|yes| G[Inform user: skipped reasons]
  F -->|no| H[phase2c]
  G --> H
```

---

## Phase 2b — Direct-ticket mode (validate)

`model_tier: **lightweight**` · subagent `validation`: **lightweight** (one per ticket, parallel)

```mermaid
flowchart TD
  A[Read investigation-comment rules] --> C["Parallel subagents per ticket: validate-flaky-ticket"]
  C --> D{Result status}
  D -->|error| E[Inform user; skip ticket]
  D -->|needs_assignment_check| F[User confirms ownership]
  D -->|ok| G[Add slim_record]
  F --> G
  G --> H[phase2c]
  E --> H
```

---

## Phase 2c — Prior-investigation gate & claim

`model_tier: **lightweight**`

If **claim-ticket** yields **no** claimed issues (every candidate skipped at the prior gate or not approved), the workflow stops—there is nothing left to investigate.

```mermaid
flowchart TD
  A[Classify each ticket: prior attempts] --> B{Any prior outcome not FIXED?}
  B -->|no| D[Auto-approved]
  B -->|yes| C[Per-ticket user prompt: continue vs skip]
  C -->|skip| E[Do not claim]
  C -->|continue| D
  D --> F[claim-ticket for approved]
  F --> G{At least one issue claimed?}
  G -->|no| Stop([Stop: nothing to investigate])
  G -->|yes| P3[phase3 investigation]
```

---

## Phase 2d — Local mode (slim records)

`model_tier: **lightweight**`

```mermaid
flowchart TD
  A[Parse test specs] --> B[Locate test: LSP then CRG then grep]
  B --> C{Found?}
  C -->|no| Skip[Skip spec; warn]
  C -->|yes| D[Optional: read --log into provided_log_text]
  D --> E["Build slim_record: local_id, package, test_name"]
  Skip --> E
  E --> F{Any records?}
  F -->|no| Stop([Stop])
  F -->|yes| P3[phase3 investigation]
```

---

## Phase 3 — Investigation (evidence → fix proposal)

**Parent** `model_tier: **standard**` · **3a** substeps **lightweight** · **3b** code-analyzer **standard** · **3b-ii** classifier **standard** (parent) · **3c** **standard** · **3d** Proposer **standard**, Challenger **reasoning**, Arbiter **reasoning** · **Checkpoint** **lightweight**

```mermaid
flowchart TD
  subgraph parallel["Per-ticket subagents one message"]
    A3a[3a Evidence]
    A3a -->|local| L3a["3a-i: user_log or diagnose probe"]
    A3a -->|jira| J3a["3a-ii: Trunk + optional CI + diagnose fallback"]
    L3a --> T3["3a-iii: SKIP_TOP_LEVEL if t.Run-only"]
    J3a --> T3
    T3 -->|SKIP_TOP_LEVEL| X1([Handled in subagent])
    T3 -->|continue| B3["3b code_analyzer · standard"]
    B3 --> C3{"code_mismatch?"}
    C3 -->|yes| M3[Skip 3b-ii → 3c MISMATCH path]
    C3 -->|no| Cl["3b-ii classifier · standard + gates"]
    Cl --> D3{After gates: classify as TEST?}
    D3 -->|no| Q[RETURNED_TO_QUEUE / SKIPPED per mode]
    D3 -->|yes| E3["3c pass-through then 3d debate"]
    E3 --> F3["3d debate: proposer standard · challenger reasoning · arbiter reasoning"]
    F3 --> G3[Integrity gate: subagent_calls_made]
    G3 --> CP["Checkpoint table · lightweight"]
  end
  CP --> H{Any PROCEED verdict?}
  H -->|yes| P4[phase4]
  H -->|no; JIRA modes| P6[phase6 JIRA wrap-up]
  H -->|no; local mode| SUM[Local session summary]
```

---

## Phase 4 — Apply & verify

`model_tier: **standard**` · subagent `verification`: **standard** · **4a** ownership recheck runs in parent (JIRA only)

```mermaid
flowchart TD
  A{Local mode?}
  A -->|no| B[4a: recheck-ownership per PROCEED]
  A -->|yes| C[Skip 4a]
  B --> D["Subagent: go.mod check → apply change → golangci-lint scoped → diagnose 10x"]
  C --> D
  D --> E{10 of 10?}
  E -->|yes| F[FIXED + diff]
  E -->|no| G[PARTIAL_FIX: revert + JIRA comment if JIRA]
  F --> H{lint skipped/failed?}
  G --> H
  H -->|yes| User[User gate unless --auto]
  H -->|no| I{JIRA?}
  User --> I
  I -->|yes| P5[phase5]
  I -->|no| LocalSum[Local session summary + uncommitted diff note]
```

---

## Phase 5 — Commit & PR

`model_tier: **lightweight**`

```mermaid
flowchart TD
  A[Summary table] --> B[Branch name fix/flaky-KEY-date]
  B --> C[Stage files by name only]
  C --> D[Yubikey signing hint]
  D --> E[Commit]
  E --> F[Ownership recheck · never skipped by --auto]
  F --> G["gh pr list dedupe"]
  G --> H[Push]
  H --> I["gh pr create with JIRA + Trunk links"]
  I --> P6[phase6]
```

---

## Phase 6 — JIRA terminal update

`model_tier: **lightweight**`

```mermaid
flowchart TD
  A[For each FIXED] --> B["Transition In Review + Investigation comment OUTCOME FIXED"]
  C[For each INCONCLUSIVE] --> D["Comment INCONCLUSIVE + abandon-ticket Open"]
  B --> E[Final session summary]
  D --> E
```