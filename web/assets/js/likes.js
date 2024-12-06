import {showRegistrationModal} from  "./script.js"

export function reactToggle(post , Id /*post or comment id*/){
    const likeButton = post.querySelector('.like');
    const dislikeButton = post.querySelector('.dislike');

    likeButton.addEventListener('click', () => handleReact(likeButton,dislikeButton, Id, "like", "post"));
    dislikeButton.addEventListener('click', () => handleReact(dislikeButton,likeButton, Id, "dislike", "post"));
}

export function reactToggleCom(comment , Id /*post or comment id*/){
        // Add event listeners for like and dislike buttons
    const likeButton = com.querySelector('.com_like');
    const dislikeButton = com.querySelector('.com_dislike');

    likeButton.addEventListener('click', async () => {
        await handleReact(likeButton, dislikeButton, respons.comment_id, 'like', "comment");
    });

    dislikeButton.addEventListener('click', async () => {
        await handleReact(dislikeButton, likeButton, respons.comment_id, 'dislike', "comment");
    });
}
// Function to handle user interaction
export async function handleReact(button, follow, id, reactionType, targetType) {
    // the method here can be eather put or delete

    let method = button.getAttribute("data-clicked") === "true" ? "DELETE" : "PUT";

    console.log(button.getAttribute("data-clicked"))
    console.log(method)
    try {
        const result = await getReactInfo(
            {
                reaction_type: reactionType,
                target_type: targetType,
                target_id: id,
            },
            method
        );

        if (!result.success) { 
            showRegistrationModal(); 
        } else {
            interactiveLike(button, follow); 
        }
    } catch (error) {
        console.error("Error in handleReact:", error);
    }
}

export async function getReactInfo(params, method) {
    const queryString = new URLSearchParams(params).toString();
    const url = `http://localhost:8080/react?${queryString}`;

    try {
        const response = await fetch(url, {
            method: method,
            headers: {
                'Content-Type': 'application/json',
            },
        });

        if (!response.ok) {
            // Handle non-OK responses
            const errorText = await response.text(); // Use text() for error body
            console.error("API error:", errorText);
            return { success: false, error: errorText || "Unknown error" };
        }

        // If response has no body, return success with no data
        const contentLength = response.headers.get("Content-Length");
        if (!contentLength || parseInt(contentLength) === 0) {
            return { success: true, data: null };
        }

        // Parse JSON response
        return { success: true, data: await response.json() };

    } catch (err) {
        console.error("Fetch error:", err);
        return { success: false, error: err.message };
    }
}




function interactiveLike(button , follow ){
    const add = button.querySelector(".count");
    const subtract = follow.querySelector(".count");

    // Parse the current count from the button's span text
    let count = parseInt(add.textContent, 10);
    let disCount = parseInt(subtract.textContent, 10)

    if (button.getAttribute("data-clicked") === "false") {

        count += 1; add.textContent = count; // Update the displayed count
        button.setAttribute("data-clicked", "true");
        button.style.backgroundColor = '#15F5BA'
        follow.style.backgroundColor = 'white'

        if (follow.getAttribute("data-clicked") === "true") {
            disCount -= 1; subtract.textContent = disCount; // Update the displayed count
            follow.setAttribute("data-clicked", "false");
            follow.style.backgroundColor = 'white'
        }
    }else if (button.getAttribute("data-clicked") === "true") {
        count -= 1; add.textContent = count; // Update the displayed count
        button.setAttribute("data-clicked", "false");
        button.style.backgroundColor = 'white'
    }
}

