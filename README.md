# Order Creation Service

An **Order Creation Service** built using **Go**, **MongoDB**, **Kafka**, and **AWS SQS** to manage customer orders efficiently. The service stores order details, processes messages via Kafka and SQS, and integrates with hubs.

## 🚀 Features

- **Order Management**: Store and manage customer orders.
- **Order Items**: Each order contains multiple items with SKU and quantity.
- **Hub Integration**: Orders are linked to a hub.
- **Status Tracking**: Orders have a status field to track progress.
- **Event-Driven Architecture**: Uses Kafka for real-time order events.
- **Message Queuing**: AWS SQS ensures reliable order processing.
- **Timestamps**: Created, updated, and soft delete timestamps.

## 🛠️ Tech Stack

- **Backend**: Golang
- **Database**: MongoDB
- **Message Broker**: Apache Kafka
- **Queue System**: AWS SQS

## 📂 Database Schema

### 🛒 Order Collection (`orders`)
| Field         | Type               | Description                           |
|--------------|--------------------|---------------------------------------|
| _id          | ObjectID (PK)       | Primary Key                           |
| customer_name| String              | Name of the customer                  |
| order_no     | String              | Unique order number                   |
| order_items  | Array of OrderItems | List of items in the order            |
| status       | String              | Order status (e.g., pending, shipped) |
| created_at   | DateTime            | Timestamp when the order was created  |
| updated_at   | DateTime            | Timestamp when the order was updated  |
| deleted_at   | DateTime (Nullable) | Soft delete field                     |

### 📦 Order Item Subdocument
| Field   | Type   | Description                    |
|---------|--------|--------------------------------|
| sku_id  | String | SKU ID of the ordered item    |
| quantity| Int    | Quantity of the item ordered  |
| hub_id  | String | ID of the warehouse hub       |

## 📡 Event-Driven Communication

### 📌 Kafka Integration
- **Producer**: Publishes path of the file containing orders.
- **Consumer**: Listens for order-related events for processing.

### 📌 AWS SQS
- **Used for asynchronous order processing and decoupling services.**
- **Ensures fault tolerance and message persistence.**

## 🛠️ Setup Instructions

### 1️⃣ Clone the Repository
```sh
git clone https://github.com/Suyash-Rawat-Omniful/Order_Creation_Service.git
cd Order_Creation_Service
