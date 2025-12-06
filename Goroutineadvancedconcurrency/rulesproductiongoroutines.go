

# 🚀 **PRODUCTION CONCURRENCY RULES — The 12 Golden Laws**



# ✅ **1. NEVER share memory without synchronization**

Use one of:

* `sync.Mutex`
* `sync.RWMutex`
* `atomic.Value` / `atomic.Bool`
* Channels (ownership transfer)

❌ Do not read/write same variable from multiple goroutines.

---

# ✅ **2. ALWAYS give goroutines an exit path (ctx.Done())**

Every goroutine must include:

```go
select {
case <-ctx.Done():
    return
}
```

Without this → **GOROUTINE LEAKS**.

---

# ✅ **3. USE buffered channels for queues**

Don't use unbuffered channels for high throughput.

Correct:

```go
jobs := make(chan job, 100)
```

Why?
✔ prevents deadlocks
✔ prevents goroutine blocking
✔ allows backpressure

---

# ✅ **4. NEVER close a channel from the receiver side**

Rule:

> Only the **sender** closes the channel.

Closing from the wrong place → panic and system crash.

---

# ✅ **5. DO NOT send to a closed channel**

Always ensure:

* channel is open, or
* send happens before close

Use:

```go
select {
case ch <- value:
case <-ctx.Done():
    return
}
```

---

# ✅ **6. ALWAYS combine channels with select**

Never block directly on `<-ch`.

Correct:

```go
select {
case v := <-ch:
case <-ctx.Done():
    return
}
```

This prevents deadlocks & leaks.

---

# ✅ **7. FOR PARALLEL TASKS → Use errgroup.WithContext**

Benefits:

* automatic cancellation
* automatic error propagation
* clean waiting

This is the gold standard in production.

---

# ✅ **8. LIMIT goroutine count using WORKER POOLS**

Never spawn unlimited goroutines in loops.

Correct:

```go
workers := 10
```

Practical benefit:

✔ stable memory
✔ predictable concurrency
✔ avoids CPU spikes

---

# ✅ **9. FOR long-running services → Add TIMEOUTS everywhere**

Use:

```go
context.WithTimeout(ctx, 1*time.Second)
```

Timeouts prevent:

❌ infinite retry
❌ stuck goroutines
❌ blocked channels

---

# ✅ **10. NEVER block on send when using unbuffered channels**

Wrong (leaks):

```go
ch <- value
```

Correct:

```go
select {
case ch <- value:
case <-ctx.Done():
    return
}
```

---

# ✅ **11. ALWAYS drain channels before exit**

Typical pattern:

```go
close(work)
wg.Wait()
close(results)
```

If you don’t drain → goroutines will get stuck forever.

---

# ✅ **12. DESIGN for graceful shutdown**

Steps:

1. cancel root context
2. stop accepting new work
3. workers finish remaining tasks
4. wait for completion
5. close all resources

This is required for:

* Kubernetes
* systemd services
* gRPC servers
* Payment and trading applications

---

# 🔥 **BONUS: PRODUCTION SIGNAL FLOW RULE**

Every goroutine should behave like this:

```
INPUT → PROCESS → OUTPUT → EXIT
```

Not:

```
INPUT → PROCESS → WAIT FOREVER
```

---

# ⭐ **FINAL ONE-LINE SUMMARY**

> **In production Go systems: use channels for communication, mutexes for shared memory, context for cancellation, worker pools for concurrency limits, and select for leak prevention.**

---

# 🚀 If you want, I can give:

### ✔ A printable **Production Concurrency Checklist**

### ✔ A **diagram** of the entire rule system

### ✔ A **10/10 enterprise template** used by Razorpay / Netflix

### ✔ A **mock interview Q&A** for senior backend roles

Just say **"give checklist"** or **"give diagram"**.
