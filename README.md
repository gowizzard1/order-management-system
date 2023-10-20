# Order Management System

This project is a simple order management system implemented using microservices. The system includes two services: an "Orders Service" for managing customers and products and a "Payment Service" for handling payments through M-Pesa. The backend is written in Golang, and the services run on Google Cloud using Terraform for infrastructure configuration.

## Table of Contents
1. [Overview](#overview)
2. [Getting Started](#getting-started)
3. [Orders Service](#orders-service)
4. [Payment Service](#payment-service)
5. [Terraform and Google Cloud Run Deployment](#terraform-and-google-cloud-run-deployment)
6. [Docker](#docker)
7. [Unit Testing and Logging](#unit-testing-and-logging)
8. [GitHub Workflow](#github-workflow)

## Overview

This project aims to provide a simple yet functional order management system with the following components:

- **Orders Service**: Handles CRUD operations on orders, customers, and products. Each order comprises a customer and one or more products. Customers can place multiple orders. Google Cloud Firestore is used for database storage.

- **Payment Service**: Manages order payments. It receives order information and processes payments through M-Pesa. A callback endpoint is implemented to update the order status in the Orders Service once payment processing is complete.

## Getting Started

1. **Set Up Your GitHub Repository**:
  - Create a new, public repository on GitHub for this project. You will use this repository to submit your work.
  - Initialize Git on your local machine and connect your local repository to your GitHub repository.

2. **Clone the Repository**:
   ```bash
   git clone https://github.com/gowizzard1/order-management-system.git
   cd order-management-system
