# Forum Project 🧑‍💻💬

## Overview 📚

This project is a **web forum** designed to facilitate communication between users, allowing them to create posts, comments, and interact through likes/dislikes. Users can filter posts by categories, created time, and user likes. The system leverages **SQLite** for database storage and **Docker** for containerization.

---

## Table of Contents 📑
- [Project Objectives 🎯](#project-objectives)
- [Features ✨](#features)
- [Database Structure 🗄️](#database-structure)
- [Authentication 🔑](#authentication)
- [Communication 💬](#communication)
- [Likes and Dislikes ❤️👎](#likes-and-dislikes)
- [Filtering Posts 🔍](#filtering-posts)
- [Docker Setup 🐳](#docker-setup)
- [File Structure 📁](#file-structure)
- [Usage 🛠️](#usage)
- [Contributors 🤝](#contributors)

---

## Project Objectives 🎯

The main objectives of this project are to:

- **Enable communication** between users via posts and comments.
- **Associate categories** with posts.
- Allow users to **like and dislike** posts and comments.
- **Filter posts** based on categories, creation date, and likes.

---

## Features ✨

### 📧 User Registration & Login
- Users can **register** and **login** securely.
- The system supports **sessions** using cookies (with an expiration date).
- The registration requires a **unique email**, **username**, and **password** (password is encrypted if implemented).

### 💬 Communication
- Registered users can create **posts** and **comments**.
- Posts can be categorized for better organization.

### ❤️👎 Likes & Dislikes
- Registered users can **like** or **dislike** posts and comments.
- The number of likes and dislikes is visible to all users.

### 🔍 Filtering Posts
- **Filter posts** based on:
  - Categories
  - Date created
  - User’s liked posts

### 🗄️ SQLite Database
- **SQLite** is used to store data such as users, posts, comments, and reactions (likes/dislikes).

---

## Database Structure 🗄️

The database is structured as follows:

- **Users**: Stores user information like email, username, and encrypted password.
- **Posts**: Stores posts with their title, content, and associated categories.
- **Comments**: Stores comments linked to posts.
- **Reactions**: Tracks likes and dislikes for posts and comments.
- **Categories**: Defines the available categories for posts.

---

## Authentication 🔑

To register, the user must provide:
- **Email** (must be unique).
- **Username**.
- **Password** (encrypted).

The system verifies user credentials during login using:
- **Email**
- **Password** (encrypted check).

---

## Communication 💬

- **Registered users** can create posts and add comments.
- **Posts** can be categorized using predefined categories (e.g., technology, sports).
- **Non-registered users** can **view posts and comments**, but cannot create or interact with them.

---

## Likes and Dislikes ❤️👎

- **Likes** and **Dislikes** are only available to registered users.
- The number of likes and dislikes is visible to everyone.
- Users can **like/dislike** both posts and comments.

---

## Filtering Posts 🔍

- Users can filter posts by:
  - **Categories** (e.g., tech, news)
  - **Posts created** by the logged-in user
  - **Posts liked** by the logged-in user

---

## Docker Setup 🐳

This project utilizes **Docker** for easy containerization and deployment.

### Prerequisites
- Install [Docker](https://www.docker.com/get-started).

### Running the Application

1. **Build the Docker image**:
    ```bash
    make build
    ```

2. **Run the Docker container**:
    ```bash
    make run
    ```

3. **Stop and clean up**:
    ```bash
    make stop
    make clean
    ```

4. **Push updates to Git**:
    ```bash
    make push
    ```

---

## File Structure 📁

Here’s an overview of the project’s directory structure:

```
.
├── cmd
│   └── main.go
├── db
│   └── data.db
├── internal
│   ├── database
│   │   ├── db.go
│   │   └── models
│   │       ├── categoriesTable.go
│   │       ├── commentsTable.go
│   │       ├── postsTable.go
│   │       ├── reactionsTable.go
│   │       └── usersTable.go
│   ├── handlers
│   │   ├── category.go
│   │   ├── comment.go
│   │   ├── login.go
│   │   ├── newPost.go
│   │   ├── post.go
│   │   ├── react.go
│   │   ├── register.go
│   │   └── user.go
│   ├── middleware
│   │   └── auth.go
│   └── utils
│       └── utils.go
|── web
|   ├── assets
|   │   ├── css
|   │   │   ├── base.css
|   │   │   ├── categories.css
|   │   │   ├── error.css
|   │   │   ├── login.css
|   │   │   ├── new_post.css
|   │   │   ├── posts.css
|   │   │   └── register.css
|   │   ├── icons
|   │   │   └── logo.png
|   │   └── js
|   │       ├── categories.js
|   │       ├── comments.js
|   │       ├── likes.js
|   │       ├── login.js
|   │       ├── posts.js
|   │       ├── register.js
|   │       └── script.js
|   ├── templates
|   │   ├── base.html
|   │   ├── categories.html
|   │   ├── error.html
|   │   ├── login.html
|   │   ├── new_post.html
|   │   ├── posts.html
|   │   ├── register.html
|   │   └── sideBar.html
|   └── templatesFunctions.go
├── Dockerfile
├── go.mod
├── go.sum
├── Makefile
└── README.md

15 directories, 47 files
```

---

## Usage 🛠️

### To run the project:

1. Clone the repository:
    ```bash
    git clone <repository_url>
    cd forum
    ```

2. Build and run using Docker:
    ```bash
    make docker
    ```

3. Open your browser and navigate to:
    ```
    http://localhost:8080
    ```

### Setting Environment Variables

Before running the application, make sure to set the following environment variables:

- **PORT**: The port where the application will run (default is `8080`).
- **DB_PATH**: The path to your SQLite database file (default is `db/data.db`).

You can set them using the following commands:

```bash
export PORT=8080
export DB_PATH=db/data.db
```

---

## Contributors 🤝

- **Oussama BENALI** 🧑‍💻
- **Ayoub El Haddad** 👨‍💻
- **Youssef Hajjaoui** 👨‍💻
- **Ilyass Atlassi** 👨‍💻
- **Ilyass Mohamed Foukahi** 👨‍💻

Thank you for your contributions! 💙
