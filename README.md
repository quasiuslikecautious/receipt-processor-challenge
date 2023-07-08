# Receipt Processor Challenge

<!-- ABOUT THE PROJECT -->
## About The Project

This is an implementation of the receipt-processor-challenge, matching the specifications outlined in the api.yml file. It was built using Go, and more specifically Gin for the API framework, and Docker to containerize the application.

<p align="right">(<a href="#readme-top">back to top</a>)</p>
<!-- GETTING STARTED -->
## Getting Started
To get a local copy up and running follow these steps.

### Prerequisites

To run this api, you will need to have docker installed
* Docker [installation docs]("https://docs.docker.com/get-docker/")

<br />

### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/quasiuslikecautious/receipt-processor-challenge.git
   ```
1. Setup the docker image

    ```sh
    # run this command in the project root
    docker build --tag receipt-processor-challenge .
    ```

<p align="right">(<a href="#readme-top">back to top</a>)</p>


<!-- USAGE EXAMPLES -->
## Usage

### Starting the API

To start the API, simply run 

```sh
docker run --network host -d receipt-processor-challenge
```

in the project's root, to start a docker container running the image.

By default, the api will run on port 8080, though this can be changed by changing the port number defined in the DefaultConfig funciton in /config/config.go

### Stopping the API
To stop the API, use

```sh
docker ps
```

to find the container that is running the receipt-processor-challenge image, and then run

```sh
docker stop <container_id OR container_name>
# OR
docker rm <container_id OR container_name>
```