# The golang package for managing Midjourney through Discord API.

The `mittelweg` package provides a Golang SDK for managing Midjourney through Discord API. It is in development.

## Installation

```bash
go get github.com/bernardoforcillo/mittelweg
```

## Usage

For use in your project, add the `mittelweg` package to your project. As follows:

```go
import "github.com/bernardoforcillo/mittelweg"
```

### How to initialize a client

```go
client := midjourney.NewMidjourneyClient("your token", "channel id")

```

### How to generate an image

The `Imagine` command is used to generate an image. It takes two arguments: a string that is the prompt for the image, and a boolean that force the execution to wait until the image is generated.

```go
generatedImage, err := client.Imagine("prompt", waitUntilGenerated)
if err != nil {
    log.Fatalf("call client.Imagine failed, err: %+v", err)
}
```

### How to upscale an image

The `Upscale` command is used to upscale an image. It takes two arguments: the index of the image to upscale, the index must be between 0 (included) and 4(excluded). and a boolean that force the execution to wait until the image is generated.

```go
upscaledImage, err := generatedImage.Upscale(index, true)
if err != nil {
	log.Fatalf("call.Upscale failed, err: %+v", err)
}
```

`SearchGeneratedMessage` and `SearchUpscaledMessage` are two utility functions
that can be used to search for messages in a channel.

## License

This project is licensed under the MIT license. See the [license.md](license.md) file for more details.
