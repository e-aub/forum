export function makePost( element , reactInfo){
  console.log(reactInfo)
    let liked = false ;
    let disliked = false ;

    let likeCount = reactInfo.liked_by? reactInfo.liked_by.length : 0 ;
    let disLikeCount = reactInfo.disliked_by? reactInfo.disliked_by.length : 0 ; 

    if (!!reactInfo.user_reaction){
      liked = reactInfo.user_reaction === "like"
      disliked = !liked
    }else{
      liked = false 
      disliked = false;
    }
    const innerHTML = `
    <article class="post">
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
                <i class="fas fa-thumbs-up"></i> <span class="count">${likeCount}</span>
            </button>
            <button 
                data-clicked="${liked}" class="dislike" 
                style="background-color: ${disliked ? '#15F5BA' : 'white'};">
                <i class="fas fa-thumbs-down"></i> <span class="count">${disLikeCount}</span>
            </button>
            <button class="comment-button">Comments</button>
          </nav>
        </footer>
    </article> 
    `;
    return innerHTML
  }