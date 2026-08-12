# Game model

Working specification of how a game of stock_ticker plays out. Mechanics marked
*TBD* are still open. Last updated 2026-08-11.

## Turn structure

Players act **asynchronously** within a turn — there is no per-player turn
order. A player submits their actions (orders, etc.) whenever they like, then
signals **EOT** (end of turn).

A turn closes when **every active player has sent EOT**, or when the
**idle timeout** expires, whichever comes first. The timeout exists so a single
absent player cannot stall the table indefinitely. *TBD: timeout duration, and
what happens to a player who times out — auto-EOT for this turn, or marked idle
and excluded from the barrier until they act again?*

## Arbitration

Once the turn closes, the daemon runs the **arbitration phase**, entirely
server-side:

1. Resolve price movements.
2. Revalue every portfolio.
3. Retire any player who is bankrupt.
4. Check the end condition.

The daemon then notifies all players of the results: portfolio value
increase/decrease, retirements, and the new market state.

Arbitration is atomic from the clients' point of view — no client observes a
half-arbitrated market.

## Joining

New players may request admission at any time, but **admission is only granted
at an EOT boundary**, never mid-turn. This keeps the EOT barrier's participant
set fixed for the duration of a turn: the count of players a turn is waiting on
cannot change under it.

A joining player's seed money comes from a formula *TBD* — presumably scaled to
the current state of the game so that joining late is neither a windfall nor
hopeless.

## Leaving

Players may leave at **any** time, including mid-turn. A departure shrinks the
barrier immediately: if the leaver was the last player the turn was waiting on,
the turn closes at once.

## End condition

The game ends when **one non-bankrupt player remains**.

Special case: if only two players remain and one abandons, the game ends
**immediately** — the remaining player is proclaimed winner without waiting for
an EOT or a further arbitration pass.

## Consequences for the daemon

The EOT barrier is the sharp edge of this design. Three separate things mutate
the participant set — EOT arrival, join, leave — and each one can trip the
barrier or the end condition. That is much easier to reason about if
"record the event, then re-evaluate the barrier and the end condition" is a
single atomic step, which argues for one state-owning goroutine over mutexes
scattered across HTTP handlers.

See [REST-API.md](REST-API.md) for how this is exposed over the wire.
