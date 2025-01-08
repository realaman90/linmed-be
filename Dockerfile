#Start with official golang image
FROM golang:1.23

# Set GOPROXY for reliable module downloads
ENV GOPROXY=https://proxy.golang.org

#Set up working directory
WORKDIR /app

#Copy go modules and souce file
COPY go.mod go.sum ./
RUN go mod download 

# Install Air (use specific or corrected module path)
RUN go install github.com/air-verse/air@latest

#Copy the rest of the application code
COPY . .

#Build the application
RUN go build -o main .

#Expose port 8080
EXPOSE 8080

# Use Air to start the app
CMD ["air"]