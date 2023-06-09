# DB

In order to run the tests in db_test.go, an instance of postgres
must be setup. This can be done using the command below

```bash
docker run -e POSTGRES_HOST_AUTH_METHOD=trust -p 5432:5432 postgres
```

"Tests" can be run by using VSCode's "run test" or "debug test" buttons
that become highlighted above each test.

Make sure to run Test_InitDB before any others, as this will set up
all database tables for you.
