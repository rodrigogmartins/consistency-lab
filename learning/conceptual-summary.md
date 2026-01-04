# Conceptual Summary — Consistency Trade-offs in Practice

This document summarizes the core architectural concepts demonstrated by the **Consistency Lab** experiment.

Its goal is not to reintroduce distributed systems theory, but to consolidate **practical, client-observable lessons** derived from measuring strong and eventual consistency under real latency, replication delay, and partial failures.

---

## 1. Consistency Is a Client-Observable Property

In distributed systems, consistency is defined by **what clients can observe**, not by how replication is implemented internally.

From a client’s perspective, consistency answers questions such as:

- Can I read my own write immediately?
- Will different replicas return the same value at the same time?
- What anomalies are possible after a successful write?

A system may replicate data correctly internally and still expose:

- stale reads
- read-after-write anomalies
- temporary divergence between replicas

Architecturally, **observable behavior matters more than internal mechanisms**.

## 2. Strong Consistency Shifts Cost to the Write Path

Strong consistency guarantees that once a write succeeds, all subsequent reads observe the latest state.

To achieve this, writes must:

- coordinate across replicas
- wait for remote acknowledgment
- fail or block if any required replica is unavailable

### Observable effects:

- Higher write latency (especially p95/p99)
- Reduced availability during partial failures
- Predictable and consistent reads

Strong consistency prioritizes **correctness over availability**.

This model is commonly required when violating invariants is unacceptable.

## 3. Eventual Consistency Shifts Cost to the Read Path

Eventual consistency allows writes to complete locally without waiting for replication.

Replication occurs asynchronously, which means:

- writes are fast and highly available
- reads may temporarily return stale data
- replicas converge over time

### Observable effects:

- Very low write latency
- High availability, even under failures
- Temporary inconsistency during replication windows

Eventual consistency prioritizes **availability and performance over immediate correctness**.

---

## 4. Partial Failures Reveal Real Trade-offs

Under ideal conditions, both models behave acceptably.
Under partial failures, their differences become explicit.

- **Strong consistency**
  - Protects correctness
  - Rejects writes if replication fails
  - Sacrifices availability to preserve invariants

- **Eventual consistency**
  - Preserves availability
  - Accepts writes despite replication issues
  - Exposes temporary inconsistency to clients

Neither behavior is inherently “better”.
The correct choice depends on **domain constraints and failure tolerance**.

## 5. Convergence Is a First-Class Metric

Eventual consistency does not mean “never consistent”.
It means **time-bounded inconsistency**.

Convergence answers:

- How long until all replicas observe the same value?
- What is the worst-case visibility delay?
- Are convergence times acceptable for the business?

Measuring **time to visibility** and its percentiles (p50, p95) is essential to evaluate whether eventual consistency is viable for a given use case.

## 6. Strong Consistency Hides Convergence Cost

Strong consistency does not eliminate convergence.
It **front-loads** its cost.

The client pays convergence cost during the write operation:

- higher write latency
- possible write failures

Once the write succeeds, replicas are already consistent from the client’s perspective.

This explains why:

- strong consistency has slower writes
- eventual consistency has faster writes but visible inconsistency

## 7. Tail Latency Matters More Than Averages

Average latency often hides critical behavior.

Architectural decisions are driven by:

- p95 latency
- p99 latency
- worst-case behavior under stress

Tail latency reflects:

- coordination overhead
- retries
- contention
- failure scenarios

Systems that look fast on average may still provide poor user experience under load.

## 8. Availability Is What the Client Experiences

Availability is not an internal metric.
It is defined by whether client requests succeed.

Client-observed availability considers:

- HTTP success responses
- timeouts
- retries
- perceived failures

A system can be “up” internally and still unavailable to clients.

This experiment measures **availability as experienced by the client**, which is what ultimately matters.

## 9. Consistency Is a Business Decision

Choosing a consistency model is not a purely technical decision.

It determines:

- which failures are acceptable
- how much inconsistency is tolerable
- whether correctness or availability is prioritized

Different parts of the same system often require different guarantees.

Most real-world architectures use **hybrid approaches**, combining strong and eventual consistency based on domain requirements.

## Final Thought

Consistency models are not abstractions to be debated.
They are trade-offs to be **measured, understood, and aligned with business constraints**.

This experiment exists to make those trade-offs observable.
