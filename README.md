# image-store
Image store project to manage albums

## Technology & framework used
- GO
- GIN for REST api calls
- MySQL for database
- GORM for ORM
- DIG for dependency injection
- VIPER for config loading
- ZAP for logging
- GINKGO for testing

This project is using microservice architecture and has following Microservices and their roles:
1. Album Microservice
  - Create Album
  - Delete Album
3. Image Microservice
  - Get all images of Album
  - Delete all images of Album (This will be called from album microservice)
  - Upload Image in Album
  - Get Image from Album
4. User Microservice
  - Create User
  - Get User
5. Auth Microservice
  - Login
  - Register
  - JWT Generate
  - Validate JWT
6. API Gateway
  - Validate all requests foe valid jwt and add `user-id` header
  - Route requests to downstream microservices

## Flow of app
- User opens a web portal.
- Login or Register and get jwt token
- Call user details api
- Call get all albums api and load all the albums on ui
- User can delete album and all the images belong to that album will be deleted
- Open single album and get all the images of that album
- Upload new image
- Delete image

## Users are restricted and can see their own albums only.

## Concurreny
- On image upload the image gets processed for resizing into various sizes. GO routines are used for that.
- On image delete all the resized images will be deleted concurrently using GO routines.
- On album delete all the images inside that album and different sizes of image will be deleted concurrently using GO routines.
- Bulk upload images in GO routines.

## Security
- JWT mechanism is used for securing APIs.
- Password are stored in DB using hash.

## Integration Testing using GINKGO
- Test cases are in gateway under tests folder.
### Following scenarios are covered
1. GET all albums of user - check for empty album list.
2. GET all albums of user - check for albums in list
3. CREATE album success
4. CREATE album error
5. DELETE album succes
6. DELETE album error on passing wrong id
7. GET all images of album - check for empty list
8. GET all images of album - check for images in list
9. GET all images of album - error is wrong album id
10. UPLOAD image - success
11. UPLOAD image - error on wrong album id
12. GET/DOWNLOAD image - success
13. GET/DOWNLOAD image - error on wrong album id

## Module testing
- GINKGO is used to perform module testing.
- Album and Image are main microservices and tests for them are in tests folder.

## Improvements
- Images are stored locally inside container. In ideal scenario Object storage like S3 should be used.
- User and Album data should be cached using Redis.
- Images should come from CDN.
- gRPC for inter communication of microservices.

## How to run
### Locally using docker compose

`All TESTS will be executed on startp`

- Clone repo
- Install docker compose
- Run docker-compose build
- Run docoker-compose up
- Open url: `http://localhost:9999`

- On restart if db data persist run this command to remove db volume `docker volume rm image-store_mysql-data`

- If restarting the services please logout once from browser to remove the access token.

### Kubernetes
- TBD



