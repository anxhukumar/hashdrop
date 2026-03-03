# Scaling Strategy

## Overview

Hashdrop is currently deployed using a single EC2 instance with an embedded SQLite database.  
This design keeps infrastructure simple and operational costs low while allowing focus on application logic and security guarantees.

However, this setup introduces clear trade-offs:

- Single point of failure (compute + database co-located)
- Limited horizontal scalability
- Manual scaling of compute resources
- Database tightly coupled to the application instance

If Hashdrop were to scale beyond its current workload, there are two primary evolution paths.

---

## Path 1: Traditional Scalable Server Architecture

This approach preserves the current Go API server model while removing single points of failure and enabling horizontal scaling.

### Infrastructure Changes

- Introduce an **Application Load Balancer (ALB)** in front of the API servers
- Move EC2 instances into an **Auto Scaling Group**
- Replace SQLite with a managed relational database (e.g., PostgreSQL on AWS RDS)

### Resulting Improvements

- Traffic distributed across multiple instances
- Automatic instance scaling based on load
- Separation of compute and database layers
- High availability via multi-instance deployment
- Managed backups and failover at the database layer (via RDS)

For cost optimization at scale, encrypted object storage could also be migrated from S3 to a compatible alternative such as **Cloudflare R2**, reducing egress costs while preserving the zero-trust storage model.

This path maintains full control over the backend server while aligning with common production deployment patterns.

---

## Path 2: Serverless Cloud-Native Architecture

Alternatively, the backend could be redesigned using fully managed services.

### Infrastructure Changes

- Replace EC2 API server with **AWS Lambda**
- Use **API Gateway** for HTTP request routing and TLS termination
- Replace SQLite with **DynamoDB**
- Retain S3 and CloudFront for object storage and delivery

### Resulting Improvements

- Automatic horizontal scaling without managing servers
- No instance patching or infrastructure maintenance
- Fine-grained scaling per request
- Reduced operational overhead

This approach trades low-level infrastructure control for operational simplicity and elasticity.

---

## Design Intent

The current EC2-based deployment is intentionally simple and cost-conscious.  
It demonstrates application-layer security design, encryption guarantees, and end-to-end system coordination without prematurely optimizing infrastructure.

The scaling paths above outline how Hashdrop can evolve to meet higher availability, throughput, and production-grade requirements.