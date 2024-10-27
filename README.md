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