#Start with official golang image
FROM golang:1.23

#Set up working directory
WORKDIR /app

#Copy go modules and souce file
COPY go.mod go.sum ./
RUN go mod download

#Copy the rest of the application code
COPY . .

#Build the application
RUN go build -o main .

#Expose port 8080
EXPOSE 8080

#Run the application
CMD ["./main"]