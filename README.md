# image-store
Image store project to manage albums

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

## Flow of app
- User opens a web portal.
- Login or Register and get jwt token
- Call user details api
- Call get all albums api and load all the albums on ui
- User can delete album and all the images belong to that album will be deleted
- Open single album and get all the images of that album
- Upload new image
- Delete image

Users are restricted and can see their own albums only.

## Concurreny
- On image upload the image gets processed for resizing into various sizes. GO routines are used for that.
- On image delete all the resized images will be deleted concurrently using GO routines.
- On album delete all the images inside that album and their sized will be deleted concurrently using GO routines.
- Bulk upload images in GO routines.

## Security
- JWT mechanism is used for securing APIs.
- Password are stored in DB using hash.

## Points
- Images are stored locally inside container. In ideal scenario Object storage like S3 should be used.
- Some data should be cached using Redis.

## How to run
### Locally using docker compoe
- TBD

### Kubernetes
- TBD



