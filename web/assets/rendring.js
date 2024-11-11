import { getComment } from "./script.js";

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
        </div>
        <button class="comment-button">Comments</button>
        <div class="likes">
            <button class="like"><span class="number-like">${element.Likes.like}</span> ⬆</button>
            <button class="dislike"><span class="number-dislike">${element.Likes.dislike}</span> ⬇</button>
        </div>
        `;

        likesEvent(post.querySelector('.likes'), element.PostId)

        let display_comment = false; let hiden = false
        post.querySelector('.comment-button').addEventListener('click', async (e) => {
            if (hiden) {
                post.querySelector('.comments-section').style.display = 'block'
                hiden = false
            } else {
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
                    post.querySelector('.comments-section').style.display = 'none'
                    hiden = true;
                }
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
                    <div class="likes">
                        <button class="like"><span class="number-like">${comment.Likes.like}</span> ⬆</button>
                        <button class="dislike"><span class="number-dislike">${comment.Likes.dislike}</span> ⬇</button>
                    </div>
                    `
                    comment_part.insertAdjacentElement('beforeend', com)
                    likesEvent(com.querySelector('.likes'), post_id)
                }
                comment.value = ''
            }
        } catch (error) {
            console.error(error);
        }
    })
}


export const likesEvent = (parentClass, post_id) => {
    const like = parentClass.querySelector('.like')
    const dislike = parentClass.querySelector('.dislike')
    like.addEventListener('click', async e => {
        const res = await fetch(`http://localhost:8080/api/likes?postId=${post_id}&type=like`, { method: 'POST' })
        const data = await res.json()
        if (res.status === 401) {
            alert(res)
        } else if (res.ok) {
            like.querySelector('.number-like').textContent = data.like
            dislike.querySelector('.number-dislike').textContent = data.dislike
        }
    })
    dislike.addEventListener('click', async e => {
        const res = await fetch(`http://localhost:8080/api/likes?postId=${post_id}&type=dislike`, { method: 'POST' })
        if (res.status === 401) {
            alert(res)
        } else if (res.ok) {
            const data = await res.json()
            dislike.querySelector('.number-dislike').textContent = data.dislike
            like.querySelector('.number-dislike').textContent = data.dislike
        }
    })
}