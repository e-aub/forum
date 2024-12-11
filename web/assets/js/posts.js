// import { getReactInfo, initializeLikeButtons } from "./likes.js";
import { initializeCommentSection } from "./comments.js";
import { getReactInfo } from "./likes.js";


export async function renderPosts(posts) {
  const postsContainer = document.querySelector(".posts");
  postsContainer.innerHTML = "";

  for (const post of posts) {
    const postElement = document.createElement("div");
    postElement.classList.add("post");

    try {
      const reactInfo = await getReactInfo(post.PostId);
      postElement.innerHTML = generatePostHTML(post, reactInfo);
      postsContainer.appendChild(postElement);

      // Initialize likes and comments
      initializeLikeButtons(postElement, post);
      initializeCommentSection(postElement, post);
    } catch (error) {
      console.error("Error rendering post:", error);
    }
  }
}

function generatePostHTML(post, reactInfo) {
  return `
  <div class="post-header">
      <h3>${post.Title}</h3>
      <span>${new Date(post.CreatedAt).toLocaleDateString()}</span>
  </div>
  <div class="post-body">
      <p>${post.Content}</p>
  </div>
  <div class="post-footer">
      <button class="like like-button" data-clicked="false">👍 Like (<span class="count">${reactInfo.likes || 0}</span>)</button>
      <button class="dislike dislike-button" data-clicked="false">👎 Dislike (<span class="count">${reactInfo.dislikes || 0}</span>)</button>
      <button class="toggle-comments">Show Comments</button>
  </div>
  <div class="comments-section" style="display: none;">
      <textarea placeholder="Add a comment..." class="comment-input"></textarea>
      <button class="comment-submit">Submit</button>
      <div class="comments"></div>
  </div>
`;
}
