import { getReactInfo, reactToggle } from "./likes.js";
import { showRegistrationModal } from "./script.js";

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

  const commentInput = postElement.querySelector(".comment-input");
  commentInput.addEventListener("keydown", async (event) => {
    if (event.key === "Enter" && !event.shiftKey) {
      const commentInput = postElement.querySelector(".comment-input");
      if (commentInput.value.trim()) {
        await addComment(post.PostId, commentInput.value.trim(), commentsSection.querySelector(".comments"), commentsSection);
        commentInput.value = ""
      }
    }
  })
  commentInput.addEventListener('input', () => {
    console.log(commentInput.scrollHeight);
    
    commentInput.style.height = 'auto'
    commentInput.style.height = commentInput.scrollHeight + "px"
  });
}

async function loadComments(postId, commentsContainer) {
  commentsContainer.innerHTML = ""
  try {
    const response = await fetch(`/comments?post=${postId}`);
    if (!response.ok) throw new Error("Failed to load comments.");
    const comments = await response.json();
    if (!comments)return
    comments.forEach(async comment => {
      const reaction = await getReactInfo({
        target_type: "comment",
        target_id: comment.comment_id,
      }, "GET")
      const commentSection = createCommentElement(comment, reaction)
      reactToggle(commentSection, comment.comment_id, 'comment')
      commentsContainer.appendChild(commentSection)
    });
  } catch (error) {
    console.error("Error loading comments:", error);
  }
}

async function addComment(postId, content, commentsContainer, commentsection) {
  try {
    const response = await fetch(`/comments`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        post_id: postId,
        content: content
      }),
    });

    switch (response.status) {
      case 400:
        const error = commentsection.querySelector('.error-comment')
        error.textContent = 'comment must be only 150 character'
        break;
      case 401:
        showRegistrationModal()
        break;
      case 201:
        const newComment = await response.json();
        const reaction = await getReactInfo({
          target_type: "comment",
          target_id: newComment.comment_id,
        }, "GET")
        const commentSection = createCommentElement(newComment, reaction)
        reactToggle(commentSection, newComment.comment_id, 'comment')
        commentsContainer.prepend(commentSection)
        break
    }
  } catch (error) {
    console.error("Error adding comment:", error);
  }
}

function createCommentElement(comment, reaction) {
  const commentElement = document.createElement("div");
  commentElement.classList.add("comment");
  let liked = false;
  let disliked = false;
  let likeCount = reaction.data.liked_by ? reaction.data.liked_by.length : 0;
  let disLikeCount = reaction.data.disliked_by ? reaction.data.disliked_by.length : 0;

  if (!!reaction.data.user_reaction) {
    liked = reaction.data.user_reaction === "like"
    disliked = !liked
  } else {
    liked = false
    disliked = false;
  }

  commentElement.innerHTML = `
      <p><strong>👤 ${comment.user_name}:</strong> ${comment.content}</p>
      <div class="reaction-section comment-likes">
        <button class="like-button ${liked ? 'clicked' : ''}" data-clicked=${liked}>
          👍 Like (<span class="count">${likeCount}</span>)
        </button>
        <button class="dislike-button ${disliked ? 'clicked' : ''}" data-clicked=${disliked}>
          👎 Dislike (<span class="count">${disLikeCount}</span>)
        </button>
      </div>
    `;
  return commentElement;
}