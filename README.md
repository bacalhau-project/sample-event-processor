# Sample Event Processor

As it stands today, Bacalhau doesn't support daemon like jobs, and every job is expected to complete within a defined timeout window. We are planning to support daemon like jobs natively in the future, but for now, this sample application demonstrates how to implement a long running job that can be restarted from the last checkpointed event.

This is a sample Bacalhau application that demonstrates executing long running jobs that process events or logs from a local file, checkpoints its progress as of the last processed event, and then resumes processing from the last checkpointed event when the application is restarted.

The sample application assumes a local Bacalhau instance that is tailing and consuming events from a local file, and checkpoints its progress on another local file. When submitting another instance of the job, it aborts if there is an already running instance of the job, which can be useful to run the job as a cron job.

## Running the sample application
### Run Bacalhau Server
```
./script.sh run-bacalhau
```
This command will create `events` and `checkpoint` files under `var` directory, and run a Bacalhau server that allow-list those local files to be accessed by Bacalhau jobs. `events` will only be exposed as read-only, while `checkpoint` as read-write since the job must be able to update the checkpoint offset.

### Publish Events
```
./script.sh publish-events
```
Publish a new entry to `events` file every second until it is stopped.

### Run the sample application
```
./script.sh submit-job
```
Will check if there is an already running instance of the job, and abort if there is one. Otherwise, it will submit a new instance of the job that will process events from `events` file, and checkpoint its progress on `checkpoint` file. It has a default timeout of 5 minutes. You can try submitting the job after it is completed, and it will resume from the last checkpointed event.

## Updating the sample application
### Update the code
The sample application is written in Go and can be found under `app` directory. You can update the code by providing your own implementation of `EventProcessor` to process the individual events, and your own implementation of `checkpointer` to checkpoint the progress. 

### Build a new docker image
```
export IMAGE_REPOSITORY=<your-image-repository>
export IMAGE_NAME=<your-image-name>
exoort IMAGE_TAG=<your-image-tag>
./script push-image
```

### Re-Synthesize your job definition
```
export IMAGE_REPOSITORY=<your-image-repository>
export IMAGE_NAME=<your-image-name>
exoort IMAGE_TAG=<your-image-tag>
./script.sh synthesize-job
```

### Re-Run the job
```
./script.sh submit-job
```

## Limitations and Future Work
- This is just a sample application and not meant to be used in production. It is not tested for performance and scalability.
- Checkpoints are persisted after each event, which can be expensive.
- The sample application depends on this [PR](https://github.com/bacalhau-project/bacalhau/pull/2499) to support writable volumes.
- Future work is to support long running daemon like jobs natively in Bacalhau, such that users only need to submit the job once, and Bacalhau will take care of restarting the job if it fails or crashes.
- Future work is to support functions that scale from zero, reduce cold startup latency and support multiple invocations of the same function using the same container instance.
- Future work is to support functions that can be invoked on a schedule, and functions that can be invoked in response to events from other services.