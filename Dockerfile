FROM golang:alpine
RUN apk update && apk add --no-cache git && apk add --no-cach bash && apk add build-base

RUN mkdir /app
WORKDIR /app

# Copy the source from the current directory to the working Directory inside the container
COPY . .
COPY Environment.env .

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# Build the Go app
RUN go build -o /web_lifo

RUN rm config.json && mv config.json.docker config.json

# Expose port 8765 to the outside world
EXPOSE 8765

# Run the executable
CMD [ "/web_lifo" ]