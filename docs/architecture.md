# System Architecture

```mermaid
flowchart LR
  Client[Merchant Dashboard / API Client] --> Fiber[Fiber HTTP API]
  Fiber --> Auth[auth module]
  Fiber --> Merchant[merchant module]
  Fiber --> Payment[payment module]
  Fiber --> Wallet[wallet module]
  Fiber --> Settlement[settlement module]
  Fiber --> Notification[notification module]
  Fiber --> Audit[auditlog module]
  Auth --> Redis[(Redis)]
  Payment --> Redis
  Fiber --> Postgres[(PostgreSQL)]
  Payment --> Rabbit[RabbitMQ]
  Rabbit --> Worker[Worker Consumer]
  Worker --> Notification
  Notification --> Webhook[Merchant Webhook]
  Fiber --> Jaeger[Jaeger]
```

The application is a modular monolith. Business modules are deployed together, share one database, and communicate through Go interfaces and events. Module boundaries mirror future service boundaries, so `payment`, `wallet`, `settlement`, and `notification` can be extracted later with minimal API and data-flow changes.
