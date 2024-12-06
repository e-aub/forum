import {showRegistrationModal} from  "./script.js"

export async function handleReact(button, follow , Id, type , target_Type) {
    // Send API request
   const response = getReactInfo({
      "is_own_react": "true",
      "target": "post",
      "target_id": Id,
    }, "PUT").then((response) =>{
      if (!response.ok) {// user is not logged in
          showRegistrationModal(); 
      }else{ // only update the like if no error
          interactiveLike(button, follow)
      }
    }
  );
}

export async function getReactInfo(params, method){
    let queryString = new URLSearchParams(params).toString();
    const url = `http://localhost:8080/react?${queryString}`
    try {
        const response = await fetch(url, {
          method: method,
          headers: {
            'Content-Type': 'application/json'
          }
        });
      return  await response.json();
    } catch (err) {
        console.log(err);
        return;
      }
}

export function reactToggle(post , Id /*post or comment id*/){
    const likeButton = post.querySelector('.like');
    const dislikeButton = post.querySelector('.dislike');

    likeButton.addEventListener('click', () => handleReact(likeButton,dislikeButton, Id, "like", "post"));
    dislikeButton.addEventListener('click', () => handleReact(dislikeButton,likeButton, Id, "dislike", "post"));
}



function interactiveLike(button , follow ){
    const add = button.querySelector(".count");
    const subtract = follow.querySelector(".count");

    // Parse the current count from the button's span text
    let count = parseInt(add.textContent, 10);

    if (button.getAttribute("data-clicked") === "false") {

        count += 1; add.textContent = count; // Update the displayed count
        button.setAttribute("data-clicked", "true");
        button.style.backgroundColor = '#15F5BA'
        follow.style.backgroundColor = 'white'

        if (follow.getAttribute("data-clicked") === "true") {
            count -= 1; subtract.textContent = count; // Update the displayed count
            follow.setAttribute("data-clicked", "false");
            follow.style.backgroundColor = 'white'
        }
    }else if (button.getAttribute("data-clicked") === "true") {
        count -= 1; add.textContent = count; // Update the displayed count
        button.setAttribute("data-clicked", "false");
        button.style.backgroundColor = 'white'
    }
}

