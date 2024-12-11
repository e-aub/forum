import { initializeCommentSection } from "./comments.js";
import { getReactInfo, reactToggle } from "./likes.js";

export async function renderPosts(posts) {
  const postsContainer = document.querySelector(".posts");
  postsContainer.innerHTML = "";

  for (const post of posts) {
    const postElement = document.createElement("div");
    postElement.classList.add("post");

    try {
      const reactInfo = await getReactInfo({
        target_type: "post",
        target_id: post.PostId,
      }, "GET");
      postElement.innerHTML = generatePostHTML(post, reactInfo);
      postsContainer.appendChild(postElement);

      // Initialize likes and comments
      reactToggle(postElement, post.PostId, "post");
      initializeCommentSection(postElement, post);
    } catch (error) {
      console.error("Error rendering post:", error);
    }
  }
}


function generatePostHTML(post, reactInfo) {
  let liked = false;
  let disliked = false;
  let likeCount = reactInfo.data.liked_by ? reactInfo.data.liked_by.length : 0;
  let disLikeCount = reactInfo.data.disliked_by ? reactInfo.data.disliked_by.length : 0;
  console.log(reactInfo.data.liked_by, reactInfo);

  if (!!reactInfo.data.user_reaction) {
    liked = reactInfo.data.user_reaction === "like"
    disliked = !liked
  } else {
    liked = false
    disliked = false;
  }
  console.log(likeCount, disLikeCount);

  return `
  <div class="post-header">
      <h3>${post.Title}</h3>
      <h3>${post.UserName}</h3>
      <span>${new Date(post.CreatedAt).toLocaleDateString()}</span>
  </div>
  <div class="post-body">
      <p>${post.Content}</p>
  </div>
  <div class="post-footer">
      <button class="like like-button" data-clicked="${liked}">👍 Like (<span class="count">${likeCount}</span>)</button>
      <button class="dislike dislike-button" data-clicked="${disliked}">👎 Dislike (<span class="count">${disLikeCount}</span>)</button>
      <button class="toggle-comments">Show Comments</button>
  </div>
  <div class="comments-section" style="display: none;">
      <textarea placeholder="Add a comment..." class="comment-input"></textarea>
      <button class="comment-submit">Submit</button>
      <div class="comments"></div>
  </div>
`;
}