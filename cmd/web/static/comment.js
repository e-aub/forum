export const getComment = async () => {
    try {
        const res = await fetch('http://localhost:8080/api/comments')
        if (res.ok) {
            const allComment = await res.json()
            for (let comment of allComment.data) {
                const container = document.querySelector('.comments')
                const newDiv = createComment({
                    id: comment.comment_id,
                    user: comment.user_id,
                    post: comment.post_id,
                    comment: comment.content
                })
                container.insertAdjacentElement('beforeend', newDiv)
            }
        }
    } catch (error) {
        console.error(error);
    }
}

export const addComment = (comment) => {
    const button = document.querySelector('.add-button')
    button.addEventListener('click', async () => {
        const obj = {
            post: '1',
            user: '1',
            comment: comment.value,
        }
        const data = new URLSearchParams(obj)

        if (obj.comment) {
            try {
                const res = await fetch(`http://localhost:8080/api/comments?${data}`, { method: 'POST' })
                if (res.ok) {
                    const container = document.querySelector('.comments')
                    const newDiv = createComment(obj)
                    container.insertAdjacentElement('beforeend', newDiv)
                }
            } catch (error) {
                console.log(error)
            }
        }
    })
}

export const deleteComment = () => {
    try {
        const allCom = document.querySelectorAll('.comment')
        allCom.forEach(ele => {
            const button = ele.querySelector('.delete-button')
            const id = button.getAttribute('id')
            button.addEventListener('click', async (e) => {
                console.log(id);
                const res = await fetch(`http://localhost:8080/api/comments?comment=${id}&user=${1}`)
                if (res.ok) {
                    ele.remove()
                }
            })
        })
    } catch (error) {
        console.error(error);
    }
}

const createComment = (data) => {
    const div = document.createElement('div');
    div.className = 'comment';
    div.innerHTML = `
        <div class="comment-content">
            <div>${data.post}</div>
            <div class="meta">User ID: ${data.user}</div>
            <p>${data.comment}</p>
        </div>
        <button id="${data.id}" class="delete-button">Delete</button>
    `;
    return div;
}