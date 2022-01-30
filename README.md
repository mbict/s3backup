#S3Backup

A really simple backup uploader for S3.
It supports deleting odler backup files if x amount of files exceed the files threshold.

Make sure the following environment variables are set:

- AWS_SECRET_ACCESS_KEY
- AWS_ACCESS_KEY_ID
- AWS_REGION
- AWS_ENDPOINT

And run the application with:
```
s3backup bucketname my/backup/dir /home/file/to/upload.txt
```

Or supply a number as the last param with max allowed files in the directory
```
s3backup bucketname my/backup/dir /home/file/to/upload.txt 3
```

Example to run it from docker with mounting a local volume
```
docker run -it \                                                                                                                              ✔  20:09:51 
  -v $(pwd)/:/data \
  --env AWS_SECRET_ACCESS_KEY=s3+access+key+here \
  --env AWS_ACCESS_KEY_ID=secretkey \
  --env AWS_REGION=ams3 \
  --env AWS_ENDPOINT=https://ams3.digitaloceanspaces.com \
  mbict/s3backup \
  mbict git-backup /data/README.md
```