# Go Cloud Build Trigger Runner
This is a Go application for initiating Cloud Build triggers on Google Cloud Platform.

## Prerequisites
Before using this application, ensure you have the following in place:

* Google Cloud Account
* Google Cloud Project
* Google Cloud Authentication

## Configuration

Log in to gcloud auth:

```bash
gcloud auth application-default login
```


Clone this repository to your local machine:

```bash
git clone https://github.com/DanielSuhett/deploy-to-cloudbuild.git
```

Install dependencies using the following command:
```bash
go mod tidy
```
Build the application:
```bash
go build -o cloud-build-runner
```
Run the configuration:
```bash
./cloud-build-runner config 
```
 Input `YOUR_PROJECT` with the appropriate values for your project.

## Run

### Deploy
```bash
./cloud-build-runner deploy --trigger YOUR_TRIGGER --branch YOUR_BRANCH
```
Replace`YOUR_TRIGGER`, and `YOUR_BRANCH` with the appropriate values for your project.

---

### Status
```bash
./cloud-build-runner status --trigger YOUR_TRIGGER
```
Replace `YOUR_TRIGGER` with the appropriate values for your project.

## License
This project is licensed under the MIT License - see the LICENSE file for details.
