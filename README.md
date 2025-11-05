# One Million Checkboxes - Development Plan

A real-time multiplayer checkbox game that scales from 100 to 1,000,000 checkboxes across thousands of concurrent users.

## Project Overview

This project replicates the viral [onemillioncheckboxes.com](https://onemillioncheckboxes.com/) concept while solving complex distributed systems challenges including real-time synchronization, horizontal scaling, and abuse prevention.

**Tech Stack:** Go backend, TypeScript frontend, Redis for state management, WebSocket for real-time communication.

## Development Phases

### Phase 1: Foundation (10×10 Grid - 100 checkboxes)
**Objective:** Establish reliable real-time synchronization and core architecture patterns

**System Design Focus:**
- State consistency across all connected clients
- WebSocket connection lifecycle management
- Efficient update propagation model
- Atomic data persistence patterns
- Bandwidth-optimized binary protocol

**Phase 1 Objectives:**
- [ ] Implement connection registry with proper cleanup
- [ ] Build broadcast mechanism for state updates
- [ ] Design and implement binary protocol
- [ ] Add connection health checks and heartbeat
- [ ] Implement basic rate limiting (per connection)
- [ ] Handle WebSocket connection failures gracefully
- [ ] Ensure zero state divergence between clients
- [ ] Add proper error handling and recovery
- [ ] Create Redis pipeline operations for bulk updates
- [ ] Validate architecture can scale to Phase 2

**Success Criteria:**
- 100+ concurrent connections without performance degradation
- Sub-50ms update propagation between clients
- Zero memory leaks over 24-hour test runs
- Graceful network partition and reconnection handling

---

### Phase 2: Scale (100×100 Grid - 10,000 checkboxes)
**Objective:** Solve update flooding and implement horizontal scaling

**System Design Focus:**
- Update batching to prevent system collapse
- Multi-server architecture with shared state
- Network efficiency optimization
- Load balancing across server instances
- Performance monitoring and observability

**Phase 2 Objectives:**
- [ ] Implement update batching (50ms time windows)
- [ ] Add inter-server communication (Redis pub/sub)
- [ ] Create connection load balancing mechanism
- [ ] Build metrics collection system (latency, connections, updates)
- [ ] Add circuit breakers for Redis failures
- [ ] Implement delta updates (only send changes)
- [ ] Create automated load testing framework
- [ ] Add geographic distribution testing
- [ ] Optimize Redis operations for bulk changes
- [ ] Implement backpressure mechanisms

**Success Criteria:**
- 1,000+ concurrent connections with linear scaling
- Sub-100ms end-to-end update latency at scale
- Redis operations under 10ms p99
- Graceful degradation during Redis failures
- Consistent performance under burst traffic

---

### Phase 3: Global Scale (1000×1000 Grid - 1,000,000 checkboxes)
**Objective:** Achieve internet-scale reliability and handle malicious actors

**System Design Focus:**
- Multi-region deployment strategies
- Conflict resolution for simultaneous updates
- Advanced abuse prevention mechanisms
- Operational excellence and monitoring
- Cost optimization at scale

**Phase 3 Objectives:**
- [ ] Implement multi-region Redis clustering
- [ ] Add conflict resolution for simultaneous updates
- [ ] Build advanced rate limiting (IP-based, adaptive)
- [ ] Create comprehensive monitoring dashboard
- [ ] Implement automated failover mechanisms
- [ ] Add CDN integration for static assets
- [ ] Design disaster recovery procedures
- [ ] Build coordinated attack detection system
- [ ] Implement auto-scaling triggers
- [ ] Add cost optimization monitoring

**Success Criteria:**
- 10,000+ concurrent connections globally
- 99.9% uptime under normal conditions
- Sub-200ms global update propagation
- Resilient to coordinated spam attacks
- Handles Redis cluster node failures gracefully
- Operates profitably at target scale

## Cross-Phase Evolution

### Security Maturity
- **Phase 1:** Basic input validation, connection limits
- **Phase 2:** DDoS protection, refined rate limiting
- **Phase 3:** Advanced threat detection, abuse mitigation

### Observability Growth
- **Phase 1:** Basic logging, error tracking
- **Phase 2:** Performance metrics, alerting systems
- **Phase 3:** Comprehensive monitoring, predictive analytics

### Consistency Models
- **Phase 1:** Strong consistency (single Redis)
- **Phase 2:** Strong consistency with replication lag
- **Phase 3:** Eventual consistency with conflict resolution

## Architecture Decisions to Validate

### Phase 1 Questions
- [ ] Redis as single source of truth vs. in-memory + Redis backup?
- [ ] Push vs. pull model for state synchronization?
- [ ] Optimistic client updates vs. server-authoritative?
- [ ] Acceptable latency threshold for "real-time" feel?

### Phase 2 Questions  
- [ ] Optimal update batch size and timing strategy?
- [ ] Inter-server communication pattern (pub/sub vs. direct)?
- [ ] Connection affinity vs. stateless server design?
- [ ] Backpressure mechanisms when clients lag?

### Phase 3 Questions
- [ ] CAP theorem tradeoffs for geographic distribution?
- [ ] Conflict resolution algorithm for simultaneous updates?
- [ ] Economic model for infrastructure scaling?
- [ ] Abuse prevention without impacting legitimate users?

## Quick Start

### Option 1: Docker (Recommended)

**Prerequisites:** Docker and Docker Compose installed

```bash
# Start everything with one command
docker-compose up

# Or run in background
docker-compose up -d

# View logs
docker-compose logs -f server

# Stop everything
docker-compose down
```

Server will start on `http://localhost:8080`

WebSocket endpoint: `ws://localhost:8080/ws`

### Option 2: Local Development

**Prerequisites:** Go 1.24+ installed, Redis server

```bash
# 1. Start Redis (in a separate terminal)
redis-server

# 2. Navigate to server directory
cd server

# 3. Start the Go server
go run cmd/server/main.go
```

Server will start on `http://localhost:8080`

WebSocket endpoint: `ws://localhost:8080/ws`

### Test Connection

```bash
# Send a test message
wscat -c ws://localhost:8080/ws
> {"cmd":"SET","index":42,"value":"true"}
```

---

## Development Status

**Current Phase:** Phase 1 - Foundation
- ✅ Connection registry with thread-safe tracking
- ✅ Binary protocol for bandwidth optimization
- ✅ Redis bitfield storage for checkboxes
- ⏳ Broadcasting mechanism (in progress)
- ⏳ Frontend implementation (in progress)

## Key Metrics to Track

- **Latency:** Update propagation time between clients
- **Throughput:** Updates per second the system can handle  
- **Reliability:** Uptime and error rates under load
- **Scalability:** Linear performance scaling with resources
- **Cost:** Infrastructure expense per concurrent user

---

*This plan progressively tackles fundamental distributed systems problems while building toward a production-ready system that can handle real-world scale and abuse patterns.*