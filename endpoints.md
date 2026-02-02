# User endpoints

### ERRORS FOR ALL SYSTEM 

- 401 - Unauthorized

    response 
    ```
    {
        message: "user is unauthorized"
    }
    ```

- 403 - Forbidden

    response 
    ```
    {
        message: "forbidden request"
    }
    ```

- 404 - Not found 

    response
    ```
    {
        message: "not found"
    }
    ```

- 500 - Internal server error

    response 
    ```
    {
        message: "internal server error"
    }
### POST /users/auth/sign-up 

- request 
    ```
    {
        name: string,
        email: string,
        password: string
    }
    ```
    response - 200
    ```
    {
        access_token: string
    }
    ```


- error
    ```
    email is already taken - 409
    {
        "message": "This email is already taken"
    }
    ```

### POST /users/auth/sign-in

- request
    ```
    {
        email: string,
        password: string
    }
    ```

    response - 200
    ```
    {
        access_token: string, 
        refresh_token: string
    }
    ```

### GET /user

- response
    ```
    {
        "id": 1,
        "email": "my@email.com"
    }
    ```

### POST /user/auth/refresh

- request
    ```
    {
        refresh: string
    }
    ```

    response
    ```
    {
        access_token: string
    }
    ```


# Task endpoints

### POST /tasks - создание таски

- request
    ```
    {
        title: string,
        text: string,
        status: ("IN_PROGRESS, COMPLETE, CREATE)
    }
    ```
    response
    ```
    {
        message: "task is created"
    }
    ```


### PUT /tasks/{task_id} - обновление таски

- request 
    ```
    {
        title: string, 
        text: string,
        status: string 
    }
    ```
    response
    ```
    {
        message: "task is update"
    }
    ```

### DELETE /tasks/{task_id} - удаление таски

- response
    ```
    {
        message: "task is deleted"
    }
    ```

### GET /tasks?limit=int&offset=int - получение тасок 

- response
    ```
    {
        [
            {
                title: string, 
                text: string
            },
            {
                title: string, 
                text: string
            },
            ....
        ]
    }
    ```

### GET /tasks/{task_id}

- response
    ```
    {
        title: name, 
        text: text
    }
    ```