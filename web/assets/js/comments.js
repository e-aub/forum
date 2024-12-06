import { handleReact } from "./likes.js";

export async function commentToggle(post, element, display_comment){
    post.querySelector('.comment-button').addEventListener('click', async (e) => {
        if (!display_comment) {
            const comment = document.createElement('div');
            comment.classList.add('comments-section');
            comment.innerHTML = `
            <div class="comments-list">
            </div>
            <div class="add_comment" >
                <textarea placeholder="Add a comment..."  class="comment-input"></textarea>
                <button class="comment-submit">Submit</button>
            <div>
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
};

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
              <div class="likes">
                  <button data-clicked="false" id="likeButton" class="com_like" style="background-color: white;">
                      Like <span class="count">${true}</span>
                  </button>
                  <button data-clicked="false" id="dislikeButton" class="com_dislike" style="background-color: white;">
                      Dislike <span class="count">${true}</span>
                  </button>
              </div>
          `
          comment_part.insertAdjacentElement('beforeend', com);
        }
        comment.value = ''
      }
    } catch (error) {
      console.error(error);
    }
  })
}


export const getComment = async (post, id) => {
    try {
        const res = await fetch(`http://localhost:8080/comments?post=${id}`)
        if (res.ok) {
            const allComment = await res.json()
            if (allComment) {
                for (let comment of allComment) {
                    const com = document.createElement('div');
                    com.classList.add('comment');
                    com.innerHTML = `
                        <div class="one_comment">
                            <p><i class="fa fa-user"></i> ${comment.user_name}:<i> ${comment.content}</i> </p> 
                            <div class="actions">
                                <button data-clicked="${comment.clicked}"  class="com_like" 
                                style="background-color: ${comment.clicked ? '#15F5BA' : 'white'};">
                                    <i class="fas fa-thumbs-up"></i> <span class="count">${comment.like_count}</span>
                                </button>

                                <button data-clicked="${comment.disclicked}"  class="com_dislike" 
                                style="background-color: ${comment.disclicked ? '#15F5BA' : 'white'};">
                                    <i class="fas fa-thumbs-down"></i> <span class="count">${comment.dislike_count}</span>
                                </button>
                            </div>
                        <div>
                    `;

                    post.insertAdjacentElement('beforeend', com);

                    // Add event listeners for like and dislike buttons
                    const likeButton = com.querySelector('.com_like');
                    const dislikeButton = com.querySelector('.com_dislike');

                    likeButton.addEventListener('click', async () => {
                        await handleReact(likeButton, dislikeButton,comment.comment_id, 'like', "comment");
                    });

                    dislikeButton.addEventListener('click', async () => {
                        await handleReact(dislikeButton, likeButton, comment.comment_id, 'dislike', "comment");
                    });
                }
            }
        }
    } catch (error) {
        console.error(error);
    }
}
