
## EmoteChat Backend API Documentation

Welcome to the EmoteChat API Documentation. This document provides information about the available endpoints for interacting with the EmoteChat backend.

### Authentication

#### Register User

Register a new user account.

- **Method**: POST
- **Endpoint**: `/api/auth/register`

##### Request Body

```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "yourpassword"
}
```

##### Response

- Status Code: 200 OK
- Content: Success message

---

#### Login User

Log in a user and obtain an authentication token.

- **Method**: POST
- **Endpoint**: `/api/auth/login`

##### Request Body

```json
{
  "email": "john@example.com",
  "password": "yourpassword"
}
```

##### Response

- Status Code: 200 OK
- Content: Authentication token

---

#### Logout User

Log out the currently authenticated user.

- **Method**: POST
- **Endpoint**: `/api/auth/logout`

##### Authorization Header

```
Authorization: Bearer <authentication_token>
```

##### Response

- Status Code: 200 OK
- Content: Success message

---

#### Validate User

Check if the current user's authentication token is valid.

- **Method**: GET
- **Endpoint**: `/api/auth/validate`

##### Authorization Header

```
Authorization: Bearer <authentication_token>
```

##### Response

- Status Code: 200 OK
- Content: User details

### User Profile

#### Get All Users

Retrieve a paginated list of all users.

- **Method**: GET
- **Endpoint**: `/api/profile/users`

##### Authorization Header

```
Authorization: Bearer <authentication_token>
```

##### Query Parameters

- `page` (optional): Page number (default: 1)
- `page_size` (optional): Number of users per page (default: 30)

##### Response

- Status Code: 200 OK
- Content: List of users and pagination details

### Conversations

#### Fetch Conversations

Retrieve a paginated list of conversations.

- **Method**: GET
- **Endpoint**: `/api/conversations`

##### Authorization Header

```
Authorization: Bearer <authentication_token>
```

##### Query Parameters

- `page` (optional): Page number (default: 1)
- `page_size` (optional): Number of conversations per page (default: 30)

##### Response

- Status Code: 200 OK
- Content: List of conversations and pagination details

### Password Management

#### Request Password Reset

Request a password reset for a user.

- **Method**: POST
- **Endpoint**: `/api/password/reset-request`

##### Request Body

```json
{
  "email": "john@example.com"
}
```

##### Response

- Status Code: 200 OK
- Content: Success message

---

#### Reset Password

Reset a user's password using a reset token.

- **Method**: POST
- **Endpoint**: `/api/password/reset`

##### Request Body

```json
{
  "token": "reset_token",
  "password": "newpassword"
}
```

##### Response

- Status Code: 200 OK
- Content: Success message

---

Please note that some endpoints require authentication via an Authorization header containing a valid Bearer token. Replace `<authentication_token>` with the actual authentication token obtained during login.
