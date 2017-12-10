# ADAS for Anki Overdrive
This is an ADAS which drives Anki cars on the track.

## Building and running it locally
To build and start locally, run:
```
go run *.go
```

## Building and publishing the docker image
To build and publish the Docker image run:
```
./build.sh
```
Note: This assumes that you are logged into the docker registry. Currently only the public Docker Hub is used. So makesure you are logged in via `docker login`.

## Running it in Docker Compose
See...
TODO: Add ref to `edge-docker` project

## References
[Introduction of Sarama](https://medium.com/@Oskarr3/implementing-cqrs-using-kafka-and-sarama-library-in-golang-da7efa3b77fe)
[Explanation of code below](https://engineering.randrr.com/getting-started-with-kafka-using-go-5a89f8555009)
[Exmaple Code on GitHub](https://github.com/randrr/kafka-example)
