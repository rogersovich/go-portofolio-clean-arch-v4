package utils

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
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
	endpoint := os.Getenv("MINIO_ENDPOINT")
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
	endpoint := os.Getenv("MINIO_ENDPOINT")
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
