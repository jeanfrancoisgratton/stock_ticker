# Daemon API design

Design notes for the stock_ticker daemon's REST API. Last updated 2026-08-11.

## Shape

The daemon exposes a **REST API** targeting Richardson maturity **level 2**:
real resources, proper HTTP verbs, meaningful status codes. Hypermedia
(level 3 / HATEOAS) is explicitly out of scope.

Explicitly rejected: the level-0 pattern of a single `POST /` endpoint carrying
a `{"command": "..."}` field dispatched by a `switch`. The old
`__oldcode/prometheusListener` code worked that way and was deleted on
2026-08-11 for exactly this reason. It discards caching, idempotency, status
codes, per-route middleware, and readable logs.

## Actions are resources

Game actions are verbs, but they are modelled as **resources that get created**,
not as RPC method names:

    POST /v1/games/{g}/rolls      not  GET /rollDice
    POST /v1/games/{g}/orders     not  POST / {"command":"buy"}

The response *is* the created object, and `GET`ting the collection back gives
game history for free.

Where an action genuinely will not nounify, the escape hatch is a custom method
(`POST /v1/games/{g}:reset`). Use sparingly.

## Sketch

    POST   /v1/games                        create              -> 201 + game id
    GET    /v1/games/{g}                    full state snapshot
    DELETE /v1/games/{g}                    end/abandon

    POST   /v1/games/{g}/players            request to join     -> 202 (see below)
    GET    /v1/games/{g}/players/{p}        holdings + cash
    DELETE /v1/games/{g}/players/{p}        leave

    GET    /v1/games/{g}/prices             all 10 stock prices
    POST   /v1/games/{g}/rolls              advance the ticker
    GET    /v1/games/{g}/rolls?since=N      history

    POST   /v1/games/{g}/orders             buy/sell            -> 201 or 409
    POST   /v1/games/{g}/players/{p}/eot    end of turn
    GET    /v1/games/{g}/events?since=N     push channel (SSE)

Note `202 Accepted` on join rather than `201`: admission is only granted at the
next EOT boundary, so the request is acknowledged now and fulfilled later. The
client learns it is in via the event stream.

Status codes carry meaning: `409` for "can't afford that" / "not your turn",
`404` unknown game, `403` wrong player token.

## Push

REST is strictly request/response and cannot push, but arbitration results must
reach every client unprompted. **SSE** (Server-Sent Events) is the chosen
mechanism: one long-lived `GET`, server writes `data: {...}` as events occur,
server-to-client only — which is exactly the required shape, since client
actions already travel as ordinary REST calls. Keeps the whole API inspectable
with `curl`. WebSockets would be more machinery than this needs.

Every state change bumps a monotonic version counter; clients pass `?since=N`
and events carry the version, which eliminates "did I miss an update" bugs.

## State ownership

The daemon is the **sole owner of truth**. Clients never compute a price,
validate their own purchase, or decide whose turn it is — they render state and
submit intents. Non-negotiable in a game with money in it.

All mutations are serialized. See the closing section of
[GAME-MODEL.md](GAME-MODEL.md) for why the EOT barrier pushes toward a single
state-owning goroutine rather than mutexes in the handlers.

## Wire format

**Never serialize display strings.** `shared.StockType` crosses the wire as its
int value or its `MessageID()` (`"stock_space"`), never `String()`. The client
owns presentation; the daemon owns identity. Localization was removed from the
project on 2026-08-11, but this separation was kept precisely so the protocol
never depends on wording.

Errors should likewise carry a stable code, not a prose message.

## Transport

Local-only play can serve this same API over a **Unix domain socket**, which
gets filesystem permissions as authentication for free. Networked play needs
TLS plus a per-player bearer token issued at join.
