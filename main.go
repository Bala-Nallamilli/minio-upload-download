package main

import (
	"archive/zip"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

var minioClient *minio.Client

func init() {
	var err error
	minioClient, err = minio.New("127.0.0.1:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("5SYt1IisdoBK9kjgE3YU", "9ESxEckUyeJeEqkabXYz1Q9vwHaDW2vuH1mmxdeD", ""),
		Secure: false, // Set to true if your MinIO server uses HTTPS
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Optionally, create the bucket if it doesn't already exist
	bucketName := "kalki"
	location := "local"

	err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check if the bucket already exists (which is not an error)
		exists, errBucketExists := minioClient.BucketExists(context.Background(), bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}
}

type User struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func ParseUser(c *gin.Context) {
	var user User
	if err := c.ShouldBind(&user); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("User: %#v\n", user)
}

func main() {
	router := gin.Default()

	router.POST("/upload-multiple", func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		files := form.File["evidence"] // Adjust the form field name as needed

		bucketName := "kalki" // Specify your MinIO bucket name

		for _, file := range files {
			if strings.ToLower(filepath.Ext(file.Filename)) == ".exe" {
				continue // Skip .exe files
			}

			openedFile, err := file.Open()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			defer openedFile.Close()

			// Define object name in the bucket using file's name
			// Consider sanitizing or hashing the file name to avoid conflicts
			objectName := filepath.Base(file.Filename)

			// Upload the file to the bucket
			info, err := minioClient.PutObject(context.Background(), bucketName, objectName, openedFile, file.Size, minio.PutObjectOptions{ContentType: "application/octet-stream"})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			log.Printf("Successfully uploaded %s of size %d\n", info.Key, info.Size)
		}

		c.String(http.StatusOK, "Files uploaded successfully!")
	})

	router.GET("/download", func(c *gin.Context) {
		downloadAllHandler(c)
	})

	router.Run(":8080")
}

func downloadAllHandler(c *gin.Context) {
	bucketName := "kalki" // Specify your MinIO bucket name

	// Create a ZIP writer
	c.Writer.Header().Set("Content-Type", "application/zip")
	c.Writer.Header().Set("Content-Disposition", "attachment; filename=\"all_files.zip\"")
	zipWriter := zip.NewWriter(c.Writer)
	defer zipWriter.Close()

	// List all objects in the bucket
	objectCh := minioClient.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{})
	for object := range objectCh {
		if object.Err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": object.Err.Error()})
			return
		}

		// Create a new file in the ZIP archive
		zipFile, err := zipWriter.Create(object.Key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Get the object from the bucket
		object, err := minioClient.GetObject(context.Background(), bucketName, object.Key, minio.GetObjectOptions{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer object.Close()

		// Copy the object data into the ZIP file
		if _, err := io.Copy(zipFile, object); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
}
