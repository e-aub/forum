export function makePost( element , reactInfo){
    if (!!reactInfo.reaction_id == true){
       liked = reactInfo.reaction_id === "like"
       disliked = !liked
    }else{
        liked, disliked = false;
    }
    innerHTML = `
    <article>
        <header>
            <hgroup>
              <h1><i class="fa fa-user"></i> ${element.UserName}</h1>
              <p>${element.CreatedAt}</p>
            <hgroup>
        </header>
        <main>
          <div class="post-content">
              <h2>${element.Title}</h2>
              <p> ${element.Content}</p>
          </div>
        </main>
        <footer>
          <nav>
            <button data-clicked="${liked}" class="like"  
                style="background-color: ${liked ? '#15F5BA' : 'white'};">
                <i class="fas fa-thumbs-up"></i> <span class="count">${element.LikeCount}</span>
            </button>
            <button 
                data-clicked="${liked}" class="dislike" 
                style="background-color: ${disliked ? '#15F5BA' : 'white'};">
                <i class="fas fa-thumbs-down"></i> <span class="count">${element.DislikeCount}</span>
                    </button>
            <button class="comment-button">Comments</button>
          </nav>
        </footer>
    </article> 
    `;
    return innerHTML
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