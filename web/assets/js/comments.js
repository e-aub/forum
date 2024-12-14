import { reactToggle } from "/assets/js/likes.js";

export function initializeCommentSection(postElement, post) {
  const toggleCommentsButton = postElement.querySelector(".toggle-comments");
  const commentsSection = postElement.querySelector(".comments-section");

  toggleCommentsButton.addEventListener("click", async () => {
    if (commentsSection.style.display === "none") {
      commentsSection.style.display = "block";
      toggleCommentsButton.textContent = "Hide Comments";
      await loadComments(post.PostId, commentsSection.querySelector(".comments"));
    } else {
      commentsSection.style.display = "none";
      toggleCommentsButton.textContent = "💬 Show Comments";
    }
  });

  const commentSubmitButton = postElement.querySelector(".comment-submit");
  commentSubmitButton.addEventListener("click", async () => {
    const commentInput = postElement.querySelector(".comment-input");
    if (commentInput.value.trim()) {
      const succes = await addComment(post.PostId, commentInput.value.trim(), commentsSection.querySelector(".comments"),commentsSection);
      succes ? commentInput.value = "" : alert('do not play in wrong place')
    }
  });

  const commentInput = postElement.querySelector(".comment-input");
  commentInput.addEventListener("keydown", async (event) => {
    if (event.key === "Enter") {
      const commentInput = postElement.querySelector(".comment-input");
      if (commentInput.value.trim()) {
        await addComment(post.PostId, commentInput.value.trim(), commentsSection.querySelector(".comments"),commentsSection);
        commentInput.value = ""
      }
    }
  })
}

async function loadComments(postId, commentsContainer) {
  commentsContainer.innerHTML = ""
  try {
    const response = await fetch(`/comments?post=${postId}`);
    if (!response.ok) throw new Error("Failed to load comments.");
    const comments = await response.json();
    comments.forEach(comment => commentsContainer.appendChild(createCommentElement(comment)));
  } catch (error) {
    console.error("Error loading comments:", error);
  }
}

async function addComment(postId, content, commentsContainer,commentsection) {
  try {
    const response = await fetch(`/comments`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        post_id: postId,
        content: content
      }),
    });
    if (response.status == 400) {
      const error = commentsection.querySelector('.error-comment')
      console.log(error);
      error.textContent = 'comment must be only 150 character'
      return      
    }
    const newComment = await response.json();
    console.log(newComment)
    const comment = createCommentElement(newComment);
    reactToggle(comment, newComment.comment_id, "comment");
    commentsContainer.appendChild(comment);
  } catch (error) {
    console.error("Error adding comment:", error);
  }
}

function createCommentElement(comment) {
  const commentElement = document.createElement("div");
  commentElement.classList.add("comment");
  commentElement.innerHTML = `
      <p><strong>👤 ${comment.user_name}:</strong> ${comment.content}</p>
      <div class="comment-likes">
        <button class="like-button">👍 Like (<span class="count">${0}</span>)</button>
        <button class="dislike-button">👎 Dislike (<span class="count">${0}</span>)</button>
      </div>
    `;
  return commentElement;
}