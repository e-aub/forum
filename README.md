# Web Forum Project

This project involves building a web forum with SQLite for data storage, user authentication, and Docker integration. The forum will include features such as user registration, posting, commenting, liking/disliking, and filtering posts by categories.

## Team Roles and Responsibilities

The project is divided into multiple sections to ensure that each team member contributes to both the backend (SQL and Go) and frontend (HTML) aspects of the project. Below are the responsibilities assigned to each person.

### 1. User Management (Registration, Login, Profile)
**ilyass atlassi**
- **SQL:**
  - Create the `users` table with fields:
    - `id`
    - `email`
    - `username`
    - `password`
    - `session tracking`
- **Backend:**
  - Implement registration and login functionality.
  - Set up bcrypt for password encryption.
  - Handle session creation using cookies and UUID (bonus).
- **Frontend:**
  - Create simple HTML pages for user registration, login, and profile (e.g., displaying their own posts, likes).
- **Testing:**
  - Write unit tests for registration and login functionalities.

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

### 3. Comments Management (Comments on Posts)
**ilyass foukahi:**
- **SQL:**
  - Create the `comments` table with fields:
    - `comment_id`
    - `post_id` (foreign key to `posts`)
    - `user_id` (foreign key to `users`)
    - `content`
    - `created_at`
    - `updated_at`
- **Backend:**
  - Implement CRUD operations for comments (only registered users can comment).
  - Ensure comments are tied to specific posts and displayed in order.
- **Frontend:**
  - Create HTML forms for adding comments below posts.
  - Display comments under each post for both registered and non-registered users.
- **Testing:**
  - Write unit tests for comment functionality.

### 4. Likes and Dislikes (Posts and Comments)
**oussama benali:**
- **SQL:**
  - Create the `likes` table with fields:
    - `like_id`
    - `user_id`
    - `post_id` or `comment_id`
    - `type` (like or dislike)
- **Backend:**
  - Implement functionality for liking or disliking posts and comments (only registered users).
  - Ensure the number of likes and dislikes is visible to all users.
- **Frontend:**
  - Add like/dislike buttons next to posts and comments.
  - Display the current count of likes/dislikes.
- **Testing:**
  - Write unit tests for like/dislike functionality.

### 5. Filtering and Categories (Subforums, Filtering by Posts, Likes)
**ayoub elhadad:**
- **SQL:**
  - Create the `categories` table with fields:
    - `category_id`
    - `name`
    - `description`
  - Create a join table `post_categories` to associate posts with categories.
- **Backend:**
  - Implement filtering mechanisms to display posts by category, user-created posts, and liked posts.
  - Handle filtering logic in the backend.
- **Frontend:**
  - Create a category page where users can filter posts based on category.
  - Add filtering options for posts created by the logged-in user and liked posts.
- **Testing:**
  - Write unit tests for filtering functionality.

## Additional Responsibilities (Shared Tasks)

### Error Handling and HTTP Status Codes (everyone contributes):
- Ensure that all responses have the appropriate HTTP status codes (e.g., 404 for not found, 500 for server errors).
- Handle database and validation errors gracefully and display useful error messages to users.

### Docker Integration (pair work):
- Two team members will work together to:
  - Create a `Dockerfile` and ensure that the entire application, including the database, runs inside a container.
  - Create Docker images and scripts to easily set up and run the forum.
  - Ensure that Docker is correctly set up for the SQLite database and other dependencies.

### Unit Testing (collaborative effort):
- Each team member is responsible for writing unit tests for the parts of the project they work on (registration, posts, comments, likes, etc.).
- Coordinate with each other to ensure all major components are covered by tests.
