import { getComment } from "./script.js";
import {addLikeDislikeListeners} from "./likes.js"
import { handleReact } from "./likes.js";

export function RenderPost(args) {
    const container = document.querySelector(".container");
    container.innerHTML = "";

    args.forEach((element, index) => {
        const post = document.createElement('div');
        post.classList.add('post');

        post.innerHTML = `
        <div class="post-header">
            <span class="post-index"> ${element.Title}</span>
        </div>
        <div class="post-content">
            <p><strong>User name:</strong> ${element.UserName}</p>
            <p><strong>Content:</strong> ${element.Content}</p>
            <p><strong>Time:</strong> ${element.Created_At}</p>
            <p><strong>Category:</strong> ${element.Categories.join(', ')}</p>
        </div>
        <button class="comment-button">Comments</button>
     <div class="likes">
    <button 
        data-clicked="${element.Clicked}" 
        class="like" 
        style="background-color: ${element.Clicked ? '#15F5BA' : 'white'};">
        Like <span class="count">${element.LikeCount}</span>
    </button>
    <button 
        data-clicked="${element.DisClicked}" 
        class="dislike" 
        style="background-color: ${element.DisClicked ? '#15F5BA' : 'white'};">
        Dislike <span class="count">${element.DislikeCount}</span>
    </button>
</div>
        </div>
        `;

        addLikeDislikeListeners(post, element.PostId);

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
                const res = await fetch(`http://localhost:8080/api/comments?post=${post_id}&comment=${comment.value}`, { method: 'POST' })
                const respons = await res.json()
                if (res.status === 401) {
                    alert(respons)
                } else if (res.ok) {

                    const com = document.createElement('div');
                    com.classList.add('comment');
com.innerHTML = `
    <strong>${respons.user_name}:</strong>
    ${comment.value}
    <div class="likes">
        <button data-clicked="false" id="likeButton" class="com_like" style="background-color: white;">
            Like <span class="count">${respons.like_count}</span>
        </button>
        <button data-clicked="false" id="dislikeButton" class="com_dislike" style="background-color: white;">
            Dislike <span class="count">${respons.dislike_count}</span>
        </button>
    </div>
`;
                    comment_part.insertAdjacentElement('beforeend', com)
                    // Add event listeners for like and dislike buttons
                    const likeButton = com.querySelector('.com_like');
                    const dislikeButton = com.querySelector('.com_dislike');

                    likeButton.addEventListener('click', async () => {
                        await handleReact(likeButton, dislikeButton, respons.comment_id, 'like', "comment");
                    });

                    dislikeButton.addEventListener('click', async () => {
                        await handleReact(dislikeButton, likeButton, respons.comment_id, 'dislike', "comment");
                    });

                }
                comment.value = ''
                
            }
        } catch (error) {
            console.error(error);
        }
    })
}
