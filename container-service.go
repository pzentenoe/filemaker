package filemaker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

const (
	containerUploadPath = "fmi/data/%s/databases/%s/layouts/%s/records/%s/containers/%s/%s"
)

// ContainerService provides methods for uploading files to container fields in FileMaker.
// Container fields are FileMaker fields that can store files such as images, PDFs, and other documents.
type ContainerService interface {
	// UploadFile uploads a file from the local filesystem to a container field.
	// The file is read from the filePath and uploaded to the specified container field.
	UploadFile(ctx context.Context, database, layout, recordID, fieldName, filePath, token string) (*ResponseData, error)

	// UploadFileWithRepetition uploads a file to a specific repetition of a container field.
	// Repetitions are used in FileMaker for repeating fields (1-based index).
	UploadFileWithRepetition(ctx context.Context, database, layout, recordID, fieldName, filePath, token string, repetition int) (*ResponseData, error)

	// UploadData uploads binary data to a container field.
	// This is useful when you have file content in memory rather than on disk.
	UploadData(ctx context.Context, database, layout, recordID, fieldName, filename string, data []byte, token string) (*ResponseData, error)

	// UploadDataWithRepetition uploads binary data to a specific repetition of a container field.
	UploadDataWithRepetition(ctx context.Context, database, layout, recordID, fieldName, filename string, data []byte, token string, repetition int) (*ResponseData, error)

	// Download downloads the content of a container field from the provided URL.
	// Returns a ReadCloser that the caller must close.
	Download(ctx context.Context, url string, token string) (io.ReadCloser, error)
}

type containerService struct {
	client *Client
}

// NewContainerService creates a new ContainerService instance.
func NewContainerService(client *Client) ContainerService {
	return &containerService{
		client: client,
	}
}

// Download downloads the content of a container field from the provided URL.
// The URL is typically obtained from a container field in a record.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - url: The URL of the container data
//   - token: Optional session token (required for secure container data)
//
// Returns an io.ReadCloser containing the file data. The caller is responsible for closing it.
func (c *containerService) Download(ctx context.Context, url string, token string) (io.ReadCloser, error) {
	ctx = ensureContext(ctx)

	if err := validateURL(url); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add Authorization header if token is provided
	// For FileMaker Data API, secure container data often requires the session token
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	version := c.client.getVersion()
	req.Header.Set("User-Agent", fmt.Sprintf("filemaker/%s", version))

	// Execute the request
	// We don't use retryConfig.executeWithRetry here because we want to return the body stream directly
	// and retry logic usually involves reading/closing the body to check for errors, which conflicts with streaming.
	resp, err := c.client.httpClient.Do(req)
	if err != nil {
		return nil, &NetworkError{
			Message: "failed to download file",
			Err:     err,
		}
	}

	if resp.StatusCode >= 400 {
		defer func() { _ = resp.Body.Close() }()
		return nil, &FileMakerError{
			HTTPStatus: resp.StatusCode,
			Message:    http.StatusText(resp.StatusCode),
		}
	}

	return resp.Body, nil
}

// UploadFile uploads a file from the local filesystem to a container field.
// This is the most common way to upload files to FileMaker container fields.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - database: Name of the FileMaker database
//   - layout: Name of the layout containing the container field
//   - recordID: ID of the record to upload the file to
//   - fieldName: Name of the container field (e.g., "Photo" or "Table::Photo")
//   - filePath: Absolute path to the file on the local filesystem
//   - token: Valid session token from Connect or ConnectWithDatasource
//
// Example:
//
//	response, err := service.UploadFile(ctx, "MyDatabase", "Customers", "123", "Photo", "/path/to/image.jpg", token)
//
// Returns the response from the FileMaker Data API.
func (c *containerService) UploadFile(ctx context.Context, database, layout, recordID, fieldName, filePath, token string) (*ResponseData, error) {
	return c.UploadFileWithRepetition(ctx, database, layout, recordID, fieldName, filePath, token, 1)
}

// UploadFileWithRepetition uploads a file to a specific repetition of a container field.
// Repetitions are used in FileMaker for repeating fields. The repetition parameter is 1-based.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - database: Name of the FileMaker database
//   - layout: Name of the layout containing the container field
//   - recordID: ID of the record to upload the file to
//   - fieldName: Name of the container field
//   - filePath: Absolute path to the file on the local filesystem
//   - token: Valid session token
//   - repetition: Repetition number (1-based index, typically 1 for non-repeating fields)
//
// Returns the response from the FileMaker Data API.
func (c *containerService) UploadFileWithRepetition(ctx context.Context, database, layout, recordID, fieldName, filePath, token string, repetition int) (*ResponseData, error) {
	ctx = ensureContext(ctx)

	if err := c.validateCommonParams(database, layout, recordID, fieldName, token); err != nil {
		return nil, err
	}

	if err := validateFilePath(filePath); err != nil {
		return nil, err
	}

	if err := validateRepetition(repetition); err != nil {
		return nil, err
	}

	// Read the file
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, &ValidationError{
			Field:   "filePath",
			Message: fmt.Sprintf("failed to read file: %v", err),
		}
	}

	filename := filepath.Base(filePath)
	return c.uploadDataInternal(ctx, database, layout, recordID, fieldName, filename, fileData, token, repetition)
}

// UploadData uploads binary data to a container field.
// This is useful when you have file content in memory or from another source.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - database: Name of the FileMaker database
//   - layout: Name of the layout containing the container field
//   - recordID: ID of the record to upload the data to
//   - fieldName: Name of the container field
//   - filename: Name to give the file in FileMaker (e.g., "photo.jpg")
//   - data: Binary content of the file
//   - token: Valid session token
//
// Example:
//
//	data := []byte{...} // file content from HTTP upload, etc.
//	response, err := service.UploadData(ctx, "MyDatabase", "Customers", "123", "Photo", "photo.jpg", data, token)
//
// Returns the response from the FileMaker Data API.
func (c *containerService) UploadData(ctx context.Context, database, layout, recordID, fieldName, filename string, data []byte, token string) (*ResponseData, error) {
	return c.UploadDataWithRepetition(ctx, database, layout, recordID, fieldName, filename, data, token, 1)
}

// UploadDataWithRepetition uploads binary data to a specific repetition of a container field.
//
// Parameters:
//   - ctx: Context for cancellation and timeout
//   - database: Name of the FileMaker database
//   - layout: Name of the layout containing the container field
//   - recordID: ID of the record to upload the data to
//   - fieldName: Name of the container field
//   - filename: Name to give the file in FileMaker
//   - data: Binary content of the file
//   - token: Valid session token
//   - repetition: Repetition number (1-based index)
//
// Returns the response from the FileMaker Data API.
func (c *containerService) UploadDataWithRepetition(ctx context.Context, database, layout, recordID, fieldName, filename string, data []byte, token string, repetition int) (*ResponseData, error) {
	ctx = ensureContext(ctx)

	if err := c.validateCommonParams(database, layout, recordID, fieldName, token); err != nil {
		return nil, err
	}

	if err := validateFilename(filename); err != nil {
		return nil, err
	}

	if err := validateFileData(data); err != nil {
		return nil, err
	}

	if err := validateRepetition(repetition); err != nil {
		return nil, err
	}

	return c.uploadDataInternal(ctx, database, layout, recordID, fieldName, filename, data, token, repetition)
}

// validateCommonParams validates common parameters used across container methods.
func (c *containerService) validateCommonParams(database, layout, recordID, fieldName, token string) error {
	if err := validateDatabase(database); err != nil {
		return err
	}

	if err := validateLayout(layout); err != nil {
		return err
	}

	if err := validateRecordID(recordID); err != nil {
		return err
	}

	if err := validateFieldName(fieldName); err != nil {
		return err
	}

	if err := validateToken(token); err != nil {
		return err
	}

	return nil
}

// uploadDataInternal handles the actual multipart upload to FileMaker.
func (c *containerService) uploadDataInternal(ctx context.Context, database, layout, recordID, fieldName, filename string, data []byte, token string, repetition int) (*ResponseData, error) {
	version := c.client.getVersion()
	url := c.client.getURL()

	// Build the upload path
	path := fmt.Sprintf(containerUploadPath, version, database, layout, recordID, fieldName, fmt.Sprintf("%d", repetition))
	completeURL := fmt.Sprintf("%s/%s", url, path)

	// Create multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the file to the form
	part, err := writer.CreateFormFile("upload", filename)
	if err != nil {
		return nil, &ValidationError{
			Field:   "upload",
			Message: fmt.Sprintf("failed to create form file: %v", err),
		}
	}

	_, err = part.Write(data)
	if err != nil {
		return nil, &ValidationError{
			Field:   "upload",
			Message: fmt.Sprintf("failed to write file data: %v", err),
		}
	}

	// Close the multipart writer to finalize the form
	err = writer.Close()
	if err != nil {
		return nil, &ValidationError{
			Field:   "upload",
			Message: fmt.Sprintf("failed to close multipart writer: %v", err),
		}
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, completeURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", fmt.Sprintf("filemaker/%s", version))
	req.Header.Set("Accept", "application/json")

	// Execute the request with retry logic
	var responseData ResponseData
	err = c.client.retryConfig.executeWithRetry(ctx, func() error {
		resp, err := c.client.httpClient.Do(req)
		if err != nil {
			return &NetworkError{
				Message: "failed to upload file",
				Err:     err,
			}
		}
		defer func() { _ = resp.Body.Close() }()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return &NetworkError{
				Message: "failed to read response body",
				Err:     err,
			}
		}

		// Parse response
		if len(respBody) > 0 {
			err = json.Unmarshal(respBody, &responseData)
			if err != nil {
				return fmt.Errorf("failed to unmarshal response: %w", err)
			}

			// Check for FileMaker errors
			if len(responseData.Messages) > 0 {
				if fmErr := ParseFileMakerError(&responseData, resp.StatusCode); fmErr != nil {
					return fmErr
				}
			}
		}

		// Check HTTP status
		if resp.StatusCode >= 400 {
			return &FileMakerError{
				HTTPStatus: resp.StatusCode,
				Message:    http.StatusText(resp.StatusCode),
			}
		}

		return nil
	})

	if err != nil {
		return &responseData, err
	}

	return &responseData, nil
}

// ContainerFileInfo provides information about a file to be uploaded.
type ContainerFileInfo struct {
	// FieldName is the name of the container field
	FieldName string
	// FilePath is the path to the file on disk (mutually exclusive with Data)
	FilePath string
	// Data is the binary content of the file (mutually exclusive with FilePath)
	Data []byte
	// Filename is the name to give the file in FileMaker (required when using Data)
	Filename string
	// Repetition is the repetition number for repeating fields (1-based, default 1)
	Repetition int
}

// NewContainerFileInfo creates a new ContainerFileInfo for a file on disk.
func NewContainerFileInfo(fieldName, filePath string) *ContainerFileInfo {
	return &ContainerFileInfo{
		FieldName:  fieldName,
		FilePath:   filePath,
		Repetition: 1,
	}
}

// NewContainerDataInfo creates a new ContainerFileInfo for in-memory data.
func NewContainerDataInfo(fieldName, filename string, data []byte) *ContainerFileInfo {
	return &ContainerFileInfo{
		FieldName:  fieldName,
		Filename:   filename,
		Data:       data,
		Repetition: 1,
	}
}

// WithRepetition sets the repetition number for the container field.
func (c *ContainerFileInfo) WithRepetition(repetition int) *ContainerFileInfo {
	c.Repetition = repetition
	return c
}
