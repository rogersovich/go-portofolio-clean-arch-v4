package utils

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type UploadFileInput struct {
	FileHeader *multipart.FileHeader
	File       multipart.File
}

type UploadResponse struct {
	FileURL  string
	FileName string
}

func GenerateAdditionalInfo(input UploadFileInput, folder string) (fileName string, contentType string, fileSize int64) {
	rawFileName := input.FileHeader.Filename

	// Ekstrak ekstensi file (misal: .jpg, .png)
	ext := filepath.Ext(rawFileName)

	// Gunakan UUID
	uniqueID := uuid.New().String()

	// Nama file baru
	fileName = fmt.Sprintf("%s/%d_%s%s", folder, time.Now().Unix(), uniqueID, ext)
	contentType = input.FileHeader.Header.Get("Content-Type")
	fileSize = input.FileHeader.Size

	return fileName, contentType, fileSize
}

func GenerateMinioClient() (*minio.Client, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT_UPLOAD")

	accessKeyID := os.Getenv("MINIO_KEY_ID")
	secretAccessKey := os.Getenv("MINIO_KEY_SECRET")
	useSSL := os.Getenv("MINIO_SSL")

	secure, err := strconv.ParseBool(useSSL)
	if err != nil {
		return nil, err
	}

	// Init client
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: secure,
	})
	if err != nil {
		return nil, err
	}

	return minioClient, nil
}

func BuildMinioURL(endpoint, bucket, fileName string) string {
	return fmt.Sprintf("%s://%s/%s/%s", GetProtocol(), endpoint, bucket, fileName)
}

func ValidateSize(size int64) (err []FieldError) {
	// Optional: File size validation (e.g. max 2MB)
	if size > 2*1024*1024 {
		errors := GenerateFieldErrorResponse("avatar_file", "File size exceeds 2MB")
		return errors
	}

	return nil
}

func ValidateExtension(fileName string, allowedExtensions []string) (err []FieldError) {
	// Jika nil atau kosong, pakai default
	if len(allowedExtensions) == 0 {
		allowedExtensions = []string{".jpg", ".jpeg", ".png", ".webp"}
	}
	ext := strings.ToLower(filepath.Ext(fileName))
	if !slices.Contains(allowedExtensions, ext) {
		message := fmt.Sprintf("File must be %s", FormatAllowedExtensions(allowedExtensions))
		errors := GenerateFieldErrorResponse("avatar_file", message)
		return errors
	}

	return nil
}

func FormatAllowedExtensions(exts []string) string {
	n := len(exts)
	if n == 0 {
		return ""
	} else if n == 1 {
		return exts[0]
	} else if n == 2 {
		return fmt.Sprintf("%s or %s", exts[0], exts[1])
	}

	// Multiple values
	return fmt.Sprintf("%s or %s",
		strings.Join(exts[:n-1], ", "),
		exts[n-1],
	)
}

func GenerateFieldErrorResponse(field, message string) []FieldError {
	errors := []FieldError{
		{
			Field:   field,
			Message: message,
		},
	}
	return errors
}

func HandlUploadFile(file *multipart.FileHeader, folderName string) (UploadResponse, error) {
	// Step 4: Open file
	openedFile, err := file.Open()
	if err != nil {
		return UploadResponse{}, err
	}
	defer openedFile.Close()

	// Step 5: Upload to MinIO
	payload := UploadFileInput{
		FileHeader: file,
		File:       openedFile,
	}
	uploadedData, err := UploadFile(context.Background(), payload, folderName)
	if err != nil {
		return UploadResponse{}, err
	}

	return *uploadedData, nil
}

func UploadFile(ctx context.Context, input UploadFileInput, folder string) (*UploadResponse, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT_VIEW")
	bucketName := os.Getenv("MINIO_BUCKET")

	// Init client
	minioClient, err := GenerateMinioClient()
	if err != nil {
		return nil, err
	}

	// Generate additional info
	fileName, contentType, fileSize := GenerateAdditionalInfo(input, folder)

	// Upload to MinIO
	_, err = minioClient.PutObject(ctx, bucketName, fileName, input.File, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return nil, err
	}

	// Generate public URL (jika public)
	fileURL := BuildMinioURL(endpoint, bucketName, fileName)

	// Generate response
	avatar := UploadResponse{
		FileURL:  fileURL,
		FileName: fileName,
	}

	return &avatar, nil
}

func DeleteFromMinio(ctx context.Context, objectPath string) error {
	bucketName := os.Getenv("MINIO_BUCKET")

	// Init client
	minioClient, err := GenerateMinioClient()
	if err != nil {
		return err
	}

	return minioClient.RemoveObject(ctx, bucketName, objectPath, minio.RemoveObjectOptions{})
}

// Function to parse URL and extract the image key
func MinioParseURLToImageKey(urlImages []string, bucketName string) ([]string, error) {
	var imageKeys []string

	// Iterate through each URL in the slice
	for _, imageURL := range urlImages {
		// Parse the URL
		parsedURL, err := url.Parse(imageURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse URL: %v", err)
		}

		// Extract the object key by removing the base URL and bucket name
		objectKey := strings.TrimPrefix(parsedURL.Path, "/"+bucketName+"/")

		// If objectKey is empty after trimming, return an error
		if objectKey == "" {
			return nil, fmt.Errorf("invalid URL or object key not found for URL: %s", imageURL)
		}

		// Append the object key to the result slice
		imageKeys = append(imageKeys, objectKey)
	}

	return imageKeys, nil
}

func DeleteBulkImagesInBatches(bucketName string, objectKeys []string, batchSize int) error {
	// Init client
	minioClient, err := GenerateMinioClient()
	if err != nil {
		return err
	}

	if len(objectKeys) == 0 {
		return nil // Tidak ada yang perlu dihapus
	}

	// Calculate the number of batches.
	numBatches := (len(objectKeys) + batchSize - 1) / batchSize // Integer division rounding up

	// Use a WaitGroup to wait for all batches to complete.
	var wg sync.WaitGroup
	errCh := make(chan error, numBatches) // Buffered channel untuk menampung error dari setiap batch

	for i := 0; i < numBatches; i++ {
		wg.Add(1)
		start := i * batchSize
		end := start + batchSize
		if end > len(objectKeys) {
			end = len(objectKeys)
		}
		batchKeys := objectKeys[start:end]

		go func(keys []string) {
			defer wg.Done()

			// Create a channel for the current batch of objects to delete.
			deleteObjectsCh := make(chan minio.ObjectInfo, len(keys))
			for _, key := range keys {
				deleteObjectsCh <- minio.ObjectInfo{Key: key}
			}
			close(deleteObjectsCh)

			// Use RemoveObjectsWithContext to handle context cancellation/timeouts
			errorResultCh := minioClient.RemoveObjects(context.Background(), bucketName, deleteObjectsCh, minio.RemoveObjectsOptions{})

			// Check for errors during deletion.  Crucially, *drain the entire channel*.
			// RemoveObjects sends *all* errors to the channel, and if you don't
			// drain it, the goroutine can leak.
			for rErr := range errorResultCh {
				errCh <- fmt.Errorf("error menghapus objek '%s': %w", rErr.ObjectName, rErr.Err)
			}

		}(batchKeys)
	}

	wg.Wait()
	close(errCh) // Close the error channel

	// Check for any errors.
	combinedErr := handleErrors(errCh) // helper to combine errors
	if combinedErr != nil {
		return combinedErr
	}
	return nil
}

// handleErrors handles errors from the error channel.
func handleErrors(errCh <-chan error) error {
	var errorStrings []string
	for err := range errCh {
		errorStrings = append(errorStrings, err.Error())
	}
	if len(errorStrings) > 0 {
		return fmt.Errorf("multiple errors occurred:\n%s", strings.Join(errorStrings, "\n"))
	}
	return nil
}
