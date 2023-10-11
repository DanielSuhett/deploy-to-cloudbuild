# Go Cloud Build Trigger Runner
This is a Go application for initiating Cloud Build triggers on Google Cloud Platform.

## Prerequisites
Before using this application, ensure you have the following in place:

* Google Cloud Account
* Google Cloud Project
* Google Cloud Authentication


## Configuration

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

Run the application with the following arguments:
```bash
./cloud-build-runner deploy --project YOUR_PROJECT --trigger YOUR_TRIGGER --branch YOUR_BRANCH
```
 Replace `YOUR_PROJECT`, `YOUR_TRIGGER`, and `YOUR_BRANCH` with the appropriate values for your project.

## License
This project is licensed under the MIT License - see the LICENSE file for details.
