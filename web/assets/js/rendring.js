import { getComment } from "./script.js";
import { makePost, getReactInfo } from "./likes.js";
import { commentToggle } from "./comments.js";

export function RenderPost(posts) {
  const container = document.querySelector(".posts-container");
  container.innerHTML = "";

  posts.forEach((element) => {
    const post = document.createElement('div');
    post.classList.add('post');
    const reactInfo = getReactInfo({
      "is_own_react": "true",
      "target": "post",
      "target_id": element.PostId,
    }, "GET",).then(reactInfo => {
      post.innerHTML =  makePost(element, reactInfo)
    }).then(()=>{
      let display_comment = false;
      commentToggle(post, element, display_comment);
    }).then(()=> {
      container.append(post);
    });
////////////////////////////////////////

  });
}

export async function addReactionButtons(targetType, target, targetId) {
  const reactionContainer = target.querySelector('.reaction-container');
  let params = {
    "user": "true",
    "target": targetType,
    "target_id": targetId,
  };
  let queryString = new URLSearchParams(params).toString();

  try {
    const response = await fetch(`http://localhost:8080/react?${queryString}`, {
      method: "GET",
      headers: {
        'Content-Type': 'application/json'
      }
    });
    var currentPostReaction = await response.json();
    if (!response.ok) {
      throw new Error("Network response was not ok");
    }
    var likeButtonClicked = !!currentPostReaction.reaction_id; 
  } catch (err) {
    console.log(err);
    return;
  }
  
  const src = likeButtonClicked
    ? "/assets/icons/" + currentPostReaction.reaction_id + ".png"
    : "/assets/icons/noReaction.png";
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
  <div class="reaction-option" data-text="love">
    <img src="/assets/icons/love.png" alt="Love">
    <span>Love</span>
  </div>
  <div class="reaction-option" data-text="haha">
    <img src="/assets/icons/haha.png" alt="Haha">
    <span>Haha</span>
  </div>
  <div class="reaction-option" data-text="wow">
    <img src="/assets/icons/wow.png" alt="Wow">
    <span>Wow</span>
  </div>
  <div class="reaction-option" data-text="sad">
    <img src="/assets/icons/sad.png" alt="Sad">
    <span>Sad</span>
  </div>
  <div class="reaction-option" data-text="angry">
    <img src="/assets/icons/angry.png" alt="Angry">
    <span>Angry</span>
  </div>
</div>`;

  const reactionButton = target.querySelector('#reactionButton');
  const reactions = target.querySelector('#reactions');
  const selectedReactionImage = target.querySelector('#selectedReactionImage');
  const selectedReactionText = target.querySelector('#selectedReactionText');

  //show reactions on hover
  reactionButton.addEventListener('mouseenter', () => {
    reactions.style.display = 'flex';
  });

  //hide reaction whem the mouse leace raction container
  reactionButton.addEventListener('mouseleave', () => {
    setTimeout(() => {
      if (!reactions.matches(':hover')) {
        reactions.style.display = 'none';
      }
    }, 100);
  });

  reactionButton.addEventListener('click', async () => {
    if (likeButtonClicked) {
      //remove raction
      let params = {
        "target": targetType,
        "target_id": targetId,
      };
      let queryString = new URLSearchParams(params).toString();
      const response = await fetch(`http://localhost:8080/react?${queryString}`, {
        method: "DELETE",
        headers: {
          'Content-Type': 'application/json'
        }
      });
      if (!response.ok) {
        console.log("Network response was not ok");
        return;
      }
      selectedReactionImage.src = '/assets/icons/noReaction.png';
      selectedReactionText.textContent = 'Like';
      likeButtonClicked = !likeButtonClicked;
      
    } else {
      // alert("not clicked")
      // add like raacion
      let params = {
        "type": "like",
        "target": targetType,
        "target_id": targetId,
      };
      let queryString = new URLSearchParams(params).toString();
      const response = await fetch(`http://localhost:8080/react?${queryString}`, {
        method: "PUT",
        headers: {
          'Content-Type': 'application/json'
        }
      })
      
      if (!response.ok) {
        console.log("Network response was not ok");
        return;
      }
      selectedReactionImage.src = '/assets/icons/like.png';
      selectedReactionText.textContent = 'Like';
      likeButtonClicked = !likeButtonClicked;
    }
  });

  reactions.addEventListener('mouseleave', () => {
    reactions.style.display = 'none';
  });

  // Handle reaction selection
  reactions.addEventListener('click', async (event) => {
    const reactionOption = event.target.closest('.reaction-option');
    if (reactionOption) {
      const reactionImage = reactionOption.querySelector('img').src;
      const reactionId = reactionOption.dataset.text;


      let params = {
        "type": reactionId,
        "target": targetType,
        "target_id": targetId,
      };
      let stringParams = new URLSearchParams(params).toString();
      try {
        const response = await fetch(`http://localhost:8080/react?${stringParams}`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json'
          },
        });
        if (!response.ok) {
          throw new Error("error while adding reaction");
        }
        likeButtonClicked = true;
        selectedReactionImage.src = reactionImage;
        selectedReactionText.textContent = reactionId;
        reactions.style.display = 'none';
      } catch (error) {
        console.error(error);
      }
    }
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
