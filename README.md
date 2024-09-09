# DICOM Service

A rudimentary REST-based API service that:

- stores DICOM files
- queries DICOM header attributes by tag
- generates a PNG representation of the DICOM file

This is not meant to conform with the DICOMweb™ standard

## How to use

Run the server with:

```
go run main.go
```

Run tests with:

```
go test ./...
```

### Uploading

Upload with the following curl command

```sh
curl -F instance=@<file_path> http://localhost:8080/instance
```

This will return a json response:

```json
{ "instanceId": "1.3.12.2.1107.5.2.6.24119.30000013121716094326500000535" }
```

### Query a header attribute

Query a tag on a particular instance with the following curl command

```sh
curl http://localhost:8080/instance/1.3.12.2.1107.5.2.6.24119.30000013121716094326500000535?tag=%280008%2C0020%29
```

If the tag is set on the instance, this will return the following json response:

```json
{
  "tag": { "Group": 8, "Element": 32 },
  "VR": 12,
  "rawVR": "DA",
  "valueLength": 8,
  "value": ["20131217"]
}
```

### Get PNG image of DICOM file

```sh
curl -i --output <output_file> http://localhost:8080/instance/1.3.12.2.1107.5.2.6.24119.30000013121716094326500000535/image
```
