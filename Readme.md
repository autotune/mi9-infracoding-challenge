Mi9 Infrastructure Coding Challenge
===

## About
This project was created to help recruits for Mi9 TechOps demonstrate their knowledge of solving systems problems.

## Structure
Main.go - entrypoint, binds to ports 3000, 8080 and 8081.
a martini application is listening on 3000 and serving static files from /public, and managing POST requests to /candidates

backend.go supplies the backend node logic that's listening on port 8080 and 8081. 

scenarios*.go provide a convenient framework to write testing scenarios for an external service; a load balancer in this case.

## Inspiration
[Mi9-Coding-Challenge](mi9-coding-challenge.herokuapp.com) - test your chops writing a JSON service; Join great people creating the future of media.