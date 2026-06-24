# Graph Report - twitter  (2026-06-24)

## Corpus Check
- 18 files · ~3,906 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 81 nodes · 104 edges · 9 communities detected
- Extraction: 85% EXTRACTED · 15% INFERRED · 0% AMBIGUOUS · INFERRED: 16 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Community 0|Community 0]]
- [[_COMMUNITY_Community 1|Community 1]]
- [[_COMMUNITY_Community 2|Community 2]]
- [[_COMMUNITY_Community 3|Community 3]]
- [[_COMMUNITY_Community 4|Community 4]]
- [[_COMMUNITY_Community 5|Community 5]]
- [[_COMMUNITY_Community 6|Community 6]]
- [[_COMMUNITY_Community 7|Community 7]]
- [[_COMMUNITY_Community 8|Community 8]]

## God Nodes (most connected - your core abstractions)
1. `EmailRedirect` - 8 edges
2. `strategyRegistry` - 6 edges
3. `SetAuthCookies()` - 5 edges
4. `strategyFromRequest()` - 5 edges
5. `AuthCallback()` - 5 edges
6. `saveUserFromIDToken()` - 5 edges
7. `userFromIDTokenClaims()` - 5 edges
8. `init()` - 4 edges
9. `doTokenForm()` - 4 edges
10. `Init()` - 4 edges

## Surprising Connections (you probably didn't know these)
- `saveUserFromIDToken()` --calls--> `IDTokenClaims()`  [INFERRED]
  controllers/auth.go → auth/jwt.go
- `AuthCallback()` --calls--> `SetAuthCookies()`  [INFERRED]
  controllers/auth.go → auth/cookies.go
- `AuthRefresh()` --calls--> `SetAuthCookies()`  [INFERRED]
  controllers/auth.go → auth/cookies.go
- `init()` --calls--> `LoadEnv()`  [INFERRED]
  main.go → initializers/env.go
- `init()` --calls--> `ConnectDB()`  [INFERRED]
  main.go → initializers/db.go

## Communities (11 total, 3 thin omitted)

### Community 0 - "Community 0"
Cohesion: 0.25
Nodes (16): AuthCallback(), authCookiePath(), AuthLogin(), AuthLogout(), AuthRefresh(), claimBool(), claimString(), clearAuthCookies() (+8 more)

### Community 1 - "Community 1"
Cohesion: 0.2
Nodes (7): AccessTokenFromCtx(), cookieSameSiteFiber(), SetAuthCookies(), ValidateAccessTokenAny(), AuthMe(), subjectFromAccessToken(), RequireAuth()

### Community 2 - "Community 2"
Cohesion: 0.22
Nodes (7): init(), main(), ConnectDB(), SyncDB(), LoadEnv(), devCORSOrigins(), Stack()

### Community 3 - "Community 3"
Cohesion: 0.22
Nodes (3): Config, loadAuth0Config(), strategyRegistry

### Community 4 - "Community 4"
Cohesion: 0.31
Nodes (5): doTokenForm(), ExchangeAuthorizationCode(), RefreshTokens(), tokenEndpoint(), TokenResponse

### Community 5 - "Community 5"
Cohesion: 0.25
Nodes (4): Init(), IDTokenClaims(), newJWTValidator(), NewEmailRedirect()

## Knowledge Gaps
- **4 isolated node(s):** `Config`, `Strategy`, `TokenResponse`, `User`
  These have ≤1 connection - possible missing edges or undocumented components.
- **3 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Init()` connect `Community 5` to `Community 3`?**
  _High betweenness centrality (0.359) - this node is a cross-community bridge._
- **Why does `EmailRedirect` connect `Community 6` to `Community 4`, `Community 5`?**
  _High betweenness centrality (0.240) - this node is a cross-community bridge._
- **Are the 3 inferred relationships involving `SetAuthCookies()` (e.g. with `RequireAuth()` and `AuthCallback()`) actually correct?**
  _`SetAuthCookies()` has 3 INFERRED edges - model-reasoned connections that need verification._
- **What connects `Config`, `Strategy`, `TokenResponse` to the rest of the system?**
  _4 weakly-connected nodes found - possible documentation gaps or missing edges._