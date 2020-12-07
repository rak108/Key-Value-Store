# Key-Value-Store

A basic Key-Value store implemented in Golang.
Response is sent back in JSON format. 

CURL or Postman can be used to interact with the key-value store server<br>
For CURL:<br>
GET request to fetch all values in the key-value store : ```curl -X GET http://localhost:8080/existing```<br>
GET request : ```curl -X GET  http://localhost:8080/<key>```<br>
POST request : ```curl -X POST http://localhost:8080/<key>-<value>```<br>
PUT request : ```curl -X PUT http://localhost:8080/<key>-<value>```<br>
DELETE request : ```curl -X DELETE  http://localhost:8080/<key>```<br>
