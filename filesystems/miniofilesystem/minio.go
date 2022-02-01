package miniofilesystem

import (
	"context"
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/youngjae-lim/gosnel/filesystems"
)

type Minio struct {
	Endpoint string
	Key      string
	Secret   string
	UseSSL   bool
	Region   string
	Bucket   string
}

// getCredentials gets a minio client.
func (m *Minio) getCredentials() *minio.Client {
	client, err := minio.New(m.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(m.Key, m.Secret, ""),
		Secure: m.UseSSL,
	})
	if err != nil {
		log.Println(err)
	}

	return client
}

// Put puts a file to a Minio remote file system.
func (m *Minio) Put(fileName, folder string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// get a file name - strip off a path before actual file name
	objectName := path.Base(fileName)

	// get a minio client
	client := m.getCredentials()

	uploadInfo, err := client.FPutObject(ctx, m.Bucket, fmt.Sprintf("%s/%s", folder, objectName), fileName, minio.PutObjectOptions{})
	if err != nil {
		log.Println("Failed with FPutObject")
		log.Println(err)
		log.Println("UploadInfo:", uploadInfo)
		return err
	}

	return nil
}

// List returns a list of all contents of a given bucket in the minio remote file system.
func (m *Minio) List(prefix string) ([]filesystems.Listing, error) {
	var listing []filesystems.Listing

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// get a minio client
	client := m.getCredentials()

	// get channels of ObjectInfo
	objectCh := client.ListObjects(ctx, m.Bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	// loop through ObjectInfo channels
	for object := range objectCh {
		if object.Err != nil {
			fmt.Println(object.Err)
			return listing, object.Err
		}

		if !strings.HasPrefix(object.Key, ".") {
			// convert file size(byte) to megabyte
			b := float64(object.Size)
			kb := b / 1024
			mb := kb / 1024

			// set Listing
			item := filesystems.Listing{
				Etag:         object.ETag,
				LastModified: object.LastModified,
				Key:          object.Key,
				Size:         mb,
			}

			// append the item to a slice of Listing
			listing = append(listing, item)
		}
	}

	return listing, nil
}

// Delete deletes a list of files in a given bucket of the minio remote file system.
func (m *Minio) Delete(itemsToDelete []string) bool {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// get a minio client
	client := m.getCredentials()

	opts := minio.RemoveObjectOptions{
		// set the bypass governance header to delete an object locked with GOVERNANCE mode
		GovernanceBypass: true,
	}

	for _, item := range itemsToDelete {
		err := client.RemoveObject(ctx, m.Bucket, item, opts)
		if err != nil {
			fmt.Println(err)
			return false
		}
	}

	return true
}

// Get downloads contents of an object from a remote file system to the server where myapp is running
func (m *Minio) Get(destination string, items ...string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// get a minio client
	client := m.getCredentials()

	for _, item := range items {
		err := client.FGetObject(ctx, m.Bucket, item, fmt.Sprintf("%s/%s", destination, path.Base(item)), minio.GetObjectOptions{})
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}
