# Graph Report - twitter  (2026-06-25)

## Corpus Check
- 27 files · ~5,147 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 121 nodes · 152 edges · 12 communities detected
- Extraction: 84% EXTRACTED · 16% INFERRED · 0% AMBIGUOUS · INFERRED: 24 edges (avg confidence: 0.8)
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
- [[_COMMUNITY_Community 9|Community 9]]
- [[_COMMUNITY_Community 10|Community 10]]
- [[_COMMUNITY_Community 11|Community 11]]

## God Nodes (most connected - your core abstractions)
1. `EmailRedirect` - 8 edges
2. `tweetRepository` - 7 edges
3. `currentUserFromToken()` - 7 edges
4. `strategyRegistry` - 6 edges
5. `init()` - 5 edges
6. `SetAuthCookies()` - 5 edges
7. `strategyFromRequest()` - 5 edges
8. `AuthCallback()` - 5 edges
9. `saveUserFromIDToken()` - 5 edges
10. `userFromIDTokenClaims()` - 5 edges

## Surprising Connections (you probably didn't know these)
- `init()` --calls--> `LoadEnv()`  [INFERRED]
  main.go → initializers/env.go
- `init()` --calls--> `ConnectDB()`  [INFERRED]
  main.go → initializers/db.go
- `init()` --calls--> `InitRepositories()`  [INFERRED]
  main.go → initializers/repositories.go
- `AuthCallback()` --calls--> `SetAuthCookies()`  [INFERRED]
  controllers/auth.go → auth/cookies.go
- `AuthRefresh()` --calls--> `SetAuthCookies()`  [INFERRED]
  controllers/auth.go → auth/cookies.go

## Communities (15 total, 3 thin omitted)

### Community 0 - "Community 0"
Cohesion: 0.23
Nodes (17): IDTokenClaims(), AuthCallback(), authCookiePath(), AuthLogin(), AuthLogout(), AuthRefresh(), claimBool(), claimString() (+9 more)

### Community 1 - "Community 1"
Cohesion: 0.12
Nodes (5): InitRepositories(), NewTweetRepository(), tweetRepository, NewUserRepository(), userRepository

### Community 2 - "Community 2"
Cohesion: 0.13
Nodes (4): EmailRedirect, Init(), newJWTValidator(), NewEmailRedirect()

### Community 3 - "Community 3"
Cohesion: 0.18
Nodes (9): init(), main(), ConnectDB(), SyncDB(), LoadEnv(), devCORSOrigins(), Stack(), ensureTweetUserForeignKey() (+1 more)

### Community 4 - "Community 4"
Cohesion: 0.33
Nodes (10): createTweetRequest, CreateTweet(), currentUserFromToken(), DeleteMyTweet(), GetMyTweets(), GetTweetByID(), LikeTweet(), parseTweetID() (+2 more)

### Community 5 - "Community 5"
Cohesion: 0.2
Nodes (7): AccessTokenFromCtx(), cookieSameSiteFiber(), SetAuthCookies(), ValidateAccessTokenAny(), AuthMe(), subjectFromAccessToken(), RequireAuth()

### Community 6 - "Community 6"
Cohesion: 0.22
Nodes (3): Config, loadAuth0Config(), strategyRegistry

### Community 7 - "Community 7"
Cohesion: 0.36
Nodes (5): doTokenForm(), ExchangeAuthorizationCode(), RefreshTokens(), tokenEndpoint(), TokenResponse

### Community 8 - "Community 8"
Cohesion: 0.29
Nodes (4): Tweet, TweetAuthor, User, TweetAuthorFromUser()

## Knowledge Gaps
- **9 isolated node(s):** `Config`, `Strategy`, `TokenResponse`, `TweetRepository`, `UserRepository` (+4 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **3 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Init()` connect `Community 2` to `Community 6`?**
  _High betweenness centrality (0.232) - this node is a cross-community bridge._
- **Why does `ValidateAccessTokenAny()` connect `Community 5` to `Community 2`, `Community 4`?**
  _High betweenness centrality (0.200) - this node is a cross-community bridge._
- **Are the 2 inferred relationships involving `currentUserFromToken()` (e.g. with `AccessTokenFromCtx()` and `ValidateAccessTokenAny()`) actually correct?**
  _`currentUserFromToken()` has 2 INFERRED edges - model-reasoned connections that need verification._
- **Are the 4 inferred relationships involving `init()` (e.g. with `LoadEnv()` and `ConnectDB()`) actually correct?**
  _`init()` has 4 INFERRED edges - model-reasoned connections that need verification._
- **What connects `Config`, `Strategy`, `TokenResponse` to the rest of the system?**
  _9 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Community 1` be split into smaller, more focused modules?**
  _Cohesion score 0.12 - nodes in this community are weakly interconnected._
- **Should `Community 2` be split into smaller, more focused modules?**
  _Cohesion score 0.13 - nodes in this community are weakly interconnected._