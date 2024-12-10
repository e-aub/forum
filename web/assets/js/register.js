function validateSignup() {
    const username = document.getElementById("signup-username").value;
    const email = document.getElementById("signup-email").value;
    const password = document.getElementById("signup-password").value;
    const confirmPassword = document.getElementById("signup-confirm-password").value;

    if (username === "" || email === "" || password === "" || confirmPassword === "") {
        alert("Please fill out all fields.");
        return false;
    }

    const emailPattern = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$/;
    if (!emailPattern.test(email)) {
        alert("Please enter a valid email address.");
        return false;
    }

    if (password.length < 6) {
        alert("Password must be at least 6 characters long.");
        return false;
    }

    if (password !== confirmPassword) {
        alert("Passwords do not match.");
        return false;
    }

    return true;
}

// JavaScript to handle form validation and submission
document.getElementById("signup-form").addEventListener("submit", async function (event) {
    event.preventDefault();  // Prevent default form submission

    // First, validate form fields using validateSignup
    if (!validateSignup()) {
        console.log("invlaadasd")
        return;
    }

    // Capture form data
    const username = document.getElementById("signup-username").value.trim();
    const email = document.getElementById("signup-email").value.trim();
    const password = document.getElementById("signup-password").value;

    // Display message area
    const messageElement = document.getElementById("responseMessage");

    // Send data to the API if validation passes
    try {
        const response = await fetch("http://localhost:8080/register", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ username, email, password }),
            credentials: "include",
        });

        if (response.ok) {
            messageElement.textContent = "Registration successful!";
            messageElement.style.color = "green";
            window.location.href = '/';  // Redirect after successful registration
        } else {
            const errorData = await response.text();  // Get error text
            messageElement.textContent = `Error: ${errorData}`;
            messageElement.style.color = "red";
        }
    } catch (error) {
        messageElement.textContent = "An error occurred during registration.";
        messageElement.style.color = "red";
    }
});