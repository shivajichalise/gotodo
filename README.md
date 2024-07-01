# Gotodo

- A simple todo API in GO
- uses sqlite3 for data storage

## Get all todos

```http
GET /api/todos
```

### Response

```javascript
[
  {
    id: string,
    todo: string,
    is_completed: bool,
  },
];
```

## Add todo

```http
POST /api/todos
```

#### Payload:

```json
{
  "todo": "Your todo task to add"
}
```

### Response

```javascript
{
    id: string,
    todo: string,
    is_completed: bool,
}
```

## Update todo

```http
PUT /api/todos/{:todo}
```

| Parameter | Type     | Description                     |
| :-------- | :------- | :------------------------------ |
| `todo`    | `string` | **Required**. Your todo task id |

#### Payload

```json
{
  "todo": "Your modified todo task"
}
```

### Response

```javascript
{
    message: "Todo updated successfully",
}
```

## Delete todo

```http
DELETE /api/todos/{:todo}
```

| Parameter | Type     | Description                     |
| :-------- | :------- | :------------------------------ |
| `todo`    | `string` | **Required**. Your todo task id |

### Response

```javascript
{
    message: "Todo deleted successfully",
}
```

## Mark todo as complete

```http
PATCH /api/todos/{:todo}/complete
```

| Parameter | Type     | Description                     |
| :-------- | :------- | :------------------------------ |
| `todo`    | `string` | **Required**. Your todo task id |

### Response

```javascript
{
    message: "Todo marked as completed",
}
```

## Status Codes

Gotodo returns the following status codes in its API:

| Status Code | Description             |
| :---------- | :---------------------- |
| 200         | `OK`                    |
| 201         | `CREATED`               |
| 404         | `NOT FOUND`             |
| 500         | `INTERNAL SERVER ERROR` |
