import { makePost } from "./post.js";
import { getReactInfo,reactToggle } from "./likes.js";
import { commentToggle } from "./comments.js";


export async function RenderPost(posts) {
  // posts-container is inside the main element
  const container = document.querySelector(".posts-container");
  container.innerHTML = ""; // Empty the inner HTML if any

  for (const element of posts) {
    const post = document.createElement("div");
    post.classList.add("post");

    try {
      // Fetch reaction info
      const reactInfo = await getReactInfo({
        target_type: "post",
        target_id: element.PostId,
      }, "GET");


      // Populate the post's HTML
      post.innerHTML = makePost(element, reactInfo.data);

      // Add toggle functionalities
      let display_comment = false;
      commentToggle(post, element, display_comment);
      reactToggle(post, element.PostId, "post");

      // Append the post to the container
      container.append(post);
    } catch (error) {
      console.error("Error rendering post:", error);
    }
  }
}
