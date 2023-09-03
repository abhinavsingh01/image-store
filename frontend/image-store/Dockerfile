FROM node:18.12.1 as build
WORKDIR /app
COPY package*.json /app/
RUN npm install
COPY ./ /app/
ARG configuration=production
RUN npm run build-prod -- --output-path=./dist/out

FROM node:14-alpine
COPY --from=build /app/dist/out /app
WORKDIR /app
RUN npm install -g http-server
EXPOSE 9999
CMD ["http-server", "-p", "9999"]



