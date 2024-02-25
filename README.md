# tdd-based-app

Inspired by [Learn Go with tests](https://quii.gitbook.io/learn-go-with-tests/build-an-application/http-server)

You have been asked to create a web server where users can track how many games players have won.

- GET /players/{name} should return a number indicating the total number of wins
- POST /players/{name} should record a win for that name, incrementing for every subsequent POST
