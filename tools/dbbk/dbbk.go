package main

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"golang.org/x/crypto/ssh"

	"cloud.google.com/go/storage"
)

/*
	# 処理フロー
	Cloud Schedulerでトピックdb-bkに対し定期パブリッシュ(12h毎)
	↓
	Cloud Pub/Subのトピックdb-bkをサブスクライブしている以下CloudFunctionのRun関数が実行される
	↓
	gcsから秘密鍵取得
	↓
	ssh
	↓
	リモートPCでdump
	↓
	gcsへアップロード

	# ローカルから実行
	GOOGLE_APPLICATION_CREDENTIALS=~/.config/gcloud/legacy_credentials/lbfdeatq\@gmail.com/adc.json go run dbbk.go

	# ref: https://godoc.org/cloud.google.com/go/storage
*/
func main() {
	if err := Run(); err != nil {
		log.Fatalf("failed to DB BK: %+v", err)
	}
}

func Run() error {
	log.Println("DATABASE BACKUP START")

	s, err := dumpDB()
	if err != nil {
		return err
	}

	if err := uploadDBDump(s); err != nil {
		return err
	}

	log.Println("DATABASE BACKUP END")
	return nil
}

func dumpDB() (string, error) {
	ctx := context.Background()
	gcsClient, err := storage.NewClient(ctx)
	if err != nil {
		return "", err
	}
	bucketName := "tag-mng-243823.appspot.com"
	bkt := gcsClient.Bucket(bucketName)

	objName := "dbbk/secret"
	obj := bkt.Object(objName)

	// Read it back.
	r, err := obj.NewReader(ctx)
	if err != nil {
		return "", err
	}
	defer r.Close()

	key := make([]byte, r.Size())
	_, err = r.Read(key)
	if err != nil {
		return "", err
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return "", err
	}

	config := &ssh.ClientConfig{
		User: "lbfdeatq",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", "ddddddo.work:22", config)
	if err != nil {
		return "", err
	}
	defer client.Close()

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	var b bytes.Buffer
	session.Stdout = &b

	cmd := "pg_dump tag-mng"
	if err := session.Run(cmd); err != nil {
		return "", err
	}

	return b.String(), nil
}

func uploadDBDump(s string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	bucketName := "tag-mng-243823.appspot.com"
	bkt := client.Bucket(bucketName)

	objName := "dbbk/dump"
	obj := bkt.Object(objName)

	// Write something to obj.
	// w implements io.Writer.
	w := obj.NewWriter(ctx)
	// Write some text to obj. This will either create the object or overwrite whatever is there already.
	if _, err := fmt.Fprintf(w, s); err != nil {
		return err
	}
	// Close, just like writing a file.
	if err := w.Close(); err != nil {
		return err
	}

	return nil
}

/*
	実際のコード

// Package p contains a Pub/Sub Cloud Function.
package p

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"golang.org/x/crypto/ssh"

	"cloud.google.com/go/storage"
)

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// RunPubSub consumes a Pub/Sub message.
func RunPubSub(ctx context.Context, m PubSubMessage) error {
	log.Println(string(m.Data))

    if err := Run(); err != nil {
		log.Fatalf("failed to DB BK: %+v", err)
	}

	return nil
}

func Run() error {
	log.Println("DATABASE BACKUP START")

	s, err := dumpDB()
	if err != nil {
		return err
	}

	if err := uploadDBDump(s); err != nil {
		return err
	}

	log.Println("DATABASE BACKUP END")
	return nil
}

func dumpDB() (string, error) {
	ctx := context.Background()
	gcsClient, err := storage.NewClient(ctx)
	if err != nil {
		return "", err
	}
	bucketName := "tag-mng-243823.appspot.com"
	bkt := gcsClient.Bucket(bucketName)

	objName := "dbbk/secret"
	obj := bkt.Object(objName)

	// Read it back.
	r, err := obj.NewReader(ctx)
	if err != nil {
		return "", err
	}
	defer r.Close()

	key := make([]byte, r.Size())
	_, err = r.Read(key)
	if err != nil {
		return "", err
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return "", err
	}

	config := &ssh.ClientConfig{
		User: "lbfdeatq",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", "ddddddo.work:22", config)
	if err != nil {
		return "", err
	}
	defer client.Close()

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	var b bytes.Buffer
	session.Stdout = &b

	cmd := "pg_dump tag-mng"
	if err := session.Run(cmd); err != nil {
		return "", err
	}

	return b.String(), nil
}

func uploadDBDump(s string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	bucketName := "tag-mng-243823.appspot.com"
	bkt := client.Bucket(bucketName)

	objName := "dbbk/dump"
	obj := bkt.Object(objName)

	// Write something to obj.
	// w implements io.Writer.
	w := obj.NewWriter(ctx)
	// Write some text to obj. This will either create the object or overwrite whatever is there already.
	if _, err := fmt.Fprintf(w, s); err != nil {
		return err
	}
	// Close, just like writing a file.
	if err := w.Close(); err != nil {
		return err
	}

	return nil
}

*/
