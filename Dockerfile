FROM golang:1.12.4
WORKDIR /srv
ADD . .

# Build the server
RUN make build
CMD "./marinad"
