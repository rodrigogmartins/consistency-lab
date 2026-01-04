# Self-Assessment Questions (Answers)

You should be able to answer all of these clearly and confidently.

Take your time — quality matters more than speed.

## 🔹 Fundamentals

1. What does consistency mean from a client’s perspective in a distributed system?

    **R:** From a client’s perspective, consistency means the guarantees a system provides about what data a client can observe after performing a write, especially whether subsequent reads reflect the most recent update.

2. Why is it misleading to discuss consistency only in terms of internal replication mechanisms?

    **R:** It is misleading because consistency is defined by client-observable behavior, not by internal replication details. A system may replicate data correctly internally and still expose stale reads or anomalies to clients.

## 🔹 Strong Consistency

3. How does strong consistency affect write latency, and why?

    **R:** Strong consistency increases write latency because a write must wait for coordination and acknowledgment from other replicas before completing. This adds network and synchronization overhead to the write path.

4. Why does strong consistency tend to reduce availability under partial failures?

    **R:** Strong consistency reduces availability under partial failures because if any required replica cannot acknowledge a write, the system must reject or delay the operation to preserve correctness.

5. In which types of systems is strong consistency usually required, and why?

    **R:** Strong consistency is usually required in systems where stale data can violate critical business invariants, such as financial ledgers, inventory systems, or distributed locking mechanisms.

## 🔹 Eventual Consistency

6. Why does eventual consistency achieve lower write latency than strong consistency?

    **R:** Eventual consistency achieves lower write latency because writes are confirmed locally without waiting for replication. Replication happens asynchronously in the background.

7. What does a stale read represent in practical terms?

    **R:** A stale read occurs when a client reads an older version of data because replication has not yet completed, even though a write has already been accepted.

8. Why is eventual consistency often described as “highly available”?

    **R:** Eventual consistency is described as highly available because writes can succeed as long as a single replica is reachable, without depending on coordination with other nodes.

## 🔹 Trade-offs and Failure Scenarios

9. How do strong and eventual consistency behave differently when replication messages are dropped?

    **R:** When replication messages are dropped, strongly consistent systems typically reject writes to preserve correctness, while eventually consistent systems accept the write and tolerate delayed or failed replication.

10. Why is rejecting writes sometimes the correct behavior in a strongly consistent system?

    **R:** Rejecting writes is sometimes the correct behavior because accepting them without full replication would break consistency guarantees and could violate critical invariants.

11. Why is accepting writes sometimes the correct behavior in an eventually consistent system?

    **R:** Accepting writes is sometimes the correct behavior because availability and low latency are prioritized over immediate consistency, and temporary inconsistency is an acceptable trade-off.

## 🔹 Convergence

12. What does time to convergence mean in this experiment?

    **R:** In this experiment, time to convergence is the time it takes until both replicas observe the same version of the written data.

13. Why is convergence an important metric when evaluating eventual consistency?

    **R:** Convergence is important because it defines how long inconsistency lasts, allowing the business to evaluate whether the inconsistency window is acceptable for the domain.

14. How does strong consistency “hide” convergence cost from the client?

    **R:** Strong consistency hides convergence cost by paying it upfront during the write operation. Once the write completes, all replicas are already consistent from the client’s perspective.

## 🔹 Metrics and Observability

15. Why are p95 and p99 latency more relevant than averages for architectural decisions?

    **R:** p95 and p99 latencies are more relevant because they capture tail behavior, which reflects worst-case user experience and system stress conditions that averages tend to hide.

16. What does “client-observed availability” mean, and why is it important?

    **R:** Client-observed availability measures whether requests succeed from the client’s perspective, regardless of internal failures. It is important because it directly impacts user experience and perceived reliability.

## 🔹 Architecture-Level Thinking

17. Why is choosing a consistency model fundamentally a business decision?

    **R:** Choosing a consistency model is a business decision because it defines which failures are acceptable and how correctness, availability, and latency are traded based on domain requirements.

18. Can a real-world system use both strong and eventual consistency at the same time? Give an example.

    **R:** Yes. Real-world systems often mix both models. For example, a payment system may use strong consistency for balance updates, while using eventual consistency for analytics, logs, or recommendation data.

19. How did this experiment change (or reinforce) your understanding of consistency trade-offs?
