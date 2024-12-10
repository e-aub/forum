document.getElementById("login-form").addEventListener("submit", async function (event) {
    event.preventDefault(); // Prevent default form submission

    // Capture form data
    const username = document.getElementById("login-username").value.trim();
    const password = document.getElementById("login-password").value;

    // Basic validation
    if (!username || password.length < 6) {
        alert("Invalid username or password. Password must be at least 6 characters long.");
        return;
    }

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
            const error = await response.text();
            alert(`Error: ${error}`);
        }
    } catch {
        alert("An error occurred during login.");
    }
});
