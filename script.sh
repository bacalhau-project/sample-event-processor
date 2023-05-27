#!/bin/bash

# Docker image configuration
IMAGE_REPOSITORY="${IMAGE_REPOSITORY:-bacalhauproject}"
IMAGE_NAME="${IMAGE_NAME:-event-processor}"
IMAGE_TAG="${IMAGE_TAG:-latest}"

# Job input configuration
EVENTS_FILE="${EVENTS_FILE:-${PWD}/var/events}"
CHECKPOINT_FILE="${CHECKPOINT_FILE:-${PWD}/var/checkpoint}"

# Job template configuration
JOB_TEMPLATE="${JOB_TEMPLATE:-job.template.yaml}"
JOB_FILE="${JOB_FILE:-job.yaml}"

# Environment variables
export BACALHAU_HOST=0.0.0.0

run_bacalhau() {
    init_files
    bacalhau serve --node-type requester,compute --allow-listed-local-paths "${EVENTS_FILE}:ro,${CHECKPOINT_FILE}:rw"
}

submit_job() {
    IS_IN_PROGRESS=$(bacalhau list --include-tag event-processor --output json | jq 'any(.[]; .State.State == "InProgress")')
    if [[ "${IS_IN_PROGRESS}" = "true" ]]; then
        echo "There is already a job in progress. Aborting."
        exit 1
    fi

    if [ ! -f "${JOB_FILE}" ]; then
        synthesize_job
    fi
    bacalhau create "${JOB_FILE}"
}

publish_events() {
    echo "Publishing events..."
    while true; do
        timestamp=$(date +"%Y-%m-%d %H:%M:%S")
        echo "$timestamp" >> "${EVENTS_FILE}"
        sleep 1
    done
}

synthesize_job() {
    sed -e "s|{{IMAGE_REPOSITORY}}|${IMAGE_REPOSITORY}|g" \
        -e "s|{{IMAGE_NAME}}|${IMAGE_NAME}|g" \
        -e "s|{{IMAGE_TAG}}|${IMAGE_TAG}|g" \
        -e "s|{{EVENTS_FILE}}|${EVENTS_FILE}|g" \
        -e "s|{{CHECKPOINT_FILE}}|${CHECKPOINT_FILE}|g" \
        "${JOB_TEMPLATE}" > "${JOB_FILE}"
}

init_files() {
    mkdir -p "$(dirname "${EVENTS_FILE}")"
    mkdir -p "$(dirname "${CHECKPOINT_FILE}")"
    touch "${EVENTS_FILE}"
    touch "${CHECKPOINT_FILE}"
}

build_image() {
    docker build --tag "${IMAGE_NAME}:latest" .
}

push_image() {
    docker buildx build --push \
        --platform linux/amd64,linux/arm64 \
        --tag "${IMAGE_REPOSITORY}/${IMAGE_NAME}:${IMAGE_TAG}" \
        --cache-from=type=registry,ref="${IMAGE_NAME}:latest" \
        .
}

setup_crontab() {
  crontab -l | { cat; echo "*/1 * * * * ${PWD}/script.sh submit-job >/var/log/bacalhau-cron.log 2>/var/log/bacalhau-cron.log"; } | crontab -
}

case $1 in
    run-bacalhau)
        run_bacalhau
        ;;
    synthesize-job)
        synthesize_job
        ;;
    submit-job)
        submit_job
        ;;
    publish-events)
        publish_events
        ;;
    build-image)
        build_image
        ;;
    push-image)
        push_image
        ;;
    setup-crontab)
        setup_crontab
        ;;
    *)
        echo "Invalid target."
        exit 1
        ;;
esac
