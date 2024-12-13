import { renderPosts } from "./posts.js";

// Fetch data and render posts
export const GetData = async (postIds = false) => {
    if (postIds == null) {
        return;
    }

    const postsContainer = document.querySelector(".posts");
    postsContainer.innerHTML = "";
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

        console.log(postIds);
        renderPage(postIds, postsContainer);
        const debouncedRenderPage = debounce(renderPage, 1000)
        window.addEventListener('scroll', () => {
            const scrollPosition = window.scrollY;
            const documentHeight = document.documentElement.scrollHeight;
            const windowHeight = window.innerHeight;
            if (scrollPosition + windowHeight >= documentHeight - 10) {
                debouncedRenderPage(postIds, postsContainer)
                // renderPage(postIds, postsContainer);

            }
        });

    } catch (err) {
        console.log(err);
    }
};

function debounce(func, delay) {
    let timer;
    return function (...args) {
        clearTimeout(timer);
        timer = setTimeout(() => func.apply(this, args), delay);
    };
}


async function renderPage(postIds, postsContainer) {
    let target = [];
    for (let i = 0; i < Math.min(10, postIds.length); i++) {
        let link = `http://localhost:8080/posts?post_id=${postIds.pop()}`;
        let postResponse = await fetch(link);

        if (postResponse.ok) {
            let post = await postResponse.json();
            target.push(post);
        } else {
            if (postResponse.status !== 404) {
                throw new Error("Response not ok");
            }
        }
    }

    console.log(target);
    if (target.length > 0) {
        await renderPosts(postsContainer, target);
    }
}

// GetData()


// Logout event
export const logoutEvent = (log) => {
    log.addEventListener('click', async () => {
        try {
            const response = await fetch('http://localhost:8080/logout', {
                method: 'POST',
                credentials: 'include'
            });
            console.log("test");
            

            if (response.ok) {
                console.log('Logged out successfully');
                // window.location.href = '/'; // Redirect after successful logout
            } else {
                console.error('Logout failed');
            }
        } catch (error) {
            console.error('Error logging out:', error);
        }
    });
};



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

export function SubmitForm(category, event) {
    event.preventDefault()
    const params = new URLSearchParams({ category: category });
    window.location.href = `/categories?${parfalseams}`;
}

