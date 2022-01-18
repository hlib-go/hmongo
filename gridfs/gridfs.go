package gridfs

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"io"
)

// 文档：https://godoc.org/go.mongodb.org/mongo-driver/mongo/gridfs

// Upload 上传文件
func Upload(db *mongo.Database, filename string, fileBytes []byte) (fileId string, err error) {
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		return
	}
	uploadStream, err := bucket.OpenUploadStream(filename)
	defer func() {
		if uploadStream != nil {
			err = uploadStream.Close()
		}
	}()
	if err != nil {
		return
	}
	_, err = uploadStream.Write(fileBytes)
	if err != nil {
		return
	}
	fileId = (uploadStream.FileID).(primitive.ObjectID).Hex()
	return
}

// Download 下载文件
func Download(db *mongo.Database, fileId string) (bytes []byte, err error) {
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		return
	}
	objid, err := primitive.ObjectIDFromHex(fileId)
	if err != nil {
		return
	}
	downloadStream, err := bucket.OpenDownloadStream(objid)
	if err != nil {
		return
	}
	for {
		var block = make([]byte, 1024)
		i, err := downloadStream.Read(block)
		if err == io.EOF && i == 0 {
			break
		}
		if err != nil {
			return nil, err
		}
		bytes = append(bytes, block...)
	}
	err = downloadStream.Close()
	return
}

// Delete 删除文件
func Delete(db *mongo.Database, fileId string) (err error) {
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		return
	}
	fileID, err := primitive.ObjectIDFromHex(fileId)
	if err != nil {
		return
	}
	err = bucket.Delete(fileID)
	return
}
