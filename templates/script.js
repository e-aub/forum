export const GetData = async () => {
    console.log("yes");
    try {
        let response = await fetch('http://localhost:8000/api');
        if (!response.ok) throw new Error("Network response was not ok");

        let data = await response.json();

        if (data) {
            for (let i = data; i > 0; i--) {
                let link = `http://localhost:8000/api?id=${i}`;
                let postResponse = await fetch(link);
                if (!postResponse.ok) throw new Error("Failed to fetch post data");

                let post = await postResponse.text();
                if (post.PostId !== 0) {
                    let posts = document.createElement('div');
                    posts.innerHTML = post;
                    document.body.appendChild(posts);
                }
            }
        }
    } catch (err) {
        console.error("Error fetching data:", err);
    }
};
