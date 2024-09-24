# API Development Project

## Overview
This project is a RESTful API built using Go, which allows for the management of products in a MySQL database. 
It supports standard CRUD operations and is designed to be efficient and secure.

## Table of Contents
- [Features](#features)
- [Technologies Used](#technologies-used)
- [API Endpoints](#api-endpoints)

## Features
- Designed and implemented API endpoints for product management (Create, Read, Update, Delete).
- Utilized Gorilla Mux for routing and structured the application for scalability.
- Integrated error handling to provide meaningful responses and ensure robust functionality.
- Established a MySQL database and configured secure user permissions for data integrity.
- Created comprehensive unit tests in app-test.go to ensure API functionality and reliability.
- Implemented a CI/CD pipeline using GitHub Actions for automated testing and deployment

## Technologies Used
- Go (Golang)
- MySQL
- Gorilla Mux (for routing)
- JSON (for data interchange)
- Github Actions

## API Endpoints
- GET /products: Retrieve all products.
- GET /products/{id}: Retrieve a product by ID.
- POST /products: Create a new product.
- PUT /products/{id}: Update an existing product by ID.
- DELETE /products/{id}: Delete a product by ID.   
