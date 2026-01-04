# Self-Assessment Questions

You should be able to answer all of these clearly and confidently.

Take your time — quality matters more than speed.

## 🔹 Fundamentals

1. What does consistency mean from a client’s perspective in a distributed system?

2. Why is it misleading to discuss consistency only in terms of internal replication mechanisms?

## 🔹 Strong Consistency

3. How does strong consistency affect write latency, and why?

4. Why does strong consistency tend to reduce availability under partial failures?

5. In which types of systems is strong consistency usually required, and why?

## 🔹 Eventual Consistency

6. Why does eventual consistency achieve lower write latency than strong consistency?

7. What does a stale read represent in practical terms?

8. Why is eventual consistency often described as “highly available”?

## 🔹 Trade-offs and Failure Scenarios

9. How do strong and eventual consistency behave differently when replication messages are dropped?

10. Why is rejecting writes sometimes the correct behavior in a strongly consistent system?

11. Why is accepting writes sometimes the correct behavior in an eventually consistent system?

## 🔹 Convergence

12. What does time to convergence mean in this experiment?

13. Why is convergence an important metric when evaluating eventual consistency?

14. How does strong consistency “hide” convergence cost from the client?

## 🔹 Metrics and Observability

15. Why are p95 and p99 latency more relevant than averages for architectural decisions?

16. What does “client-observed availability” mean, and why is it important?

## 🔹 Architecture-Level Thinking

17. Why is choosing a consistency model fundamentally a business decision?

18. Can a real-world system use both strong and eventual consistency at the same time? Give an example.

19. How did this experiment change (or reinforce) your understanding of consistency trade-offs?
