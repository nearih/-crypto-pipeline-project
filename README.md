# crypto-pipeline-project

This project is a mvp demonstration of how to stream data to pipeline via grpc and how to use channel and go routine to create data pipeline

note: I used vendor because I adjusted code in go-ftx library

## producer-app

Stream crypto ticker data from exchange and stream it to pipeline

note: default run on port:7071

## crypto-data-pipeline

Received ticker data from producer-app and save to influx db

note: default run on port:7072

## Flow

producer-app -> crypto-data-pipeline -> influxDb -> grafana


## How to run

1. run docker compose in crypto-data-pipeline to setup  environment

    ```docker-compose up -d```

2. start producer-app by run the following command in folder producer-app
    ```make dev```

3. start crypto-data-pipeline by run the following command in crypto-data-pipeline
    ```make dev```

4. send request shown in command.txt in producer-app to start the pipeline 