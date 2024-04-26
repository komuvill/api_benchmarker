# API BENCHMARKER

API Benchmarker is a CLI tool written in Go you can use to test how your REST APIs perform under heavy load.
I wrote this as course project for the Aalto University course Software Engineering with Large Language Models
during the spring of 2024. Goal for the project was to see how well a complete newbie in the Go programming
language can perform armed with the OpenAI GPT-4 model. I have done my best to review the code and structure it into it's logical components. However, I offer no guarantee the application works in all circumstances. The codebase is mainly GPT-4 generated, I have just coaxed the model to give me what I want and occasionally fixed some simple bugs in the code. Please use the application responsibly and against APIs you are allowed to stress test.

## Prerequisites

I have tested this application on the Ubuntu version `22.04.3`, Go versions `1.18.1` and `1.21.1`. You can use the dummy Python Flask API provided in this repo to try out the application. Python version used for the API was `3.10.12`. I make no promises on whether this application will work on any other OS or version of Go.

To get started make sure you have Go installed on your machine.

On Ubuntu run the following commands:

```bash
sudo apt update
sudo apt install golang-go
```

Verify the installation by checking the Go version:

```bash
go version
```

This command should output the installed version of Go.

To be able to run my application from the CLI you need to setup the Go workspace directories in your home folder and export some path variables. To do this run the following commands on your terminal:

```bash
mkdir -p ~/go/{bin,src,pkg}
echo "export GOPATH=\$HOME/go" >> ~/.bashrc
echo "export PATH=\$PATH:\$GOPATH/bin" >> ~/.bashrc
source ~/.bashrc
```

Finally, you need to build the application by running the `build.sh` script provided with this repo. Simply run it with

```bash
./build.sh
```

If you plan to use the dummy API for testing my application, make sure you have python installed. If you are using Ubuntu you probably have it already. To make sure, run the command

```bash
python3 --version
```

If this command doesn't output a Python version install it with

```bash
sudo apt install python3
```

Virtual environment is recommended to be used when installing the API dependencies. To setup the venv and install dependencies run these commands:

```bash
cd dummy_api
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
```

To start the web server run the command:

```bash
python app.py
```

Then start another terminal session for running the benchmarker.

## Usage

```bash
Usage:
  api_benchmarker [flags]

Flags:
  -b, --body string       The request body for POST/PUT requests. Prefix with @ to point to a file
  -c, --concurrency int   The level of concurrency for the requests. (default 1000)
  -d, --duration int      The duration of the test in seconds. (default 10)
  -h, --help              help for api_benchmarker
  -m, --method string     The HTTP method to use. (default "GET")
  -r, --requests int      The number of requests to perform. (default 10000)
  -u, --url string        The URL of the API endpoint to benchmark.
```

Concurrency controls how many requests the benchmarker makes at once. The duration flag defines when to stop making new requests. Any ongoing requests might exceed this time limit for the test. The HTTP client in the application has a hardcoded limit of 30 seconds for any one request. If a request is started before the time limit for test is reached, that request is handled until it succeeds or receives a timeout. The requests flag defines how many requests in total is performed. You will need to supply a request body for POST/PUT methods with the body flag. If a body is given, the application hardcodes the `application/json` header into the request.

When running the application, it creates the folder `output` in your cwd. Test results are output in both HTML and JSON formats.

## Examples with Dummy API

Note that the default concurrency value is quite high for the dummy API. Therefore, some failed requests are expected.

To simpy test GET against the dummy API with default values:

```bash
api_benchmarker -u http://127.0.0.1:5000/posts
```

To test the POST method with 500 concurrent requests, 50000 requests in total and with a time limit of 60 seconds:

```bash
api_benchmarker -u http://127.0.0.1:5000/posts -m POST -c 500 -d 60 -r 50000 -b '{"id": 2, "title": "lorem ipsum dolor"}'
```

You can achieve the same without typing out the entire request body by prefixing the body flag with @ and pointing it to the `dummy_api/sample_payload.json` file:

```bash
api_benchmarker -u http://127.0.0.1:5000/posts -m POST -c 500 -d 60 -r 50000 -b @dummy_api/sample_payload.json
```

To test the PUT method with default values:

```bash
api_benchmarker -u http://127.0.0.1:5000/posts/2 -m PUT -b @dummy_api/sample_payload.json
```

And to test the DELETE method with default values:

```bash
api_benchmarker -u http://127.0.0.1:5000/posts/2 -m DELETE
```
