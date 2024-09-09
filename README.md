# go-todo
A RESTfull Go API for a ToDo app


## Endpoints

There is an `http://localhost:8080/hello` endpoint which can be used to be sure the API is up and running well.

Besides this, there are two set of endpoints


### User Endpoints

These endpoints can be used for signing up and login.

#### POST `/user/register`
Calling `POST http://localhost:8080/user/register` will register a new user. The required body parameters are `email`, `name` and `password`.

**Request Body**:
```
{
  "email": "john@exmaple.org",
  "name": "John Doe",
  "password": "securepassword"
}
```

On success you'll get a HTTP 200 response.

#### POST `/user/login`
Calling `POST http://localhost:8080/user/login` with correct `email` and `password` will result in a JWT token which can be used to use ToDo related endpoints.

**Response**:
```
{
  "email": "john@exmaple.org",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImpvaG5AZXhtYXBsZS5vcmciLCJleHAiOjE3MjU4OTM3MzcsInVzZXJfaWQiOjF9.rMkAfU78ifi7X_63dbBBIYYsReqnoUDV5i46ISUgF-I"
}
```

So save this `token` and use it in `Authorization` header of the requests.


## ToDo endpoints

ALl these endpoints are all protected and require the `Authorization` header be set with `token` value:
![alt text](example-authorization-header.png)

This is to be sure that each user to do items are their own and other users can't access or edit them.

#### POST `/todos`

You can create a new to do item with this. This endpoint gets two parameters, `title` string and `completed` boolean.

**Request Body**:
```
{
  "title": "Create a todo app",
  "completed": false
}
```

**Response**:
```
{
  "id": 1,
  "user_id": 1,
  "title": "Create a todo app",
  "completed": false
}
```


#### GET `/todos`

This endpoint will list all the To Do items this user has

**Response**:
```
[
    {
        "id": 1,
        "user_id": 1,
        "title": "Second To-Do",
        "completed": false
    }
]
```

You can also get just one to do item if you pass the to do `id` in the URL like this: `GET /todos/1`


#### PUT `/todos/[todo id]`

You can update a single to do item (e.g. set complete to `false` or `true`).

Notice that the todo `id` must be present in the URL and also `title` and `completed` must be present in the body. 

**Request Body**:
```
{
  "title": "Create a todo app",
  "completed": true
}
```

**Response**:
```
{
  "id": 1,
  "user_id": 1,
  "title": "Create a todo app",
  "completed": true
}
```

#### DELETE `/todos/[todo id]`

This endpoint required no body and deletes the to do with the `id` given to it from database. NO GOING BACK! 

