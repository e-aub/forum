import { RenderPost } from "./rendring.js"

export const GetData = async (postIds) => {
    let target = []
    try{
        if (postIds === false) {
            postIds = [];
            let response = await fetch('http://localhost:8080/api/posts');
            if (!response.ok) throw new Error("Network response was not ok");
            let lastPostId = await response.json();
            for (let postId = 1; postId <= lastPostId; postId++) {
                postIds.push(postId);
            }
        }
        for (let i = postIds.length - 1; i >= 0; i--) {
            let link = `http://localhost:8080/api/posts?id=${postIds[i]}`;
            let postResponse = await fetch(link);
            if (!postResponse.ok) throw new Error("Failed to fetch post data");
            let post = await postResponse.json();
            if (post.PostId !== 0) {
                target.push(post)
                RenderPost(target)
            }
        }
    } catch (err) {
        console.error("Error fetching data:", err);
    }
};

export const getComment = async (post, id) => {
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
                    <div class="likes">
                        <button class="like">Like</button>
                        <button class="dislike">Dislike</button>
                    </div>
                    `
                    post.insertAdjacentElement('beforeend', com)
                }
            }
        }
    } catch (error) {
        console.error(error);
    }
}

export const logoutEvent = (log) => {
    log.addEventListener('click', async () => {
        try {
            const response = await fetch('http://localhost:8080/api/logout', {
                method: 'POST',
                credentials: 'include'
            });

            if (response.ok) {
                console.log('Logged out successfully');
                window.location.href = '/';
            } else {
                console.error('Logout failed');
            }
        } catch (error) {
            console.error('Error logging out:', error);
        }
    });
}