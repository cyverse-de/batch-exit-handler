package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func sendStatus(statusURL, jobUUID, message, state string) error {
	body := map[string]string{
		"job_uuid": jobUUID,
		"hostname": "batch",
		"message":  message,
		"state":    state,
	}

	b, err := json.Marshal(&body)
	if err != nil {
		return err
	}

	resp, err := http.Post(statusURL, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		return fmt.Errorf("status code %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func sendCleanupMessage(cleanUpURL, jobUUID string) error {
	body := map[string]string{
		"uuid": jobUUID,
	}

	b, err := json.Marshal(&body)
	if err != nil {
		return err
	}

	resp, err := http.Post(cleanUpURL, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		return fmt.Errorf("status code %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func uploadFiles(outputFolder string) error {
	cmd := exec.Command(
		"gocmd",
		"--log_level=debug",
		"put",
		"-f",
		"--no_root",
		".",
		outputFolder,
	)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func main() {
	var err error

	_, err = exec.LookPath("gocmd")
	if err != nil {
		log.Fatal("gocmd must be in the PATH")
	}

	irodsHost := os.Getenv("IRODS_HOST")
	if irodsHost == "" {
		log.Fatal("IRODS_HOST environment variable must be set")
	}
	irodsPort := os.Getenv("IRODS_PORT")
	if irodsPort == "" {
		log.Fatal("IRODS_POST environment variable must be set")
	}
	irodsUser := os.Getenv("IRODS_USER_NAME")
	if irodsUser == "" {
		log.Fatal("IRODS_USER_NAME environment variable must be set")
	}
	irodsPassword := os.Getenv("IRODS_USER_PASSWORD")
	if irodsPassword == "" {
		log.Fatal("IRODS_USER_PASSWORD environment variable must be set")
	}
	irodsZone := os.Getenv("IRODS_ZONE_NAME")
	if irodsZone == "" {
		log.Fatal("IRODS_ZONE_NAME environment variable must be set")
	}
	irodsClientUser := os.Getenv("IRODS_CLIENT_USER_NAME")
	if irodsClientUser == "" {
		log.Fatal("IRODS_CLIENT_USER_NAME environment variable must be set")
	}
	statusURL := os.Getenv("STATUS_URL")
	if statusURL == "" {
		log.Fatal("STATUS_URL environment variable must be set")
	}
	cleanupURL := os.Getenv("CLEANUP_URL")
	if cleanupURL == "" {
		log.Fatal("CLEANUP_URL environment variable must be set")
	}
	username := os.Getenv("USERNAME")
	if username == "" {
		log.Fatal("USERNAME environment variable must be set")
	}
	jobUUID := os.Getenv("UUID")
	if jobUUID == "" {
		log.Fatal("UUID environment variable must be set")
	}
	outputFolder := os.Getenv("OUTPUT_FOLDER")
	if outputFolder == "" {
		log.Fatal("OUTPUT_FOLDER environment variable must be set")
	}
	workflowStatus := os.Getenv("WORKFLOW_STATUS")
	if workflowStatus == "" {
		log.Fatal("WORKFLOW_STATUS environment variable must be set")
	}

	// Send the uploading files status
	log.Println("sending uploading files status: running")
	if err = sendStatus(statusURL, jobUUID, "uploading files", "running"); err != nil {
		log.Printf("%s\n", err)
	}

	// upload files
	log.Printf("uploading files to %s\n", outputFolder)
	if err = uploadFiles(outputFolder); err != nil {
		log.Printf("%s\n", err)

		workflowStatus = "failed"
	}

	// send finished status
	log.Printf("sending final status: %s\n", workflowStatus)
	if err = sendStatus(statusURL, jobUUID, "sending final status", workflowStatus); err != nil {
		log.Printf("%s\n", err)
	}

	// send clean up message
	log.Printf("sending cleanup message for job %s\n", jobUUID)
	if err = sendCleanupMessage(cleanupURL, jobUUID); err != nil {
		log.Printf("%s\n", err)
	}
}
