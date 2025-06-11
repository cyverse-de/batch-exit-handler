# batch-exit-handler

A simple tool to orchestrate the actions that need to take place in the exit handler of a batch workflow running in the Discovery Environment.

It will take the following steps:
* Send a status message that the workflow is uploading files.
* Upload the files by calling the `gocmd` tool.
* Send a final status that causes the state to change to `Completed` or `Failed` in the Discovery Environment.
* Send a status update through `argo-events` to the `app-exposer` service telling it to clean up the workflow.

If any of the steps should fail, the following steps will still be attempted.

## Configuration

`batch-exit-handler` uses the following environment variables:
* `IRODS_HOST` - The hostname to use for the file uploads.
* `IRODS_PORT` - The port to use when uploading files.
* `IRODS_USER_NAME` - The user to use when connecting to IRODS.
* `IRODS_USER_PASSWORD` - The password for the connection user.
* `IRODS_ZONE_NAME` - The zone for the IRODS connection.
* `IRODS_CLIENT_USER_NAME` - The user who ran the workflow.
* `STATUS_URL` - The URL that accepts status updates.
* `CLEANUP_URL` - The URL that accepts clean up requests.
* `USERNAME` - The user who ran the workflow. Duplicate of IRODS_CLIENT_USER_NAME.
* `UUID` - The invocation UUID assigned to the workflow by the Discovery Environment.
* `OUTPUT_FOLDER` - The directory the output files will be uploaded into.
* `WORKFLOW_STATUS` - The final status sent to the Discovery Environment for the workflow.

Additionally, `gocmd` must be in the `PATH`.
