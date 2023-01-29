# Flipbook

Flipbook is an event sourcing platform focusing on abstracting event storage, replication and reduction correct implementation complexity, leaving developers with a powerful, painless and minimalist API. The main features are:

- Event Store 
  - [gRPC API](#event-store-grpc-api)
  - Supported Backends
    - [AWS DynamoDB](#aws-dynamodb-event-store)
- Snapshot Store
  - gRPC API
  - Supported Backends
    - Redis
- SDKs
  - Golang
  - Python (Coming soon!)
  - NodeJS (Coming soon!)

- [x] Postgres
- [x] AWS DynamoDB
- [x] Redis
- [ ] Input Validation
- [ ] Logging
- [ ] gRPC Server Instrumentation (middlewares for observability and tls)
- [ ] EventBus
- [ ] Tests & Benchmarks
- [ ] Golang Event Sourcing SDK
- [ ] Typescript Event Sourcing SDK

## Table of Contents

- [Flipbook](#flipbook)
  - [Table of Contents](#table-of-contents)
  - [A word about complexity](#a-word-about-complexity)
  - [Event Sourcing, a real-world explanation](#event-sourcing-a-real-world-explanation)
    - [Question 1](#question-1)
      - [Answer 1: yes](#answer-1-yes)
      - [Answer 2: no](#answer-2-no)
      - [Answer 3: no, but I need the state to be "auditable"](#answer-3-no-but-i-need-the-state-to-be-auditable)
    - [Question 2](#question-2)
      - [Answer 1: yes](#answer-1-yes-1)
      - [Answer 2: no](#answer-2-no-1)
    - [Question 3](#question-3)
      - [Answer 1: yes](#answer-1-yes-2)
      - [Answer 2: no](#answer-2-no-2)
    - [Question 4](#question-4)
      - [Answer 1: yes](#answer-1-yes-3)
      - [Answer 2: no](#answer-2-no-3)
  - [Flipbook Concepts](#flipbook-concepts)
    - [Version-Neutral State](#version-neutral-state)
    - [Time-Neutral State](#time-neutral-state)
    - [Event Reducer](#event-reducer)
    - [Event Store](#event-store)
    - [Event Bus](#event-bus)
    - [Snapshot Store](#snapshot-store)
  - [Flipbook Project](#flipbook-project)
    - [Repository Structure](#repository-structure)
    - [Usage](#usage)
    - [Development](#development)
  - [Flipbook Features](#flipbook-features)
    - [Event Store - gRPC API](#event-store---grpc-api)
    - [Event Store - Supported Backends](#event-store---supported-backends)
      - [AWS DynamoDB Event Store](#aws-dynamodb-event-store)
    - [Snapshot Store - gRPC API](#snapshot-store---grpc-api)
    - [Snapshot Store - Supported Backends](#snapshot-store---supported-backends)
    - [ToDos for v1](#todos-for-v1)


## A word about complexity

We found that the complexity for event sourcing implementation can be broken down to 3 main challenges:

1. Making the best decision about using it as an architectural pattern.
2. Creating an optimal design for your states and events.
3. Implementing it correctly as a part of your software and infrastructure.

And for each statement above, we'll dive into the necessary information for anyone wishing to consider event sourcing as the right choice to solve a specific use case.

## Event Sourcing, a real-world explanation

First things first, you need to decide if event sourcing is the right architecture pattern for your use case. You can search and find a lot (really, a lot) of good definitions for *"event sourcing"*. Some more in the technical side, some more detailed, but we'll focus here on the real world use. And to make a good decision, we decided that the best criteria is to answer each question below.


### Question 1

> Do I need to be able to reproduce a state based on a sequence of changes caused by external operations?

#### Answer 1: yes

For example, a wallet ledger is good example, 'cause the actual balance for any wallet can be determined by reapplying all the events (transactions like deposits, payments, withdrawals and etc) tha 'caused some change in the final state of the balance. It's important to start familiarizing with common terms in the event sourcing world.

<table width="100%">
  <tr>
    <th>Terms</th>
    <th>Example</th>
    <th>Definition</th>
  </tr>
  <tr>
    <th>"Event"</th>
    <td></td>
    <td>An event is always a description of what happened in the past containing information about an operation that lead to some change.</td>
  </tr>
  <tr>
    <th>"State", "Aggregate"</th>
    <td></td>
    <td>"Is the combination of the original values in the object plus any modifications made to them." you can search for program state, object-oriented programming or other context-aware definitions of state in software. The core concept is always, the same.</td>
  </tr>
  <tr>
    <th>"Command", "External Operation"</th>
    <td></td>
    <td></td>
  </tr>
  <tr>
    <th>"Reduce", "Apply Event"</th>
    <td></td>
    <td></td>
  </tr>
  <tr>
    <th>"Event Sourcing"</th>
    <td></td>
    <td></td>
  </tr>
</table>

#### Answer 2: no

So you will be ok with a state persistence architecture ideal for saving the final state resulting from each operation. In this case we recommend to design the optimal state structure for your read/write operations.

- **Relational:** Postgres, MySQL, MariaDB and etc.
- **Graphs:** Neo4j, DGraph and etc.
- **Documents (NoSQL):** MongoDB, DynamoDB, Firestore and etc.
- **Search Engine:** Elasticsearch and etc.
- **Timeseries:** InfluxDB, Postgres, Prometheus, Thanos and etc.

And many, many others specialized engines for each data structure.

#### Answer 3: no, but I need the state to be "auditable"

You can follow the same approach commented on the [Answer 2]() and implement a properly [audit log](https://en.wikipedia.org/wiki/Audit_trail) for your [software/system](https://www.datadoghq.com/knowledge-center/audit-logging/).

### Question 2

> Do I know all the possible events that can mutate the state?

#### Answer 1: yes

To see if your concept of event is right, you can write down the field specification below for each event:

<table width="100%">
  <tr>
    <th>Field</th>
    <th>Type</th>
    <th>Example</th>
    <th>Description</th>
  </tr>
  <tr>
    <th>State ID</th>
    <td>Unique Identifier</td>
    <td><code>"cfd7cdf3-9f96-4932-9587-9566c16403ab"</code>, <code>"some.nickname"</code>, <code>"/unique/path/to/file"</code></td>
    <td>Unique identifier of the mutated state that generated this event.</td>
  </tr>
  <tr>
    <th>State Version</th>
    <td>Integer</td>
    <td><code>1</code>, <code>2</code>, <code>null</code></td>
    <td>A monotonic increasing sequence, greater than zero that places the event in the correct ascending order to be applied on the state. The <code>0</code> version is used to point to the last version and to store state metadata. If the state version is <code>null</code> we call it a <a href="#version-neutral-state">version-neutral state</a>.</td>
  </tr>
  <tr>
    <th>Event Timestamp</th>
    <td>Date and time</td>
    <td><code>"2022-12-04T14:46:45.402Z"</code>, <code>1670165226</code>, <code>null</code></td>
    <td>The temporal representation of a point in time when the event was generated. Yes, you can have events with this timestamp pointing to the future, and we can have events happening in the same time. For more information read about <a href="#time-neutral-state">time-neutral state</a>.</td>
  </tr>
  <tr>
    <th>Event Name</th>
    <td>String</td>
    <td><code>"Approved"</code>, <code>"Sent"</code>, <code>"Bloomed"</code></td>
    <td>An identifier for this event to be used when reducing an event sequence. It defines which logic to apply to the state using the <b>Event Payload</b> information. See <a href="#event-reducer">event reducer</a> for an extensive definition  of this concept.</td>
  </tr>
  <tr>
    <th>Event Payload</th>
    <td>Any</td>
    <td><code>{"amount":56,"from":"Winter","to":"Spring"}</code></td>
    <td>All the information about the mutation caused on the state. It only has a meaning with the <b>Event Name</b> property.</td>
  </tr>
<table>

#### Answer 2: no

Then we recommend to invest some time to design your use case through [DDD practices](https://khalilstemmler.com/articles/domain-driven-design-intro/).

### Question 3

> Does my use case implies multiple concurrent mutation operations happening over the same state?

#### Answer 1: yes

While this can be a symptom of a bad-design outcome for your data structure, sometimes it's an inevitable characteristic of your use case.

#### Answer 2: no

### Question 4

> Does my use case have a heavier workload on the write side?

#### Answer 1: yes

#### Answer 2: no

## Flipbook Concepts

### Version-Neutral State

### Time-Neutral State

### Event Reducer

### Event Store

### Event Bus

### Snapshot Store

<a name="time-bounded-event"></a>

## Flipbook Project

### Repository Structure

### Usage

### Development

## Flipbook Features

### Event Store - gRPC API

<a name="event-store-grpc-api"></a>

### Event Store - Supported Backends

<a name="aws-dynamodb-event-store"></a>

#### AWS DynamoDB Event Store

### Snapshot Store - gRPC API

### Snapshot Store - Supported Backends



### ToDos for v1

- SnapshotStore
  - [ ] Developer Documentation
  - [X] Input Validation
  - [ ] Metrics Endpoint
  - [ ] Error Handling/Logging
  - CI (GitHub Actions)
    - [ ] Test Automation
    - [ ] Container Build
  - Backend Support
    - [X] Redis (Ring Cluster)
- EventStore
  - [ ] Developer Documentation
  - [X] Input Validation
  - [ ] Metrics Endpoint
  - [ ] Error Handling/Logging
  - CI (GitHub Actions)
    - [ ] Test Automation
    - [ ] Container Build
  - Backend Support
    - [X] AWS DynamoDB
- Cloud Support
  - [ ] AWS
    - [ ] Terraform module for ECS
    - [ ] Terraform module for ElastiCache
    - [ ] Terraform module for DynamoDB
  - [ ] GCP
  - [ ] Azure