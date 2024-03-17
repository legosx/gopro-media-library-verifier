# GoPro Media Library Verifier

![Coverage Status](https://coveralls.io/repos/github/legosx/gopro-media-library-verifier/badge.svg?branch=main)

This CLI tool will make sure that your precious Gopro media is safely uploaded to the cloud.
You can also use it for other files that are not created by Gopro camera.

* verifies that your local files are already uploaded to https://plus.gopro.com/media-library/
* identifies which files are not yet uploaded

## Workflow

1. Upload your media to https://plus.gopro.com/media-library/ using the browser or Gopro apps.
2. Run this tool to verify that all your media is uploaded.
3. If any files are missing, upload them manually.

## Installation

1. You need to install golang first:
   https://go.dev/doc/install

2. Then run:

```bash
cd $GOPATH/bin
go install github.com/legosx/gopro-media-library-verifier@main
```

## Usage

### Getting API token

1. Login to https://plus.gopro.com/media-library/ or just open the page if you are already logged in.
2. Open developer tools of your browser.
3. Lookup for "search" requests in Network tab. If the results are empty, refresh or scroll the page - new request to API should go out. 
You should pick the request with `Request method: GET` and `Status Code: 200 OK`.
4. From here on, you have 2 options to specify the token for the tool:
   1. Click on the request and go to the Request Headers section. You should see "Authorization:" header with the value
   starting with "Bearer ". Copy the value after "Bearer " and use it as a token for the tool.
   2. You can just do right mouse click and select "Save as cURL". Later you can paste it in the tool.

### Run!

There are multiple ways to specify the token for the tool:

1. Just run it and it will ask you for the token:
```
gopro-media-library-verifier verify -p /path/to/your/media
? Select token prompt method: 
  ▸ Direct input
    CURL request
```
if you fetched the token value yourself, you go for "Direct input".
If you saved the CURL request, you can use "CURL request".

If the token valid, it will be saved to `~/.gopro-media-library-verifier.json` file.

2. Use `AUTH.TOKEN=your token` environment variable. The token won't be saved to the file in this case.

#### Example

```
gopro-media-library-verifier verify -p /gopro
No token found in config
✔ Direct input

Please provide your Gopro Media Library token to authenticate:
✔ Token: ******************
Token is valid
Token saved to a new config file

Identifying files that are not yet uploaded to cloud from
/gopro
based on: fileName, fileSize

Files that still can be uploaded to Gopro Media Library:
/gopro/GX014882.MP4
```

which means that you can upload `/gopro/GX014882.MP4` to Gopro Media Library and run this command again to make sure that it's uploaded.

#### Output results in a file

if you want to save the results in a file, you can use the `-o` flag:

```bash
gopro-media-library-verifier verify -p /path/to/your/media -o /path/to/output/file
```

#### Always use the same token prompt method

If you want to always use the same token prompt method and don't show other options, you can use the `-m` flag:
```
gopro-media-library-verifier verify -p /path/to/your/media -m direct
gopro-media-library-verifier verify -p /path/to/your/media -m curl
```
