document.getElementById("login-form").addEventListener("submit", async function (event) {
    event.preventDefault();

    // Capture form data
    const username = document.getElementById("login-username").value;
    const password = document.getElementById("login-password").value;
    const messageElement = document.getElementById("responseMessage");

    // Send data to the API
    try {
        const response = await fetch("http://localhost:8080/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ username, password }),
            credentials: "include",
        });

        if (response.ok) {
            window.location.href = '/'; // Redirect on success
        } else {
            const errorData = await response.text();
            messageElement.textContent = `Error: ${errorData}`;
            messageElement.style.color = "red";
        }
    } catch {
        messageElement.textContent = "An error occurred during registration.";
        messageElement.style.color = "red";
    }
});

// // Logout event
// export const logoutEvent = (log) => {
//     log.addEventListener('click', async () => {
//         try {
//             const response = await fetch('http://localhost:8080/logout', {
//                 method: 'POST',
//                 credentials: 'include'
//             });

//             if (response.ok) {
//                 console.log('Logged out successfully');
//                 window.location.href = '/'; // Redirect after successful logout
//             } else {
//                 console.error('Logout failed');
//             }
//         } catch (error) {
//             console.error('Error logging out:', error);
//         }
//     });
// };

// Integrate logout event to the button
// const logoutButton = document.getElementById('logoutBtn');
// if (logoutButton) {
//     logoutEvent(logoutButton);
// }


