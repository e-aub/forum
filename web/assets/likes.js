

export function addLikeDislikeListeners(post, postId) {
    const likeButton = post.querySelector('.like');
    const dislikeButton = post.querySelector('.dislike');

    likeButton.addEventListener('click', () => handleLike(postId));
    dislikeButton.addEventListener('click', () => handleDislike(postId));
}

async function handleLike(postId) {
    // Logic to handle the "like" action
    console.log(`Liked post with ID: ${postId}`);
    // Update like count, send API request, etc.
       try {
        // Send API request
        const response = await fetch(`/api/react/${postId}/like`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
        });

        if (!response.ok) {
            throw new Error("Failed to dislike the post");
        }
    } catch (error) {
        console.error("Error:", error);
        // Revert the UI update if the request fails
    }
}

function handleDislike(postId) {
    // Logic to handle the "dislike" action
    console.log(`Disliked post with ID: ${postId}`);
    // Update dislike count, send API request, etc.
}

