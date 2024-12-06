import { handleReact, getReactInfo, reactToggle } from "./likes.js";


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
                  <button data-clicked="false" id="likeButton" class="like" id="com_like" style="background-color: white;">
                      Like <span class="count">${true}</span>
                  </button>
                  <button data-clicked="false" id="dislikeButton" class="dislike" id="com_dislike" style="background-color: white;">
                      Dislike <span class="count">${true}</span>
                  </button>
              </div>
          `
          reactToggle(com, comment.comment_id, "comment")
          comment_part.insertAdjacentElement('beforeend', com);
        }
        comment.value = ''
      }
    } catch (error) {
      console.error(error);
    }
  })
}


export const getComment = async (element , id) => {
    try {
        const res = await fetch(`http://localhost:8080/comments?post=${id}`)
        if (res.ok) {
            const allComment = await res.json()
            if (allComment) {
                for (let comment of allComment) {
                    const com = document.createElement('div');
                    com.classList.add('comment');

                    try{
                        // Fetch reaction info
                        const reactInfo = await getReactInfo({
                            target_type: "comment",
                            target_id: id,
                        }, "GET");


                        com.innerHTML = commentTemplate(comment, reactInfo.data)
                        reactToggle(com, comment.comment_id, "comment")
                        
                        // Add event listeners for like and dislike buttons
                        element.append(com)
                    }catch (error) {
                        console.error("Error rendering post:", error);
                    }
                }
            }
        }
    } catch (error) {
        console.error(error);
    }
}

function commentTemplate(comment,reactInfo){
    console.log(reactInfo)
    let liked = false ;
    let disliked = false ;

    let likeCount = reactInfo.liked_by? reactInfo.liked_by.length : 0 ;
    let disLikeCount = reactInfo.disliked_by? reactInfo.disliked_by.length : 0 ; 

    if (!!reactInfo.user_reaction){
      liked = reactInfo.user_reaction === "like"
      disliked = !liked
    }else{
      liked = false 
      disliked = false;
    }

    const innerHTML = `
    <div class="one_comment">
        <p><i class="fa fa-user"></i> ${comment.user_name}:<i> ${comment.content}</i> </p> 
        <div class="actions">
            <button data-clicked="${liked}" class="like" id="com_like" 
            style="background-color: ${liked ? '#15F5BA' : 'white'};">
                <i class="fas fa-thumbs-up"></i> <span class="count">${likeCount}</span>
            </button>

            <button data-clicked="${disliked}" class="dislike" id="com_dislike" 
            style="background-color: ${disliked ? '#15F5BA' : 'white'};">
                <i class="fas fa-thumbs-down"></i> <span class="count">${disLikeCount}</span>
            </button>
        </div>
    <div>
    `;
    return innerHTML
}