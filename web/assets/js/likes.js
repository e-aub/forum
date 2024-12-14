import { showRegistrationModal } from "./script.js";

export function reactToggle(element /*html element*/, Id /*post or comment id*/, targetType /*post or comment*/) {
    const likeButton = element.querySelector('.post-footer button:nth-child(1)'); // Select the first button as Like
    const dislikeButton = element.querySelector('.post-footer button:nth-child(2)'); // Select the second button as Dislike

    likeButton.addEventListener('click', () => handleReact(likeButton, dislikeButton, Id, "like", targetType));
    dislikeButton.addEventListener('click', () => handleReact(dislikeButton, likeButton, Id, "dislike", targetType));
}

// Function to handle user interaction
export async function handleReact(button /*ex: like*/, follow /*ex: dislike*/, id, reactionType, targetType) {
    // Determine the HTTP method
    let method = button.getAttribute("data-clicked") === "true" ? "DELETE" : "PUT";

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

// Used to get reaction info or set them when the method is PUT or DELETE
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

// For the changes of reaction to take place on the front end instantly
function interactiveLike(button, follow) {
    const add = button.querySelector(".count");
    const subtract = follow.querySelector(".count");

    let count = parseInt(add.textContent, 10) || 0;
    let disCount = parseInt(subtract.textContent, 10) || 0;

    const buttonClicked = button.getAttribute("data-clicked") === "true";
    const followClicked = follow.getAttribute("data-clicked") === "true";

    if (!buttonClicked) {
        // If the main button wasn't clicked, like/dislike it
        add.textContent = count + 1;
        button.setAttribute("data-clicked", "true");
        button.classList.add("clicked");

        // If the other button was clicked, undo its action
        if (followClicked) {
            subtract.textContent = disCount - 1;
            follow.setAttribute("data-clicked", "false");
            follow.classList.remove("clicked");
        }
    } else {
        // If the main button was clicked, undo its action
        add.textContent = count - 1;
        button.setAttribute("data-clicked", "false");
        button.classList.remove("clicked");
    }
}

