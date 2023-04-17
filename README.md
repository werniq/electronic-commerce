# E-commerce Application

This is an e-commerce application built with Go and JavaScript, using PostgreSQL for data storage.

## Features

- User Management System: Allows administrators to manage users (create, update, delete)
- Product Management System: Allows administrators to manage products (create, update, delete)
- JWT Authentication: Provides secure authentication for users and administrators
- Pagination: Allows for easy browsing of large numbers of products
- Authorization: Requires authentication and authorization for all API calls
- Token-Based Authentication: Uses JWT to issue and verify access and refresh tokens
- PostgreSQL: Uses PostgreSQL for data storage and management

## Front end routes

- `/catalogue/{page}`: Returns a paginated list of all products
- `/create`: Allows administrators to create new products
- `/edit/{id}`: Allows administrators to edit existing products
- `/delete/{id}`: Allows administrators to delete existing products
- `/details/{id}`: Returns details about a specific product
- `/book/order/{id}`: Allows users to order a specific book
- `/register`: Allows users to register for an account
- `/login`: Allows users to login to their account
- `/logout`: Allows users to logout of their account
- `/my-profile`: Returns information about the current user's profile
- `/charge-once`: Charges a user for a single book order
- `/reset-password`: Allows users to reset their password
- `/forgot-password`: Allows users to reset their password if they forget it



## API Documentation

The following API endpoints are available for use:

- `/api/signup`: Allows users to sign up and receive an access token and refresh token
- `/api/signin`: Allows users to sign in and receive an access token and refresh token
- `/api/is-authenticated`: Allows users to check if they are authenticated
- `/api/catalogue/page/:page`: Returns a list of products with pagination support
- `/api/forgot-password`: Allows users to reset their password
- `/api/create`: Creates a new user

### User Management System

The user management system provides the following API endpoints:

- `GET /api/users`: Returns a list of all users with pagination support
- `GET /api/users/{id}`: Returns details about a specific user
- `POST /api/users`: Creates a new user
- `PUT /api/users/{id}`: Updates a user
- `DELETE /api/users/{id}`: Deletes a user

### Product Management System

The product management system provides the following API endpoints:

- `GET /api/products`: Returns a list of all products with pagination support
- `GET /api/products/{id}`: Returns details about a specific product
- `POST /api/products`: Creates a new product
- `PUT /api/products/{id}`: Updates a product
- `DELETE /api/products/{id}`: Deletes a product

### JWT Authentication

The JWT authentication system provides the following API endpoints:

- `POST /api/signin`: Authenticates a user and returns an access token and refresh token
- `POST /api/is-authenticated`: Verifies an access token and returns a new access token if valid

### Authorization

All API endpoints require authentication and authorization. The following authorization roles are available:

- `user`: Allows access to all user management API endpoints
- `admin`: Allows access to all user and product management API endpoints

### Token-Based Authentication

Token-based authentication is used to issue and verify access and refresh tokens. The following endpoints are available:

- `POST /api/signin`: Authenticates a user and issues access and refresh tokens
- `POST /api/is-authenticated`: Verifies an access token and issues a new access token

### PostgreSQL

PostgreSQL is used for data storage and management. The following tables are used:

- `users`: Stores user information
- `products`: Stores product information

## Setup

To set up the application, follow these steps:

1. Clone the repository
2. Install dependencies: `go get .`
3. Set up a PostgreSQL database and update the database configuration in the `.env` file
4. Start the application: `go run cmd/api` + `go run cmd/web`

<div align="center">
<a href="https://drive.google.com/file/d/1vUSRcBcNy61eqLWRbmnsQ_CLgYvVAg7o/view?usp=share_link"> Google Drive Video Link </a>
</div>
