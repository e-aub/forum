export const GetData = async () => {
    let target = []
    try {
        let response = await fetch('http://localhost:8080/api/posts');
        if (!response.ok) throw new Error("Network response was not ok");

        let data = await response.json();
        console.log(data);

        if (data) {
            for (let i = data; i > 0; i--) {
                let link = `http://localhost:8080/api/posts?id=${i}`;
                let postResponse = await fetch(link);
                if (!postResponse.ok) throw new Error("Failed to fetch post data");
                let post = await postResponse.json();
                if (post.PostId !== 0) {
                    target.push(post)
                    RenderPost(target)
                }
            }
        }
    } catch (err) {
        console.error("Error fetching data:", err);
    }
};

function RenderPost(args) {
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
        </div>
        <button class="comment-button">Comments</button>
        `;

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

document.getElementById('logout-button').addEventListener('click', async () => {
    try {
        const response = await fetch('http://localhost:8080/api/logout', {
            method: 'POST',
            credentials: 'include'
        });

        if (response.ok) {
            // Handle successful logout
            console.log('Logged out successfully');
            window.location.href = '/login'; // Redirect to login page
        } else {
            console.error('Logout failed');
        }
    } catch (error) {
        console.error('Error logging out:', error);
    }
});

const getComment = async (post, id) => {
    try {
        const res = await fetch(`http://localhost:8080/api/comments?post=${id}`)
        if (res.ok) {
            const allComment = await res.json()
            if (allComment) {
                for (let comment of allComment) {
                    const com = document.createElement('div');
                    com.classList.add('comment');
                    com.innerHTML = `
                <strong>${comment.user_name}:</strong> ${comment.content}
                `
                    post.insertAdjacentElement('beforeend', com)
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
                const res = await fetch(`http://localhost:8080/api/comments?post=${post_id}&comment=${comment.value}`, { method: 'POST' })
                const respons = await res.json()
                if (res.status === 401) {
                    alert(respons)
                } else if (res.ok) {
                    const com = document.createElement('div');
                    com.classList.add('comment');
                    com.innerHTML = `
                    <strong>${respons.user_name}:</strong> ${comment.value}
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

