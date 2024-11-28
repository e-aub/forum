export function addLikeDislikeListeners(post, postId) {
    const likeButton = post.querySelector('.like');
    const dislikeButton = post.querySelector('.dislike');

    likeButton.addEventListener('click', () => handleReact(likeButton,dislikeButton, postId, "like", "post"));
    dislikeButton.addEventListener('click', () => handleReact(dislikeButton,likeButton, postId, "dislike", "post"));
}

export async function handleReact(button, follow , postId, type , target_Type) {
    // Logic to handle the "like" action
    console.log(`Liked/disliked post/comment with ID: ${postId}`);
    // Update like count, send API request, etc.
       try {
        // Send API request
        const response = await fetch(`/api/react/${postId}/${type}/${target_Type}`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
        });

        if (!response.ok) {
           showRegistrationModal(); 
        }else{
            interactiveLike(button, follow)
        }
    } catch (error) {
        console.error("Error:", error);
        // Revert the UI update if the request fails
    }
}

function interactiveLike(button , follow ){
    // Access the count within the specific elements
    // Check if the button was already clicked
    const add = button.querySelector(".count");
    const subtract = follow.querySelector(".count");

    // Parse the current count from the button's span text
    let count = parseInt(add.textContent, 10);

    if (button.getAttribute("data-clicked") === "false") {
        count += 1;
        add.textContent = count; // Update the displayed count
        button.setAttribute("data-clicked", "true");
        button.style.backgroundColor = '#15F5BA'
        follow.style.backgroundColor = 'white'
        if (follow.getAttribute("data-clicked") === "true") {
            count -= 1;
            subtract.textContent = count; // Update the displayed count
            follow.setAttribute("data-clicked", "false");
            //button.style.backgroundColor = '#15F5BA'
            follow.style.backgroundColor = 'white'
        }
    }else if (button.getAttribute("data-clicked") === "true") {
        count -= 1;
        add.textContent = count; // Update the displayed count
        button.setAttribute("data-clicked", "false");
        button.style.backgroundColor = 'white'
    }
    // Increment the count
}

function showRegistrationModal() {
    // Create dialog element
    const dialog = document.createElement('dialog');


    // Create message
    const message = document.createElement('p');
    message.textContent = 'You need to be logged in to react. Please register or log in to continue.';

    // Create register button
    const registerButton = document.createElement('button');
    registerButton.textContent = 'Register Now';
    registerButton.style.marginTop = '10px';
    registerButton.style.padding = '10px 20px';
    registerButton.style.fontSize = '16px';

    // Create login button
    const loginButton = document.createElement('button');
    loginButton.textContent = 'Login';
    loginButton.style.marginTop = '10px';
    loginButton.style.padding = '10px 20px';
    loginButton.style.fontSize = '16px';

    // Create stay button
    const stayButton = document.createElement('button');
    stayButton.textContent = 'Stay';
    stayButton.style.marginTop = '10px';
    stayButton.style.padding = '10px 20px';
    stayButton.style.fontSize = '16px';

    // Add event listeners
    registerButton.addEventListener('click', () => {
        window.location.href = '/register'; // Replace with your registration URL
    });
    loginButton.addEventListener('click', () => {
        window.location.href = '/login'; // Replace with your login URL
    });
    stayButton.addEventListener('click', () => {
        dialog.close(); // Close the dialog
    });

    // Append content to dialog
    dialog.appendChild(message);
    dialog.appendChild(registerButton);
    dialog.appendChild(loginButton);
    dialog.appendChild(stayButton);

    // Append dialog to the body
    document.body.appendChild(dialog);

    // Show the dialog
    dialog.showModal();
}