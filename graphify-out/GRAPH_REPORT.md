# Graph Report - twitter  (2026-06-24)

## Corpus Check
- 6 files · ~723 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 13 nodes · 10 edges · 2 communities detected
- Extraction: 70% EXTRACTED · 30% INFERRED · 0% AMBIGUOUS · INFERRED: 3 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Community 0|Community 0]]
- [[_COMMUNITY_Community 1|Community 1]]

## God Nodes (most connected - your core abstractions)
1. `init()` - 3 edges
2. `main()` - 2 edges
3. `Stack()` - 2 edges
4. `LoadEnv()` - 2 edges
5. `ConnectDB()` - 2 edges

## Surprising Connections (you probably didn't know these)
- `init()` --calls--> `LoadEnv()`  [INFERRED]
  main.go → initializers/env.go
- `init()` --calls--> `ConnectDB()`  [INFERRED]
  main.go → initializers/db.go
- `main()` --calls--> `Stack()`  [INFERRED]
  main.go → middlewares/middlewares.go

## Communities (4 total, 1 thin omitted)

### Community 0 - "Community 0"
Cohesion: 0.4
Nodes (3): init(), ConnectDB(), LoadEnv()

## Knowledge Gaps
- **1 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `init()` connect `Community 0` to `Community 1`?**
  _High betweenness centrality (0.303) - this node is a cross-community bridge._
- **Are the 2 inferred relationships involving `init()` (e.g. with `LoadEnv()` and `ConnectDB()`) actually correct?**
  _`init()` has 2 INFERRED edges - model-reasoned connections that need verification._