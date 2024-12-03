import { getComment } from "./script.js";

const container = document.querySelector(".posts-container");
export function RenderPost(posts) {
  container.innerHTML = "";

  posts.forEach((element, index) => {
    const post = document.createElement('div');
    post.classList.add('post');
    console.log(element)
    const createdAt = new Date(element.CreatedAt);
    const formattedDate = createdAt.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });


    post.innerHTML = `
      <article>
        <header>
          <h1><i class="fa fa-user" aria-hidden="true"></i> ${element.UserName}</h1>
          <time>${formattedDate}</time>
        </header>
        <main>
          <section class="post-content">
            <h2>${element.Title}</h2>
            <p>${element.Content}</p>
          </section>
        </main>
        <footer>
          <nav>
            <div class="reaction-container"></div>
            <button class="comment-button" aria-label="View Comments">Comments</button>
          </nav>
        </footer>
      </article> 
    `;


        console.log(element)
    addReactionButtons("post", post, element.PostId)
    let display_comment = false
    post.querySelector('.comment-button').addEventListener('click', async (e) => {
      if (!display_comment) {
        const comment = document.createElement('div');
        comment.classList.add('comments-section');
        comment.innerHTML = `
                <h3>Comments</h3>
                <div class="comments-list">
                </div>
                <textarea placeholder="Add a comment..." rows="4" class="comment-input"></textarea>
                <button class="comment-submit">Submit</button>
                `
        post.appendChild(comment)
        await createComment(comment, comment.querySelector('.comments-list'), element.PostId)
        await getComment(comment.querySelector('.comments-list'), element.PostId)
        display_comment = true
      } else {
        post.querySelector('.comments-section').remove()
        display_comment = false
      }
    })
    container.append(post);
  });
}

const createComment = async (post, comment_part, post_id) => {
  const comment = post.querySelector('.comment-input')
  post.querySelector('.comment-submit').addEventListener('click', async (e) => {
    try {
      if (comment.value) {
        const res = await fetch(`http://localhost:8080/comments?post=${post_id}&comment=${comment.value}`, { method: 'POST', headers: { "Content-Type": 'application/json' } })
        const respons = await res.json()
        if (res.status === 401) {
          alert("you are unautherized")
        } else if (res.ok) {
          const com = document.createElement('div');
          com.classList.add('comment');
          com.innerHTML = `
          <strong>${respons.user_name}:</strong>
          ${comment.value}
          <div class="reaction-container"></div>

`;
          addReactionButtons("comment", com, respons.comment_id)
          comment_part.insertAdjacentElement('beforeend', com);
        }
        comment.value = ''

      }
    } catch (error) {
      console.error(error);
    }
  })
}


async function getCurrentReaction(targetType, targetId) {
  const params = {
    "user": "true",
    "target": targetType,
    "target_id": targetId,
  };
  const queryString = new URLSearchParams(params).toString();
  try {
    const response = await fetch(`http://localhost:8080/react?${queryString}`, {
      method: "GET",
      headers: {
        'Content-Type': 'application/json'
      }
    });
    if (!response.ok) throw new Error("Network response was not ok");
    return await response.json();
  } catch (err) {
    console.error(err);
    return null;
  }
}

function setReactionImage(reaction, selectedReactionImage, selectedReactionText) {
  const src = reaction ? `/assets/icons/${reaction.reaction_id}.png` : '/assets/icons/noReaction.png';
  const name = reaction ? reaction.name : 'Like';
  selectedReactionImage.src = src;
  selectedReactionText.textContent = name;
}

async function removeReaction(targetType, targetId, selectedReactionImage, selectedReactionText) {
  const params = { "target": targetType, "target_id": targetId };
  const queryString = new URLSearchParams(params).toString();
  const response = await fetch(`http://localhost:8080/react?${queryString}`, {
    method: "DELETE",
    headers: { 'Content-Type': 'application/json' }
  });
  if (!response.ok) {
    console.log("Network response was not ok");
    return;
  }
  selectedReactionImage.src = '/assets/icons/noReaction.png';
  selectedReactionText.textContent = 'Like';
}

async function addReaction(type, targetType, targetId, selectedReactionImage, selectedReactionText) {
  const params = { "type": type, "target": targetType, "target_id": targetId };
  const queryString = new URLSearchParams(params).toString();
  const response = await fetch(`http://localhost:8080/react?${queryString}`, {
    method: "PUT",
    headers: { 'Content-Type': 'application/json' }
  });
  if (!response.ok) {
    console.log("Network response was not ok");
    return;
  }
  selectedReactionImage.src = `/assets/icons/${type}.png`;
  selectedReactionText.textContent = type;
}

function setupReactionEvents(reactionButton, reactions, selectedReactionImage, selectedReactionText, targetType, targetId, likeButtonClicked) {
  reactionButton.addEventListener('mouseenter', () => reactions.style.display = 'flex');
  reactionButton.addEventListener('mouseleave', () => {
    setTimeout(() => {
      if (!reactions.matches(':hover')) reactions.style.display = 'none';
    }, 100);
  });

  reactions.addEventListener('mouseleave', () => reactions.style.display = 'none');

  reactionButton.addEventListener('click', async () => {
    if (likeButtonClicked) {
      await removeReaction(targetType, targetId, selectedReactionImage, selectedReactionText);
    } else {
      await addReaction("like", targetType, targetId, selectedReactionImage, selectedReactionText);
    }
  });

  reactions.addEventListener('click', async (event) => {
    const reactionOption = event.target.closest('.reaction-option');
    if (reactionOption) {
      const reactionId = reactionOption.dataset.text;
      await addReaction(reactionId, targetType, targetId, selectedReactionImage, selectedReactionText);
      reactions.style.display = 'none';
    }
  });
}

export async function addReactionButtons(targetType, target, targetId) {
  const reactionContainer = target.querySelector('.reaction-container');
  const currentPostReaction = await getCurrentReaction(targetType, targetId);
  if (!currentPostReaction) return;

  const likeButtonClicked = !!currentPostReaction.reaction_id;
  const src = likeButtonClicked ? `/assets/icons/${currentPostReaction.reaction_id}.png` : "/assets/icons/noReaction.png";
  const name = currentPostReaction.name || "Like";

  reactionContainer.innerHTML = `
    <button class="reaction-button" id="reactionButton">
      <img id="selectedReactionImage" src=${src} alt="${name}">
      <span id="selectedReactionText">${name}</span>
    </button>
    <div class="reactions" id="reactions">
      <div class="reaction-option" data-text="like">
        <img src="/assets/icons/like.png" alt="Like">
        <span>Like</span>
      </div>
      <div class="reaction-option" data-text="dislike">
        <img src="/assets/icons/dislike.png" alt="Dislike">
        <span>Dislike</span>
      </div>
    </div>`;

  const reactionButton = target.querySelector('#reactionButton');
  const reactions = target.querySelector('#reactions');
  const selectedReactionImage = target.querySelector('#selectedReactionImage');
  const selectedReactionText = target.querySelector('#selectedReactionText');

  setupReactionEvents(reactionButton, reactions, selectedReactionImage, selectedReactionText, targetType, targetId, likeButtonClicked);
}
