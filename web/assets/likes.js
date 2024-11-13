

export function addLikeDislikeListeners(post, postId) {
    const likeButton = post.querySelector('.like');
    const dislikeButton = post.querySelector('.dislike');

    likeButton.addEventListener('click', () => handleReact(postId, "like"));
    dislikeButton.addEventListener('click', () => handleReact(postId, "dislike"));
}

async function handleReact(postId, type ) {
    // Logic to handle the "like" action
    console.log(`Liked/disliked post with ID: ${postId}`);
    // Update like count, send API request, etc.
       try {
        // Send API request
        const response = await fetch(`/api/react/${postId}/${type}`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
        });

        if (!response.ok) {
           showRegistrationModal(); 
        }
    } catch (error) {
        console.error("Error:", error);
        // Revert the UI update if the request fails
    }
}


function showRegistrationModal() {
    // Create modal container
    const modal = document.createElement('div');
    modal.style.position = 'fixed';
    modal.style.top = '0';
    modal.style.left = '0';
    modal.style.width = '100vw';
    modal.style.height = '100vh';
    modal.style.backgroundColor = 'rgba(0, 0, 0, 0.5)';
    modal.style.display = 'flex';
    modal.style.justifyContent = 'center';
    modal.style.alignItems = 'center';
    modal.style.zIndex = '1000'; // Ensure it appears on top

    // Create modal content
    const modalContent = document.createElement('div');
    modalContent.style.backgroundColor = 'black';
    modalContent.style.padding = '20px';
    modalContent.style.borderRadius = '8px';
    modalContent.style.textAlign = 'center';
    modalContent.style.maxWidth = '400px';
    modalContent.style.boxShadow = '0 0 10px rgba(0, 0, 0, 0.1)';
    
    // Create message
    const message = document.createElement('p');
    message.textContent = 'You need to be loged int to to react. Please register or loging to continue.';

    // Create register button
    const registerButton = document.createElement('button');
    registerButton.textContent = 'Register Now';
    registerButton.style.marginTop = '10px';
    registerButton.style.padding = '10px 20px';
    registerButton.style.fontSize = '16px';

    // Create Login button
    const loginButton = document.createElement('button');
    loginButton.textContent = 'Login';
    loginButton.style.marginTop = '10px';
    loginButton.style.padding = '10px 20px';
    loginButton.style.fontSize = '16px';

    // Add event listener to the register button
    registerButton.addEventListener('click', () => {
        // Redirect to the registration page
        window.location.href = '/register'; // Replace with your registration URL
    });
        loginButton.addEventListener('click', () => {
        // Redirect to the login page
        window.location.href = '/login'; // Replace with your login URL
    });

    // Append content to modal
    modalContent.appendChild(message);
    modalContent.appendChild(registerButton);
    modalContent.appendChild(loginButton);

    modal.appendChild(modalContent);

    // Append modal to the body
    document.body.appendChild(modal);
}