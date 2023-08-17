# Your Project Name

Short description of your project.

## Table of Contents

- [Authentication](#authentication)
  - [Register User](#register-user)
  - [Login User](#login-user)
  - [Logout User](#logout-user)
- [Profile](#profile)
  - [Get All Users](#get-all-users)
- [Conversations](#conversations)
  - [Fetch Conversations](#fetch-conversations)
- [Password Reset](#password-reset)
  - [Request Password Reset](#request-password-reset)
  - [Reset Password](#reset-password)

## Authentication

### Register User

Registers a new user.

- Endpoint: `POST /api/auth/register`
- Request Body:

```json
{
  "name": "John Doe",
  "email": "johndoe@example.com",
  "password": "securepassword"
}

- Response:

```json
{
  "message": "User registered successfully"
}
