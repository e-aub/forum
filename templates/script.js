export const GetData = async () => {
    let target = []
    try {
        let response = await fetch('http://localhost:8000/api');
        if (!response.ok) throw new Error("Network response was not ok");

        let data = await response.json();

        if (data) {
            for (let i = data; i > 0; i--) {
                let link = `http://localhost:8000/api?id=${i}`;
                let postResponse = await fetch(link);
                if (!postResponse.ok) throw new Error("Failed to fetch post data");
                //console.log(typeof postResponse);
                let post = await postResponse.json();
                if (post.PostId !== 0) {
                    target.push(post)
                    RenderPost(target)
                }
            }
        }
    } catch (err) {
        console.error("Error fetching data:", err);
    }
};

function RenderPost(args) {
    const container = document.querySelector(".container");
    container.innerHTML = ""
    args.forEach((element, index) => {
        const Post = document.createElement('div')
        Post.innerHTML = `
        <ul>
        Post N ${index}
        <li>PostId ${element.PostId}</li>
        <li>UserId ${element.UserId}</li>
        <li> title : ${element.Title} </li>
        <li> time : ${element.Created_At} </li>
        <li> Content : ${element.Content} </li> 
        </ul>
        `;
        container.append(Post)
    });
}

async function New_Post() {
    const Botton = document.querySelector('.New_Post')
    Botton.addEventListener('click', (Event) => {
        RenderParam()
        const submit = document.querySelector('.submit_content')
        submit.addEventListener('click', async (event) => {
            let newPost = {}
            newPost.UserId = Math.random()
            newPost.Title = document.querySelector(".Title").value
            newPost.Content = document.querySelector('.Content').value
            newPost.Created_At = Date()
            try {
                const response = await fetch('http://localhost:8000', {
                    method: 'POST',
                    body: JSON.stringify(newPost),
                    headers: {
                        'Content-Type': 'application/json',
                    }
                });

                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                alert('Post created successfully!');
                GetData()
            } catch (error) {
                console.error('Error creating post:', error);
                alert('Failed to create post. Please try again.');
            }
        })
    })
}

function RenderParam() {
    const container = document.querySelector(".container");
    container.innerHTML = ""
    container.innerHTML = `
    New Post :
         <textarea class="Title" placeholder="write your Title"></textarea>
        <textarea class="Content" placeholder="write your Content"></textarea>
        <button class="submit_content" type="submit">Submit</button>`
}

export function Active_Events() {
    New_Post()
}