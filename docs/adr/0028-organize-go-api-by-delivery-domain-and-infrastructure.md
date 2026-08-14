# Organize the Go API by Delivery, Domain, and Infrastructure

The Go API will use `cmd/api` plus internal packages for HTTP delivery, domain logic, services, repositories, configuration, and database setup. This keeps Fiber handlers, business rules, GORM persistence, and startup wiring separate enough to test and refactor during later course phases.
