# [PROJECT-NAME]

Backend service Template

## Test Command

Navigate to ./server/tests and run command below

```sh
go test ./...
```

## Project structure

    ├── server                 - Application source code
    │   ├── api                - api related code (controleers, handlers, dto, middleware...)
    │   ├── cmd                - Go appliction entry point
    │   ├── config             - Application config from env.
    │   ├── pkg                - Packages or Dependencies
    │   ├── tests              - Test files.  
    ├── scripts                - Folder dedicated for bash script or others.
    ├── Dockerfile             - Dockerfile for building the image
    └── README.md              - Current view.

## Debugging code in VS Code

Create `launch.json` under the `.vscode` folder in the root directory of the project. Add the following configurations:

This allowed you to run the project and debug the code in Visual Studio Code.

```json
{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.1.0",
    "configurations": [
      {
        "name": "Launch Go",
        "type": "go",
        "request": "launch",
        "mode": "auto",
        "env": {
          "DEBUG": "true",
          "MODE": "development",
          "SECRET_POSTGRES_DB_PASSWORD": "postgres_password",
        },
        "program": "${workspaceFolder}/server/cmd/gin"
      },
      {
        "name": "Launch Go Cli",
        "type": "go",
        "request": "launch",
        "mode": "auto",
        "env": {
          "DEBUG": "true",
          "MODE": "development",
          "SECRET_POSTGRES_DB_PASSWORD": "postgres_password",
        },
        "program": "${workspaceFolder}/server/cmd/cli"
      }
    ]
  }
```
## Environment Configuration

Staging and Production might not be useful as we could be using configmaps and secret for k8s.

## Setup local postgresql database

### Start postgresql in your local

Please change the launch.json and environment variable respectively if you have set up user and password for your db.

```shell
$ brew install postgresql

$ brew services start postgresql
```