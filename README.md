# golang-user-registration
A Golang service which handles the user SignUp - SignIn

### How it Works ?
- Just run it using after making sure that the environment variables declared in the compose file existed :
```
make local
```
- To stop local environment
```
make stop_local
```
### Health Check :
- Access the Endpoint `http://localhost:8080/healthz` to get the status
### Access to SignUp - SignIn :
`http://localhost:8080/login` and `http://localhost:8080/signup`
