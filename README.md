### 2. Posts Management (Create, Read, Update, Delete Posts)
**youssef hajjaoui:**
- **SQL:**
  - Create the `posts` table with fields:
    - `post_id`
    - `user_id` (foreign key to `users`)
    - `title`
    - `content`
    - `created_at`
    - `updated_at`
- **Backend:**
  - Implement CRUD operations (Create, Read, Update, Delete) for posts.
  - Ensure only registered users can create or modify posts.
  - Handle pagination for displaying posts.
- **Frontend:**
  - Create HTML forms for creating and updating posts.
  - Display posts in a list format on the homepage.
- **Testing:**
  - Write unit tests for post creation and display.
# Project forum

One Paragraph of project description goes here

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

Run build make command with tests
```bash
make all
```

Build the application
```bash
make build
```

Run the application
```bash
make run
```

Live reload the application:
```bash
make watch
```

Run the test suite:
```bash
make test
```

Clean up binary from the last build:
```bash
make clean
```
