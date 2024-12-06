import { RenderPost} from "./rendring.js"

export const GetData = async (postIds) => {
    let target = []
    try{
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



export const logoutEvent = (log) => {
    log.addEventListener('click', async () => {
        try {
            const response = await fetch('http://localhost:8080/logout', {
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

export function showRegistrationModal() {
    const dialog = document.createElement('dialog');
    // Create message
    const message = document.createElement('p');
    message.textContent = 'You need to be logged in to react. Please register or log in to continue.';

    // Create register button
    const registerButton = document.createElement('button');
    registerButton.textContent = 'Register Now';

    // Create login button
    const loginButton = document.createElement('button');
    loginButton.textContent = 'Login';


    // Add event listeners
    registerButton.addEventListener('click', () => {
        window.location.href = '/register'; // Replace with your registration URL
    });
    loginButton.addEventListener('click', () => {
        window.location.href = '/login'; // Replace with your login URL
    });

    // Close dialog when clicking outside
    dialog.addEventListener('click', (event) => {
        if (event.target === dialog) {
            dialog.close();
        }
    }); 

    // Append content to dialog
    dialog.appendChild(message);
    dialog.appendChild(registerButton);
    dialog.appendChild(loginButton);

    // Append dialog to the body
    document.body.appendChild(dialog);

    // Show the dialog
    dialog.showModal();
}