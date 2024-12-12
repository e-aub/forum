import { renderPosts } from "./posts.js";

// Fetch data and render posts
export const GetData = async (postIds = false) => {
    let target = [];
    try {
        if (postIds === false) {
            postIds = [];
            let response = await fetch('http://localhost:8080/posts');
            if (!response.ok) throw new Error("Network response was not ok");
            let lastPostId = await response.json();
            for (let postId = 1; postId <= lastPostId; postId++) {
                postIds.push(postId);
            }
        }
        for (let i = postIds.length - 1; i >= 0; i--) {
            let link = `http://localhost:8080/posts?post_id=${postIds[i]}`;
            let postResponse = await fetch(link);
            if (postResponse.ok) {
                let post = await postResponse.json();
                if (post.PostId !== 0) {
                    target.push(post);
                }
            } else {
                console.log("error");
            }
        }
        renderPosts(target);
    } catch (err) {
        throw err
    }
    console.log(target);
};
GetData()

// Logout event
export const logoutEvent = (log) => {
    log.addEventListener('click', async () => {
        try {
            const response = await fetch('http://localhost:8080/logout', {
                method: 'POST',
                credentials: 'include'
            });

            if (response.ok) {
                console.log('Logged out successfully');
                window.location.href = '/'; // Redirect after successful logout
            } else {
                console.error('Logout failed');
            }
        } catch (error) {
            console.error('Error logging out:', error);
        }
    });
};

// Integrate logout event to the button
const logoutButton = document.getElementById('logoutBtn');
if (logoutButton) {
    logoutEvent(logoutButton);
}


export function showRegistrationModal() {
    const dialog = document.createElement('dialog');
    const message = document.createElement('p');
    message.textContent = 'You need to be logged in to react. Please register or log in to continue.';
    const registerButton = document.createElement('button');
    registerButton.textContent = 'Register Now';
    const loginButton = document.createElement('button');
    loginButton.textContent = 'Log In';
    dialog.appendChild(message);
    dialog.appendChild(registerButton);
    dialog.appendChild(loginButton);
    document.body.appendChild(dialog);
    dialog.showModal();
}