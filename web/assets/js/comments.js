export async function commentToggle(post, element, display_comment){
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
};


const getComment = async (post, id) => {
    try {
        const res = await fetch(`http://localhost:8080/comments?post=${id}`)
        if (res.ok) {
            const allComment = await res.json()
            if (allComment) {
                for (let comment of allComment) {
                    const com = document.createElement('div');
                    com.classList.add('comment');
                    com.innerHTML = `
                    <strong>${comment.user_name}:</strong>
                    <strong>${comment.content}:</strong>
                    <div class="reaction-container"></div>
                    `;
                    console.log(comment)
                    addReactionButtons("comment", com, comment.comment_id)

                    post.insertAdjacentElement('beforeend', com);
              
                }
            }
        }
    } catch (error) {
        console.error(error);
    }
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
          <div class="reaction-container"></div>`;
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
