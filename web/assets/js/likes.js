

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

export function reactToggle(post , element){
  
}