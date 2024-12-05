

export function commentToggle(post, element, display_comment){
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