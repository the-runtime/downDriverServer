<div align="center">
<h1 align="center">
<img src="https://raw.githubusercontent.com/PKief/vscode-material-icon-theme/ec559a9f6bfd399b82bb44393651661b08aaf7ba/icons/folder-markdown-open.svg" width="100" />
<br>
serverDowndrive
</h1>
<h3 align="center">📍 Up Your Server with serverDowndrive!</h3>
<h3 align="center">🚀 Developed with the software and tools below:</h3>
<p align="center">

<img src="https://img.shields.io/badge/NOW-001211.svg?style=for-the-badge&logo=NOW&logoColor=white" alt="NOW" />
<img src="https://img.shields.io/badge/SVG-FFB13B.svg?style=for-the-badge&logo=SVG&logoColor=black" alt="SVG" />
<img src="https://img.shields.io/badge/Markdown-000000.svg?style=for-the-badge&logo=Markdown&logoColor=white" alt="Markdown" />
<img src="https://img.shields.io/badge/Go-00ADD8.svg?style=for-the-badge&logo=Go&logoColor=white" alt="Go" />
<img src="https://img.shields.io/badge/JavaScript-F7DF1E.svg?style=for-the-badge&logo=JavaScript&logoColor=black" alt="JavaScript" />
<img src="https://img.shields.io/badge/GitHub%20Actions-2088FF.svg?style=for-the-badge&logo=GitHub-Actions&logoColor=white" alt="GitHub%20Actions" />
<img src="https://img.shields.io/badge/SQLite-003B57.svg?style=for-the-badge&logo=SQLite&logoColor=white" alt="SQLite" />
</p>

</div>

---

## 📚 Table of Contents
- [📚 Table of Contents](#-table-of-contents)
- [📍Overview](#-introdcution)
- [🔮 Features](#-features)
- [⚙️ Project Structure](#project-structure)
- [🧩 Modules](#modules)
- [🏎💨 Getting Started](#-getting-started)
- [🗺 Roadmap](#-roadmap)
- [🤝 Contributing](#-contributing)
- [🪪 License](#-license)
- [📫 Contact](#-contact)
- [🙏 Acknowledgments](#-acknowledgments)

---


## 📍Overview

The serverDowndrive GitHub project is a powerful tool that provides users with an easy and secure way to download files from a server to their Google Drive. It offers features such as speed limit, progress tracking, database recording, and authentication, as well as a job dispatcher and job queue for workers. This project ensures files are transferred quickly and efficiently, and provides users with a secure and reliable way to store their data. It is a great solution for users looking for a convenient and secure way to transfer files.

---

## 🔮 Feautres

### Distinctive Features

1. **User-Centered Design:** The code scripts for serverDowndrive are designed with the user in mind. It includes features such as authentication, processing, and progress monitoring for the file download manager, as well as a file uploader for Google Drive with features such as speed limit, progress bar, database recording, and file deletion after upload.
2. **Flexible Architecture:** The project is designed with a flexible architecture that allows it to accommodate different types of users and their specific needs. It includes a JobHandler for downloading files from a given URL and uploading them to a Google Drive, a Worker type that can receive jobs from a JobQueue channel, and a dispatcher that manages a pool of worker threads.
3. **Data Storage:** The project also includes a number of code scripts for data storage and retrieval. It contains structs for storing user data from Google, user tokens, and user histories, as well as functions for fetching and resetting user data and retrieving a history table from a database.
4. **High Performance:** The project is optimized for high performance. It includes imports for various libraries and packages to facilitate the process, a WriteCounter struct for tracking progress, and a StartDown function to handle progress info and create files in the working directory. Additionally, it utilizes the bwlimit library for setting read and write limits.

---


<img src="https://raw.githubusercontent.com/PKief/vscode-material-icon-theme/ec559a9f6bfd399b82bb44393651661b08aaf7ba/icons/folder-github-open.svg" width="80" />

## ⚙️ Project Structure


```bash
repo
├── controller
│   └── progress.go
├── database
│   └── database.go
├── fileController
│   ├── down.go
│   └── up.go
├── go.mod
├── go.sum
├── handlers
│   ├── base.go
│   ├── oauth.go
│   ├── starter.go
│   └── user.go
├── main.go
├── model
│   ├── googleData.go
│   ├── model.go
│   ├── userHistory.go
│   └── userToken.go
├── README.md
├── todo.txt
└── workers
    ├── dispatcher.go
    ├── job.go
    └── worker.go

16 directories, 54 files
```

---

<img src="https://raw.githubusercontent.com/PKief/vscode-material-icon-theme/ec559a9f6bfd399b82bb44393651661b08aaf7ba/icons/folder-src-open.svg" width="80" />

## 💻 Modules

<details closed><summary>Controller</summary>

| File        | Summary                                                                                                                                                                                                                                     | Module                 |
|:------------|:--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|:-----------------------|
| progress.go | This code script defines a structure for tracking progress for a file transfer and includes functions for creating a new progress entry, retrieving progress by user ID and progress ID, and getting a list of progress entries by user ID. | controller/progress.go |

</details>

<details closed><summary>Filecontroller</summary>

| File    | Summary                                                                                                                                                                                                                                                                                                                                                               | Module                 |
|:--------|:----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|:-----------------------|
| up.go   | This code script implements a file uploader for Google Drive, with features such as speed limit, progress bar, database recording, and file deletion after upload. It also includes imports of various libraries and packages to facilitate the process.                                                                                                              | fileController/up.go   |
| down.go | This code script provides a file controller package for downloading files from a server. It includes functions to implement transfer limit on users, a WriteCounter struct for tracking progress, and a StartDown function to handle progress info and create files in the working directory. It also utilizes the bwlimit library for setting read and write limits. | fileController/down.go |

</details>

<details closed><summary>Handlers</summary>

| File       | Summary                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             | Module              |
|:-----------|:----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|:--------------------|
| user.go    | This code script contains handlers for user registration, fetching and resetting user data, and retrieving a history table from a database. It imports the necessary packages to access a user's Google account data and manipulate a database. After parsing the form data, user parameters such as data transfer limit are set based on account type. The user data is then retrieved from Google and encoded in JSON format to be sent back to the client. The history table is also retrieved from the database and encoded in the same format. | handlers/user.go    |
| starter.go | This code script provides handlers for a download service that allows users to transfer files from a URL to their Google Drive. It includes functions to start the download, track download progress, and authenticate user accounts.                                                                                                                                                                                                                                                                                                               | handlers/starter.go |
| oauth.go   | This package handles the authentication of a user utilizing Google OAuth. It starts by importing packages and establishing a Oauth2 Config. It then creates a function for the user to log in, and another to handle the callback from Google. It also sets the user database and token database for the user. Finally, it gets the user data from Google and sets a cookie for the user.                                                                                                                                                           | handlers/oauth.go   |
| base.go    | This code script sets up a server to handle requests for a file download manager, including authentication, processing, and progress monitoring. It also creates a dispatcher and job queue for workers, as well as a file server for the API.                                                                                                                                                                                                                                                                                                      | handlers/base.go    |

</details>

<details closed><summary>Model</summary>

| File           | Summary                                                                                                                                                                                                                                                   | Module               |
|:---------------|:----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|:---------------------|
| userToken.go   | This code script defines a UserToken struct for use with the GORM library. It contains fields for userId, accessToken, refreshToken, authCode, tokenType, expiry, and a foreign key reference to a Token object.                                          | model/userToken.go   |
| userHistory.go | This code script defines a struct'SingleHistory' with fields for storing data about an individual file, such as user ID, file name, file size, and start/end times. It also includes a commented-out struct for storing multiple files in a user history. | model/userHistory.go |
| model.go       | This code script defines a'User' model that stores relevant information such as user id, name, account type, data transfer allowance, consumed data transfer, allowed speed, and allowed threads.                                                         | model/model.go       |
| googleData.go  | This code script defines a struct "GoogleUserData" which contains user data from Google, such as their ID, Email, Verified Email status, and Picture.                                                                                                     | model/googleData.go  |

</details>

<details closed><summary>Root</summary>

| File    | Summary                                                                                                                                                                                                                                                                                          | Module   |
|:--------|:-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|:---------|
| main.go | This code script is a simple HTTP server for downloading files that listens to port 8000 and uses a handler package for routing requests. It prints statements to the log about its status and closes the server when finished.                                                                  | main.go  |
| go.mod  | This code script is for serverFordownDrive, which requires various libraries for its functionality such as bwlimit, oauth2, and Google APIs. Additionally, the script requires several indirect libraries, such as cloud.google.com/go/compute, golang.org/x/crypto, and google.golang.org/grpc. | go.mod   |

</details>

<details closed><summary>Workers</summary>

| File          | Summary                                                                                                                                                                                                                                                                                                                                                                                | Module                |
|:--------------|:---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|:----------------------|
| worker.go     | This script creates a Worker type that can receive jobs from a JobQueue channel, execute the job, and report errors if needed. It also contains Start and Stop functions that allow the Worker to join and leave its WorkerPool.                                                                                                                                                       | workers/worker.go     |
| job.go        | This code script creates a JobHandler for downloading files from a given URL and uploading them to a Google Drive. It includes imports for context, oauth2, controller, database, and fileController, as well as a Job structure and a NewJob function. Additionally, it includes a DoJob function to handle the token exchange, download the file, and upload it to the Google Drive. | workers/job.go        |
| dispatcher.go | This code script creates a dispatcher that manages a pool of worker threads, allowing jobs to be assigned to them and processes to be completed. It establishes a maximum number of workers that can be used and provides a dispatch function that assigns jobs to the available workers.                                                                                              | workers/dispatcher.go |

</details>

<hr />

## 🚀 Getting Started

### ✅ Prerequisites

Before you begin, ensure that you have the go version 20 installed in your system
<!-- > `[📌  INSERT-PROJECT-PREREQUISITES]`
 -->
### 💻 Installation

1. Clone the serverDowndrive repository:
```sh
git clone https://github.com/the-runtime/serverDowndrive.git
```

2. Change to the project directory:
```sh
cd serverDowndrive
```

3. Install the dependencies:
```sh
go build -o myapp
```

### 🤖 Using serverDowndrive

```sh
./myapp
```

### 🧪 Running Tests
```sh
#run tests
```

<hr />


<!-- ## 🛠 Future Development
- [X] [📌  COMPLETED-TASK]
- [ ] [📌  INSERT-TASK]
- [ ] [📌  INSERT-TASK] -->


---

## 🤝 Contributing
Contributions are always welcome! Please follow these steps:
1. Fork the project repository. This creates a copy of the project on your account that you can modify without affecting the original project.
2. Clone the forked repository to your local machine using a Git client like Git or GitHub Desktop.
3. Create a new branch with a descriptive name (e.g., `new-feature-branch` or `bugfix-issue-123`).
```sh
git checkout -b new-feature-branch
```
4. Make changes to the project's codebase.
5. Commit your changes to your local branch with a clear commit message that explains the changes you've made.
```sh
git commit -m 'Implemented new feature.'
```
6. Push your changes to your forked repository on GitHub using the following command
```sh
git push origin new-feature-branch
```
7. Create a pull request to the original repository.
Open a new pull request to the original project repository. In the pull request, describe the changes you've made and why they're necessary.
The project maintainers will review your changes and provide feedback or merge them into the main branch.

---

## 🪪 License

This project is licensed under the `MIT License` License. 
<!-- See the [LICENSE](https://docs.github.com/en/communities/setting-up-your-project-for-healthy-contributions/adding-a-license-to-a-repository) file for additional info.
 -->
---
<!-- 
## 🙏 Acknowledgments

[📌  INSERT-DESCRIPTION]


---
 -->
