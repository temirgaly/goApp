# Project Style Guidelines & Rules

## Go Coding Standards (Uber Go Style Guide)

All Go code written for the OrderPulse microservices must strictly adhere to the [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md):

1. **Pointers to Interfaces**: Do not return pointers to interfaces. Accept interfaces and return concrete types/structs.
2. **Interface Compliance**: Verify interface compliance at compile time where appropriate:
   ```go
   var _ Handler = (*httpHandler)(nil)
   ```
3. **Receiver Names & Types**: Keep receiver names short (1-2 letters), consistent across all methods of a type. Use pointer receivers when methods mutate state or if the struct contains mutexes.
4. **Zero Value Mutexes**: Keep `sync.Mutex` / `sync.RWMutex` as zero-value fields inside structs; do not pass mutexes by value or embed them publicly.
5. **Error Handling**:
   - Wrap errors with context using `fmt.Errorf("...: %w", err)`.
   - Use `errors.Is` and `errors.As` for error matching.
   - Handle errors immediately and return early (guard clauses).
6. **Struct Initialization**: Always use field names when instantiating structs.
7. **Concurrency & Goroutines**: Always ensure goroutines have clear ownership and shutdown signals (do not leak goroutines).
8. **Package Naming**: Avoid generic package names like `util`, `common`, or `helpers`. Use cohesive, descriptive domain package names.
