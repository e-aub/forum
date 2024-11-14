import { getComment } from "./script.js";
import {addLikeDislikeListeners} from "./likes.js"

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
            <p><strong>Category:</strong> ${element.Category}</p>
            <p><strong>likes:</strong> ${element.LikeCount}</p>
            <p><strong>dislikes:</strong> ${element.DislikeCount}</p>
        </div>
        <button class="comment-button">Comments</button>
        <div class="likes">
        <button class="like" style="background-color: white;">Like</button>
        <button class="dislike" style="background-color: white;">Dislike</button>
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
                    <strong>${respons.user_name}:</strong> ${comment.value}
                    <strong>likes: ${respons.like_count}<strong>
                    <strong>dislikes: ${respons.dislike_count}<strong>
                    <div class="likes">
                        <button class="like">Like</button>
                        <button class="dislike">Dislike</button>
                    </div>
                    `
                    comment_part.insertAdjacentElement('beforeend', com)
                }
                comment.value = ''
            }
        } catch (error) {
            console.error(error);
        }
    })
}
